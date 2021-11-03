package main

import (
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/AlpsMonaco/proxy/socks5"
)

func main() {
	fmt.Println(os.Getpid())
	const port int = 7899
	const addr string = "127.0.0.1"
	s := socks5.Server{
		Address: addr,
		Port:    port,
		OnError: func(err error) {
			fmt.Println(err)
		},

		OnConnectRemote: func(host string, port int) (net.Conn, error) {
			c := socks5.Client{
				Address: addr,
				Port:    7890,
				Timeout: 0,
			}
			err := c.Connect(host, port)
			if err != nil {
				return nil, err
			}

			return &c, nil
		},
	}

	go http.ListenAndServe(":8888", nil)

	if err := s.Listen(); err != nil {
		panic(err)
	}
}
