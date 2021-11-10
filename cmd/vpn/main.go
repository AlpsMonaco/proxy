package main

import (
	"fmt"
	"net"
	_ "net/http/pprof"
	"time"

	"github.com/AlpsMonaco/proxy/socks5"
	"github.com/AlpsMonaco/proxy/vpn"
)

type Decorator struct {
	net.Conn
	c *vpn.Client
}

func (d *Decorator) Write(b []byte) (int, error) {
	err := d.c.Write(b)
	if err != nil {
		return 0, err
	}

	return len(b), nil
}

func (d *Decorator) Read(b []byte) (int, error) {
	nb, err := d.c.Read()
	if err != nil {
		return 0, err
	}
	copy(b, nb)
	return len(nb), nil
}

func main() {
	const ip = "127.0.0.1"
	const port = 7899
	const socksPort = 7898
	var s = &vpn.Server{
		IP:          ip,
		Port:        port,
		Key:         []byte{},
		ErrorHandle: func(err error) { fmt.Println(err) },
	}

	var ss = &socks5.Server{
		Address:         ip,
		Port:            socksPort,
		OnError:         func(err error) { fmt.Println(err) },
		OnClientRequest: nil,
		OnConnectRemote: func(remoteIP string, remotePort int) (net.Conn, error) {
			var c = &vpn.Client{
				ServerIP:    ip,
				ServerPort:  port,
				Key:         []byte{},
				ErrorHandle: func(err error) { fmt.Println(err) },
			}
			if err := c.Connect(remoteIP, remotePort); err != nil {
				return nil, err
			}

			return &Decorator{
				Conn: c.GetConn(),
				c:    c,
			}, nil

		},
	}

	go func() {
		if err := ss.Listen(); err != nil {
			panic(err)
		}
	}()

	go func() {
		if err := s.Listen(); err != nil {
			panic(err)
		}
	}()

	time.Sleep(3600 * time.Second)
	// var c = &vpn.Client{
	// 	ServerIP:    ip,
	// 	ServerPort:  port,
	// 	Key:         []byte{},
	// 	ErrorHandle: func(err error) { fmt.Println(err) },
	// }

	// if err := c.Connect("120.92.17.85", 80); err != nil {
	// 	panic(err)
	// }

	// err := c.Write([]byte(`GET / HTTP/1.1
	// HOST: www.google.com

	// `))
	// if err != nil {
	// 	panic(err)
	// }

	// b, err := c.Read()
	// fmt.Println(err)
	// fmt.Println(string(b))
	// b, err = c.Read()
	// fmt.Println(err)
	// fmt.Println(string(b))
}
