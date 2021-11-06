package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/AlpsMonaco/proxy/vpn"
)

func main() {
	fmt.Println(os.Getpid())
	const port int = 7899
	const addr string = "127.0.0.1"
	s := vpn.Server{
		IP:          addr,
		Port:        7899,
		Key:         []byte{},
		ErrorHandle: func(err error) { fmt.Println(err) },
	}

	go s.Listen()

	var c = vpn.Client{
		ServerIP:    addr,
		ServerPort:  port,
		Key:         []byte{},
		ErrorHandle: func(err error) { fmt.Println(err) },
		Conn:        nil,
	}
	c.Dial()

	time.Sleep(1 * time.Second)
	var b []byte = []byte{10, 0, 10, 10, 10, 10, 10, 10, 10, 10, 6, 0, 6, 6, 6, 6, 5, 0, 5, 5, 5, 17, 0, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17}
	for i := range b {
		c.Write([]byte{b[i]})
	}

	c.Write(b[0:13])
	c.Write(b[13:14])
	c.Write(b[14:21])
	c.Write(b[21:29])
	c.Write(b[29:37])
	time.Sleep(3 * time.Second)
	c.Write(b[37:38])
	time.Sleep(3600 * time.Second)
}
