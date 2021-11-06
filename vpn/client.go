package vpn

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIP    string
	ServerPort  int
	Key         []byte
	ErrorHandle func(err error)
	net.Conn
}

func (c *Client) Dial() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.ServerIP, c.ServerPort))
	if err != nil {
		return err
	}
	c.Conn = conn
	return nil
}
