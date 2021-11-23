package socks5

import (
	"fmt"
	"net"

	"github.com/AlpsMonaco/proxy/forward"
)

type ProxyConn interface {
	Proxy(client net.Conn)
}

type direct struct{ remote net.Conn }

func (p *direct) Proxy(client net.Conn) {
	forward.NewForward(client, p.remote, func(err error) { fmt.Println(err) }).Start()
}
