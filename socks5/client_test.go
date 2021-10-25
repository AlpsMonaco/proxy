package socks5

import (
	"fmt"
	"testing"

	"github.com/AlpsMonaco/proxy/util"
)

func TestAnd(t *testing.T) {
	var i int = 80
	t.Log(i & 0xFF)
	t.Log(i & 0x0F)
	t.Log(i & 0xF0)
}

func TestClient(t *testing.T) {
	var err error
	var c = &Client{
		Address: "127.0.0.1",
		Port:    7890,
		Timeout: 0,
	}

	err = c.Connect("8.134.75.115", 80)
	assert(err)

	_, err = c.Write([]byte(`GET / HTTP/1.1
HOST: www.google.com

`))
	assert(err)

	var a util.Alloctor
	a.Alloc(256)

	for {
		n, err := c.Read(a.GetBytes())
		fmt.Print(string(a.GetBytes()))
		if n < 256 {
			break
		}
		assert(err)
	}
}
