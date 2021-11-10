package vpn

import (
	"errors"
	"fmt"
	"net"
)

type Client struct {
	ServerIP    string
	ServerPort  int
	Key         []byte
	ErrorHandle func(err error)
	conn        net.Conn
	p           *Packet
}

func (c *Client) dial() (err error) {
	c.conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", c.ServerIP, c.ServerPort))
	if err != nil {
		return err
	}

	var p Packet
	c.p = &p
	p.Init()
	return nil
}

// func (c *Client) onError(err error) {
// 	if c.ErrorHandle != nil {
// 		c.ErrorHandle(err)
// 	}
// }

func (c *Client) Read() (b []byte, err error) {
	err = c.p.Next(c.conn)
	if err != nil {
		return nil, err
	}
	return c.p.GetData(), nil
}

func (c *Client) GetConn() net.Conn {
	return c.conn
}

func (c *Client) Write(b []byte) (err error) {
	if err = c.p.WriteSize(c.conn, len(b)); err != nil {
		return err
	}

	_, err = c.conn.Write(b)
	return err
}

func (c *Client) Connect(ip string, port int) error {
	var err error
	if err = c.dial(); err != nil {
		return err
	}

	v := (*Verify)(c.p.GetPointer())
	v.SetKey()
	if err = c.p.WriteBuffer(c.conn, 16); err != nil {
		return err
	}

	pr := (*ProxyRequest)(c.p.GetPointer())
	pr.SetInfo(ip, port)
	err = c.p.WriteBuffer(c.conn, pr.GetSize())
	if err != nil {
		return err
	}

	if err = c.p.Next(c.conn); err != nil {
		return err
	}

	gr := (*GeneralResponse)(c.p.GetPointer())
	if gr.Code != Success {
		return errors.New(string(gr.Msg[:gr.MsgSize]))
	}

	return nil
}
