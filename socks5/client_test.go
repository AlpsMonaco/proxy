package socks5

import (
	"testing"
	"unsafe"

	"github.com/AlpsMonaco/proxy/util"
)

func TestVersionMessageMethods(t *testing.T) {
	t.Log(unsafe.Alignof(interface{}(1)))

	var msg = Socks5_VersionMessage{
		Ver:       5,
		NumMethod: 0,
		va:        [256]byte{},
	}

	t.Log(msg.GetMethods())
	msg.SetMethod(0, 1, 2)
	t.Log(msg.GetMethods())
	t.Log(msg)

	t.Log(util.ToBinary(&msg, int(2+msg.NumMethod)))

	var b []byte = []byte{5, 7, 0, 1, 2, 3, 4, 5, 6}
	var p *Socks5_VersionMessage
	p = (*Socks5_VersionMessage)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&b))))
	t.Log(p.GetMethods())

}
