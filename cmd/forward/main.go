package main

import (
	"fmt"
	"net"
	"time"

	"github.com/AlpsMonaco/proxy/forward"
)

func main() {
	l, err := net.Listen("tcp", ":7898")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		go serveConn(conn)
	}
}

func serveConn(conn net.Conn) {
	dst, err := net.DialTimeout("tcp", "127.0.0.1:7890", 5*time.Second)
	if err != nil {
		fmt.Println(err)
	}

	f := forward.Forward{
		SrcConn: conn,
		DstConn: dst,
		OnError: func(e error) { fmt.Println(e) },
	}

	f.Start()
}
