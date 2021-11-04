package socks5

import (
	"fmt"
	"testing"

	"github.com/AlpsMonaco/proxy/util"
)

func TestServer(t *testing.T) {
	var s Server = Server{
		Address: "127.0.0.1",
		Port:    7899,
		OnError: func(err error) {
			fmt.Println(err)
		},
	}
	s.Listen()
}

func pointerMethod(p uptr) {
	c := (*Socks5_RequestMessage)(p)
	fmt.Println(c)
}

func TestPointer(t *testing.T) {
	var a util.Allocator
	a.Alloc(264)

	c := (*Socks5_RequestMessage)(a.GetPointer())
	c.va[0] = 128
	c.Atype = 2
	c.Cmd = 3
	c.Ver = 4
	c.Rsv = 0x09

	pointerMethod(a.GetPointer())
}
