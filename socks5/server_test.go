package socks5

import (
	"fmt"
	"testing"
)

func TestServer(t *testing.T) {
	var s Server = Server{
		Address: "127.0.0.1",
		Port:    7899,
		OnError: func(err error) {
			fmt.Println(err)
		},
	}
	s.Listen()
}
