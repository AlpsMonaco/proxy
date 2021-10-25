package main

import (
	"fmt"
	"github.com/AlpsMonaco/proxy/socks5"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	const port int = 7899
	const addr string = "127.0.0.1"
	s := socks5.Server{
		Address: addr,
		Port:    port,
		OnError: func(err error) {
			fmt.Println(err)
		},
	}

	go http.ListenAndServe(":8888", nil)

	if err := s.Listen(); err != nil {
		panic(err)
	}
}
