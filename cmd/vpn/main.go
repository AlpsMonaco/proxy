package main

import (
	"fmt"

	"github.com/AlpsMonaco/proxy/socks5"
	"github.com/AlpsMonaco/proxy/vpn"
)

func main() {
	go ServerSide()

	go func() {
		sockServer := socks5.Server{
			Address: ServerIP,
			Port:    SockPort,
			OnClientConnect: func(c *socks5.ClientConn) {
				// rm :=(*vpn.RequestMessage)(c.Allocator.GetPointer())

			},
			OnError: func(err error) { fmt.Println(err) },
		}

		if err := sockServer.Listen(); err != nil {
			panic(err)
		}
	}()
}

const key = "123456"
const ServerIP = "127.0.0.1"
const Port = 61124

const SockPort = 7899

func ServerSide() {
	var s = vpn.Server{
		IP:       ServerIP,
		Port:     Port,
		Password: key,
		OnError:  func(err error) { fmt.Println(err) },
	}

	if err := s.Serve(); err != nil {
		panic(err)
	}
}
