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

	if err := c.Connect("120.92.17.85", 80); err != nil {
		panic(err)
	}

	err := c.Write([]byte(`GET / HTTP/1.1
	HOST: www.google.com

	`))
	if err != nil {
		panic(err)
	}

	b, err := c.Read()
	fmt.Println(err)
	fmt.Println(string(b))
	b, err = c.Read()
	fmt.Println(err)
	fmt.Println(string(b))
}
