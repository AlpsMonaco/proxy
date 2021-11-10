package vpn

import (
	"fmt"
	"net"
	"unsafe"

	"github.com/AlpsMonaco/proxy/socks5"
	"github.com/AlpsMonaco/proxy/stream"
	"github.com/AlpsMonaco/proxy/util"
)

type Client struct {
	ServerIP    string
	ServerPort  int
	Key         []byte
	ErrorHandle func(err error)
	p           *stream.Packet
}

func (c *Client) Dial() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.ServerIP, c.ServerPort))
	if err != nil {
		return err
	}
	c.p = &stream.Packet{
		Conn: conn,
	}
	return nil
}

func (c *Client) Write(b []byte) (int, error) {
	return c.p.Write(b)
}

func (c *Client) Read(b []byte) (int, error) {
	return c.p.Read(b)
}

func (c *Client) onError(err error) {
	if c.ErrorHandle != nil {
		c.ErrorHandle(err)
	}
}

func (c *Client) Connect(ip string, port int) error {
	a := util.GetAlloctor(stream.PacketSize)
	defer util.FreeAllocator(a)
	// var n int
	var err error

	reqMsg := (*socks5.Socks5_RequestMessage)(a.GetPointer())
	socks5.FillRequestMessage(reqMsg, 0, ip, port)
	_, err = c.Write(a.GetByteSize(reqMsg.GetSize()))
	if err != nil {
		return err
	}
	_, err = c.Read(a.GetBytes())
	if err != nil {
		return err
	}

	r := (*Protocol_Response)(unsafe.Pointer(&c.p.Body[0]))
	fmt.Println(r.Code)
	fmt.Println(r.MsgSize)
	fmt.Println(string(r.Msg[:r.MsgSize]))
	return nil
}
