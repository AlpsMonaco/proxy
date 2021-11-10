package vpn

import (
	"fmt"
	"io"
	"net"
	"time"
	"unsafe"

	"github.com/AlpsMonaco/proxy/socks5"
	"github.com/AlpsMonaco/proxy/stream"
	"github.com/AlpsMonaco/proxy/util"
)

type Server struct {
	IP          string
	Port        int
	Key         []byte
	ErrorHandle func(err error)
	l           net.Listener
}

func (s *Server) Listen() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		return err
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		s.newConn(conn)
	}
}

func (s *Server) onError(err error) {
	if s.ErrorHandle != nil {
		s.onError(err)
	}
}

func (s *Server) newConn(conn net.Conn) {
	defer closeConn(conn)
	var a *util.Allocator = util.GetAlloctor(stream.PacketSize)
	defer util.FreeAllocator(a)
	p := stream.GetPacket()
	p.Conn = conn
	defer stream.FreePacket(p)

	var err error
	_, err = p.Read(a.GetBytes())
	if err != nil {
		s.onError(err)
		return
	}

	host := (*socks5.Socks5_RequestMessage)(unsafe.Pointer(&p.Body[0])).GetHost()
	port := (*socks5.Socks5_RequestMessage)(unsafe.Pointer(&p.Body[0])).GetPort()
	var remote net.Conn
	remote, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 10*time.Second)
	resp := (*Protocol_Response)(a.GetPointer())
	if err != nil {
		s.onError(err)
		resp.Code = Failed
		resp.FillMsg("失败")
	}
	defer closeConn(remote)
	resp.Code = Success
	resp.FillMsg("成功")
	_, err = p.Write(a.GetByteSize(resp.GetSize()))
	if err != nil {
		s.onError(err)
		return
	}
	s.beginProxy(p, remote, a)
}

func (s *Server) beginProxy(client, remote net.Conn, a *util.Allocator) {
	defer closeConn(client)
	defer closeConn(remote)
	var n int
	var err error

	go func() {
		for {
			n, err = remote.Read(a.GetBytes())
			if n == 0 {
				err = io.EOF
			}
			if err != nil {
				s.onError(err)
				return
			}
			n, err = client.Write(a.GetByteSize(n))
			if err != nil {
				s.onError(err)
				return
			}
		}
	}()

	func() {
		for {
			n, err = client.Read(a.GetBytes())
			if n == 0 {
				err = io.EOF
			}
			if err != nil {
				s.onError(err)
				return
			}
			n, err = remote.Write(a.GetByteSize(n))
			if err != nil {
				s.onError(err)
				return
			}
		}
	}()
}

func closeConn(conn net.Conn) {
	if err := conn.Close(); err != nil {
		conn.Close()
	}
}
