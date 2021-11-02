package vpn

import (
	"fmt"
	"net"
	"time"
	"unsafe"

	"github.com/AlpsMonaco/proxy/crypto"
	"github.com/AlpsMonaco/proxy/forward"
)

type Server struct {
	IP       string
	Port     int
	Password string
	OnError  func(error)
	e        crypto.Encryptor
}

func (s *Server) onError(err error) {
	if s.OnError != nil {
		s.OnError(err)
	}
}

func (s *Server) Serve() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		return err
	}

	s.e = new(crypto.AESEncryptor).Key([]byte(s.Password))

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		go s.newConn(conn)
	}
}

func (s *Server) newConn(conn net.Conn) {
	netBuf := make([]byte, ServerBufSize)
	var size int
	var err error

	size, err = conn.Read(netBuf)
	if err != nil {
		s.onError(err)
		closeConn(conn)
		return
	}

	b, err := s.e.Decrypt(netBuf[:size])
	if err != nil {
		s.onError(err)
		closeConn(conn)
		return
	}

	host, err := (*RequestMessage)(unsafe.Pointer(&b[0])).Parse()
	if err != nil {
		s.onError(err)
		closeConn(conn)
		return
	}

	remote, err := net.DialTimeout("tcp", host, 10*time.Second)
	if err != nil {
		s.onError(err)
		closeConn(conn)
		return
	}

	forward.NewForward(remote, &SecureConn{
		Conn: conn,
		e:    s.e,
	}, func(e error) { fmt.Println(e) }).Start()
}

func closeConn(c net.Conn) {
	err := c.Close()
	if err != nil {
		_ = c.Close()
	}
}
