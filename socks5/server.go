package socks5

import (
	"fmt"
	"net"
	"reflect"
	"unsafe"

	"github.com/AlpsMonaco/proxy/forward"
	"github.com/AlpsMonaco/proxy/util"
)

type uptr = unsafe.Pointer

type Server struct {
	Address string
	Port    int
	OnError func(error)

	OnClientRequest func(rm *Socks5_RequestMessage) error
	OnConnectRemote func(host string, port int) (net.Conn, error)
	OnProxy         func(client, remote net.Conn)

	listener net.Listener
}

func (s *Server) Listen() error {
	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", s.Address, s.Port))
	if err != nil {
		return err
	}

	if s.OnClientRequest == nil {
		s.OnClientRequest = func(rm *Socks5_RequestMessage) error {
			if rm.Ver != SOCKS5_VERSION {
				return ErrSocks5VersionNotSupported
			}
			if rm.Cmd != SOCKS5_CMD_CONNECT {
				return ErrSocks5CommandNotSupported
			}
			return nil
		}
	}

	if s.OnConnectRemote == nil {
		s.OnConnectRemote = func(host string, port int) (net.Conn, error) {
			conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
			return conn, err
		}
	}

	if s.OnProxy == nil {
		s.OnProxy = func(client, remote net.Conn) {
			forward.NewForward(client, remote, s.onError).Start()
		}
	}

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		go s.newConn(conn)
	}
}

func (s *Server) onError(err error) {
	if s.OnError != nil {
		s.OnError(err)
	}
}

func (s *Server) newConn(c net.Conn) {
	var err error
	a := util.GetAlloctor(264)

	_, err = c.Read(a.GetBytes())
	if err != nil {
		s.onError(err)
		closeConn(c)
		return
	}

	err = parseVersionMessage((*Socks5_VersionMessage)(a.GetPointer()))
	if err != nil {
		s.onError(err)
		closeConn(c)
		return
	}

	fillSelectionMessage((*Socks5_SelectionMessage)(a.GetPointer()))
	_, err = c.Write(a.GetByteSize(2))
	if err != nil {
		s.onError(err)
		closeConn(c)
		return
	}

	_, err = c.Read(a.GetBytes())
	if err != nil {
		s.onError(err)
		closeConn(c)
		return
	}

	rMsg := (*Socks5_RequestMessage)(a.GetPointer())
	err = s.OnClientRequest(rMsg)
	if err != nil {
		s.onError(err)
		closeConn(c)
		return
	}

	remote, err := s.OnConnectRemote(rMsg.GetHost(), rMsg.GetPort())
	if err != nil {
		s.onError(err)
		closeConn(c)
		return
	}

	respMsg := (*Socks5_ResponseMessage)(a.GetPointer())
	fillResponseMessage(respMsg, s)
	_, err = c.Write(a.GetByteSize(respMsg.GetSize()))
	if err != nil {
		s.onError(err)
		closeConn(c)
		return
	}

	util.FreeAllocator(a)
	s.OnProxy(c, remote)
}

func parseVersionMessage(vMsg *Socks5_VersionMessage) error {
	if vMsg.Ver != SOCKS5_VERSION {
		return ErrSocks5VersionNotSupported
	}
	return nil
}

func fillSelectionMessage(sMsg *Socks5_SelectionMessage) {
	sMsg.Ver = SOCKS5_VERSION
	sMsg.Method = SOCKS5_METHOD_NO_AUTH
}

func fillResponseMessage(respMsg *Socks5_ResponseMessage, s *Server) {
	respMsg.Ver = SOCKS5_VERSION
	respMsg.Rsv = 0x00
	respMsg.Rep = SOCKS5_REP_SUCCESS
	var i byte
	if isDomain(s.Address) {
		respMsg.Atype = SOCKS5_ATYPE_DOMAIN
		respMsg.va[0] = byte(len(s.Address))
		for i = 1; i < respMsg.va[0]+1; i++ {
			respMsg.va[i] = s.Address[i-1]
		}
	} else {
		respMsg.Atype = SOCKS5_ATYPE_IPV4
		util.IPV4AddrToByte(s.Address, (*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
			Data: uintptr(unsafe.Pointer(&respMsg.va[0])),
			Len:  4,
			Cap:  4,
		})))
		i = i + 4
	}
	respMsg.va[i] = byte(s.Port >> 8)
	respMsg.va[i+1] = byte(s.Port & 0x00FF)
}

func closeConn(c net.Conn) {
	err := c.Close()
	if err != nil {
		_ = c.Close()
	}
}
