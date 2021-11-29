package vpn

import (
	"fmt"
	"net"
)

type CipherEnum byte

const (
	Cipher_Plain CipherEnum = iota
	Cipher_AES256GCM
	Cipher_Chacha20poly1305
)

type Server struct {
	Addr        string
	Port        int
	Key         []byte
	Cipher      CipherEnum
	ErrorHandle func(err error)
	encryptor   Encryptor
}

func (s *Server) Listen() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Addr, s.Port))
	if err != nil {
		return err
	}
	s.encryptor = GetEncryptor(s.Cipher, s.Key)
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go s.newConn(conn)
	}
}

func (s *Server) newConn(client net.Conn) {

}

func (s *Server) onError(err error) {
	if s.ErrorHandle != nil {
		s.ErrorHandle(err)
	}
}

func closeConn(conn net.Conn) {
	if err := conn.Close(); err != nil {
		conn.Close()
	}
}
