package socks5

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"time"
	"unsafe"

	"github.com/AlpsMonaco/proxy/util"
)

var ErrSocks5VersionNotSupported = errors.New("Socks5 version not supported")
var ErrSocks5MethodNotSupported = errors.New("Socks5 method not supported")
var ErrSocks5CommandNotSupported = errors.New("Socks5 command not supported")
var ErrRemoteNotAllowed = errors.New("request host not allowed")

type Client struct {
	Address  string
	Port     int
	Timeout  time.Duration
	User     string
	Password string
	net.Conn
}

func (c *Client) Connect(addr string, port int) error {
	var err error
	var a util.Alloctor
	a.Alloc(264)

	if err = c.dial(&a); err != nil {
		return err
	}

	fillRequestMessage((*Socks5_RequestMessage)(a.GetPointer()), SOCKS5_CMD_CONNECT, addr, port)
	_, err = c.Write(a.GetByteSize((*Socks5_RequestMessage)(a.GetPointer()).GetSize()))
	if err != nil {
		return err
	}

	_, err = c.Read(a.GetBytes())
	if err != nil {
		return err
	}

	parseResponseMessage((*Socks5_ResponseMessage)(a.GetPointer()))
	return nil
}

func (c *Client) dial(a *util.Alloctor) error {
	var err error
	c.Conn, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d", c.Address, c.Port), c.Timeout)
	if err != nil {
		return err
	}

	fillVersionMessage(c, (*Socks5_VersionMessage)(a.GetPointer()))
	_, err = c.Write(a.GetByteSize((*Socks5_VersionMessage)(a.GetPointer()).GetSize()))
	if err != nil {
		return err
	}

	_, err = c.Read(a.GetBytes())
	if err != nil {
		return err
	}

	if err = parseSelectionMessage((*Socks5_SelectionMessage)(a.GetPointer())); err != nil {
		return err
	}

	return nil
}

func fillVersionMessage(c *Client, vMsg *Socks5_VersionMessage) {
	vMsg.Ver = SOCKS5_VERSION
	vMsg.NumMethod = 0x01
	vMsg.va[0] = 0x00
}

func parseSelectionMessage(vMsg *Socks5_SelectionMessage) error {
	if vMsg.Ver != SOCKS5_VERSION {
		return ErrSocks5VersionNotSupported
	}

	if vMsg.Method >= SOCKS5_METHOD_NOT_SUPPORT {
		return ErrSocks5MethodNotSupported
	}

	return nil
}

func fillRequestMessage(vMsg *Socks5_RequestMessage, sockcmd byte, addr string, port int) {
	vMsg.Ver = SOCKS5_VERSION
	vMsg.Cmd = sockcmd
	vMsg.Rsv = 0x00
	var i int
	if isDomain(addr) {
		// domain
		vMsg.Atype = SOCKS5_ATYPE_DOMAIN
		vMsg.va[0] = byte(len(addr))
		for i = 1; i < len(addr)+1; i++ {
			vMsg.va[i] = addr[i-1]
		}
	} else {
		// ipv4
		vMsg.Atype = SOCKS5_ATYPE_IPV4
		util.IPV4AddrToByte(addr, (*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
			Data: uintptr(unsafe.Pointer(&vMsg.va[0])),
			Len:  4,
			Cap:  4,
		})))
		i = i + 4
	}

	vMsg.va[i] = byte(port >> 8)
	vMsg.va[i+1] = byte(port & 0x00FF)
}

func parseResponseMessage(rMsg *Socks5_ResponseMessage) {
	// fmt.Println(rMsg)
}

func isDomain(s string) bool {
	for _, v := range []byte(s) {
		if v >= 58 {
			return true
		}
	}

	return false
}
