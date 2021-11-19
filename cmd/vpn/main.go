package main

import (
	"fmt"
	"net"
	"time"

	"github.com/AlpsMonaco/proxy/socks5"
	"github.com/AlpsMonaco/proxy/vpn"
)

const ServerIP = "127.0.0.1"
const Port = 7899
const SocksPort = 7898

func main() {
	var s vpn.Server = vpn.Server{
		Addr:        ServerIP,
		Port:        Port,
		Key:         []byte{},
		Cipher:      0,
		ErrorHandle: func(err error) { fmt.Println(err) },
	}

	go func() {
		if err := s.Listen(); err != nil {
			panic(err)
		}
	}()
	time.Sleep(1 * time.Second)
	fmt.Println("server start.")

	var socksServer = socks5.Server{
		Address:         ServerIP,
		Port:            SocksPort,
		OnError:         func(err error) { fmt.Println(err) },
		OnClientRequest: nil,
		OnConnectRemote: func(host string, port int) (net.Conn, error) {
			var c = vpn.Client{
				IP:          ServerIP,
				Port:        Port,
				Key:         []byte{},
				Cipher:      0,
				ErrorHandle: func(err error) { fmt.Println(err) },
			}
			err := c.Connect(host, port)
			if err != nil {
				return nil, err
			}
			return c.Conn(), nil
		},
		OnProxy: func(client net.Conn, remote net.Conn, s *socks5.Server) {
			// defer client.Close()
			// defer remote.Close()

			// a := util.GetAlloctor(stream.PacketSize)
			// defer util.FreeAllocator(a)
			// var n int
			// var err error
			// var clientBuffer = a.GetByteSize(stream.PacketSize - (1 << 8))
			// var packet = stream.Packet{
			// 	Stream: remote,
			// }
			// packet.SetAllocator(a)

			// go func() {
			// 	defer client.Close()
			// 	for {
			// 		n, err = client.Read(clientBuffer)
			// 		if err == nil && n == 0 {
			// 			err = io.EOF
			// 		}
			// 		if err != nil {
			// 			s.OnError(err)
			// 			return
			// 		}
			// 	}
			// }()

		},
	}

}
