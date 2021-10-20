package socks5

import (
	"net"
	"unsafe"

	"github.com/AlpsMonaco/proxy/util"
)

type Client struct {
}

func SendVersionMessage(conn net.Conn, msg *Socks5_VersionMessage) (int, error) {
	size := 2 + msg.NumMethod
	b := util.ToBinary(msg, int(size))
	return conn.Write(*b)
}

func RecvSelectionMessage(b []byte) *Socks5_SelectionMessage {
	return (*Socks5_SelectionMessage)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&b))))
}
