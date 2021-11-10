package main

import (
	"fmt"
	_ "net/http/pprof"
	"time"

	"github.com/AlpsMonaco/proxy/vpn"
)

func main() {
	const ip = "127.0.0.1"
	const port = 7899
	var s = &vpn.Server{
		IP:          ip,
		Port:        port,
		Key:         []byte{},
		ErrorHandle: func(err error) { fmt.Println(err) },
	}

	go func() {
		if err := s.Listen(); err != nil {
			fmt.Println(err)
		}
	}()

	time.Sleep(1 * time.Second)
	var c = &vpn.Client{
		ServerIP:    ip,
		ServerPort:  port,
		Key:         []byte{},
		ErrorHandle: func(err error) { fmt.Println(err) },
	}

	if err := c.Dial(); err != nil {
		panic(err)
	}

	if err := c.Connect("120.92.17.85", 80); err != nil {
		panic(err)
	}

	_, err := c.Write([]byte(`GET / HTTP/1.1
	HOST: www.google.com
	
	`))
	if err != nil {
		panic(err)
	}
	var b = make([]byte, 1024)
	c.Read(b)
	fmt.Println(string(b[2:]))

}
