package main

import (
	"fmt"
	"time"

	"github.com/AlpsMonaco/proxy/socks5"
	"github.com/AlpsMonaco/proxy/vpn"
)

const ListenAddr = "127.0.0.1"
const VpnPort = 7899
const SocksPort = 7898

func main() {
	var server = vpn.Server{
		Addr:        ListenAddr,
		Port:        VpnPort,
		Key:         []byte{},
		Cipher:      0,
		ErrorHandle: func(err error) { fmt.Println(err) },
	}

	go func() {
		err := server.Listen()
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(1 * time.Second)
	fmt.Println("server start")

	// var client = vpn.Client{
	// 	IP:          ListenAddr,
	// 	Port:        VpnPort,
	// 	Key:         []byte{},
	// 	Cipher:      0,
	// 	ErrorHandle: func(err error) { fmt.Println(err) },
	// }
	// err := client.Connect("www.baidu.com", 80)
	// if err != nil {
	// 	panic(err)
	// }

	// c := client.Conn()
	// msg := "GET / HTTP/1.1\r\nHost: www.baidu.com\r\n\r\n"
	// var p = stream.NewPacket()
	// p.WriteStream(c, []byte(msg))

	// err = p.Next(c)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(p.Data()))

	// time.Sleep(100 * time.Second)

	var sockserver = socks5.Server{
		Address: ListenAddr,
		Port:    SocksPort,
		OnError: func(err error) { fmt.Println(err) },
		// OnClientRequest: func(*socks5.Socks5_RequestMessage) error { return nil },
		OnConnectRemote: func(ip string, port int) (socks5.ProxyConn, error) {
			var client = vpn.Client{
				IP:          ListenAddr,
				Port:        VpnPort,
				Key:         []byte{},
				Cipher:      0,
				ErrorHandle: func(err error) { fmt.Println(err) },
			}
			err := client.Connect(ip, port)
			if err != nil {
				return nil, err
			}

			return &client, nil
		},
	}

	if err := sockserver.Listen(); err != nil {
		panic(err)
	}

}
