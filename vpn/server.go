package vpn

import (
	"fmt"
	"net"

	"github.com/AlpsMonaco/proxy/stream"
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
	s.l = l
	stream.Serve(s)
	return nil
}

func (s *Server) Accept() net.Conn {
	conn, err := s.l.Accept()
	if err != nil {
		s.OnError(err)
		return nil
	}
	return conn
}

func (s *Server) OnError(err error) {
	if s.ErrorHandle != nil {
		s.ErrorHandle(err)
	}
}

func (s *Server) Handle(b []byte, c net.Conn) {
	fmt.Println(b)
}
