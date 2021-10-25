package socks5

import (
	"fmt"
	"github.com/AlpsMonaco/proxy/forward"
	"net"
	"reflect"
	"unsafe"

	"github.com/AlpsMonaco/proxy/util"
)

type Server struct {
	Address             string
	Port                int
	BeforeClientConnect func(*ClientConn) bool
	OnClientConnect     func(*ClientConn)
	OnError             func(error)
	listener            net.Listener
}

type ClientConn struct {
	Addr   string
	Port   int
	Remote net.Conn
	Client net.Conn
}

func DefaultOnClientConnect(c *ClientConn) {
	forward.NewForward(c.Remote, c.Client, nil).Start()
}

func (s *Server) Listen() error {
	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", s.Address, s.Port))
	if err != nil {
		return err
	}

	if s.BeforeClientConnect == nil {
		s.BeforeClientConnect = func(client *ClientConn) bool {
			return true
		}
	}

	if s.OnClientConnect == nil {
		s.OnClientConnect = func(c *ClientConn) {
			forward.NewForward(c.Remote, c.Client, s.OnError).Start()
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
	var a util.Alloctor
	a.Alloc(264)

	_, err = c.Read(a.GetBytes())
	if err != nil {
		s.onError(err)
		return
	}

	err = parseVersionMessage((*Socks5_VersionMessage)(a.GetPointer()))
	if err != nil {
		s.onError(err)
		return
	}

	fillSelectionMessage((*Socks5_SelectionMessage)(a.GetPointer()))
	_, err = c.Write(a.GetByteSize(2))
	if err != nil {
		s.onError(err)
		return
	}

	_, err = c.Read(a.GetBytes())
	if err != nil {
		s.onError(err)
		return
	}

	rMsg := (*Socks5_RequestMessage)(a.GetPointer())
	err = ParseRequestMessage(rMsg)
	if err != nil {
		s.onError(err)
		return
	}

	var clientConn = ClientConn{
		Addr:   rMsg.GetHost(),
		Port:   rMsg.GetPort(),
		Client: c,
	}

	if !s.BeforeClientConnect(&clientConn) {
		_ = c.Close()
		return
	}

	respMsg := (*Socks5_ResponseMessage)(a.GetPointer())
	fillResponseMessage(respMsg, s)
	clientConn.Remote, err = net.Dial("tcp", fmt.Sprintf("%s:%d", clientConn.Addr, clientConn.Port))
	if err != nil {
		s.OnError(err)
		respMsg.Rep = SOCKS5_REP_CONNECTION_FAILED
		_, err = c.Write(a.GetByteSize(respMsg.GetSize()))
		if err != nil {
			s.OnError(err)
		}
		_ = c.Close()
		return
	}

	respMsg.Rep = SOCKS5_REP_SUCCESS
	_, err = c.Write(a.GetByteSize(respMsg.GetSize()))
	if err != nil {
		s.OnError(err)
		_ = c.Close()
		return
	}

	s.OnClientConnect(&clientConn)
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

func ParseRequestMessage(rMsg *Socks5_RequestMessage) error {
	if rMsg.Ver != SOCKS5_VERSION {
		return ErrSocks5VersionNotSupported
	}
	if rMsg.Cmd != SOCKS5_CMD_CONNECT {
		return ErrSocks5CommandNotSupported
	}
	return nil
}

func fillResponseMessage(respMsg *Socks5_ResponseMessage, s *Server) {
	respMsg.Ver = SOCKS5_VERSION
	respMsg.Rsv = 0x00
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
	respMsg.va[i] = byte(s.Port & 0x0F)
	respMsg.va[i+1] = byte(s.Port & 0xF0)
}
