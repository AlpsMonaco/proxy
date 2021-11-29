package vpn

import (
	"fmt"
	"net"
)

type Client struct {
	IP          string
	Port        int
	Key         []byte
	Cipher      CipherEnum
	ErrorHandle func(err error)
	encryptor   Encryptor
	s           net.Conn
}

func (c *Client) GetCipher() Encryptor {
	return c.encryptor
}

func (c *Client) dial() (err error) {
	c.s, err = net.Dial("tcp", fmt.Sprintf("%s:%d", c.IP, c.Port))
	c.encryptor = GetEncryptor(c.Cipher, c.Key)
	return err
}

func (c *Client) onError(err error) {
	if c.ErrorHandle != nil {
		c.ErrorHandle(err)
	}
}
