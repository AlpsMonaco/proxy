package main

import (
	"fmt"
	"time"

	"github.com/AlpsMonaco/proxy/vpn"
)

func main() {
	go ServerSide()
	time.Sleep(1 * time.Second)
	var c = vpn.Client{
		IP:       ServerIP,
		Port:     Port,
		Password: key,
	}

	if err := c.Connect("119.91.84.155", 80); err != nil {
		panic(err)
	}

	time.Sleep(1 * time.Second)

	_, err := c.Write([]byte(`GET / HTTP/1.1
Host: www.baidu.com

`))

	if err != nil {
		panic(err)
	}

	buf := make([]byte, 1024)
	_, err = c.Read(buf)
	if err != nil {
		panic(err)
	}

	fmt.Println(buf)
	fmt.Println(string(buf))

}

const key = "123456"
const ServerIP = "127.0.0.1"
const Port = 61124

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
