package vpn

import (
	"fmt"
	"net"

	"github.com/AlpsMonaco/proxy/encrypt"
)

type Server struct {
	IP        string
	Port      int
	Encryptor encrypt.Encryptor
}

func (s *Server) Serve() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		return err
	}

}
