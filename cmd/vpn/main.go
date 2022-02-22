package main

import (
	"fmt"

	"github.com/AlpsMonaco/proxy/socks5"
	"github.com/AlpsMonaco/proxy/vpn"
)

const socks5Port int = 7899
const addr string = "127.0.0.1"
const vpnPort = 7179
const key = "123"

func main() {
	go func() {
		var vpnServer = vpn.Server{IP: addr, Port: vpnPort, Key: key}
		err := vpnServer.Serve()
		if err != nil {
			panic(err)
		}
	}()

	s := socks5.Server{
		Address: addr,
		Port:    socks5Port,
		OnError: func(err error) {
			fmt.Println(err)
		},

		OnConnectRemote: func(host string, port int) (socks5.ProxyConn, error) {
			var c vpn.Client = vpn.Client{
				ServerIP:   addr,
				ServerPort: vpnPort,
				Key:        key,
			}
			err := c.Connect(host, port)
			if err != nil {
				return nil, err
			}
			return &c, nil
		},
	}

	if err := s.Listen(); err != nil {
		panic(err)
	}
}
