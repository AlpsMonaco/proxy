package vpn

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/AlpsMonaco/proxy/stream"
	"github.com/AlpsMonaco/proxy/util"
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

func (c *Client) Connect(host string, port int) error {
	var err error
	err = c.dial()
	if err != nil {
		return err
	}

	var allocator *util.Allocator = util.GetAlloctor(256)
	defer util.FreeAllocator(allocator)
	var n int

	(*HelloMessage)(allocator.GetPointer()).SetMsg(masquerade)

	_, err = c.s.Write(allocator.GetBytes())
	if err != nil {
		return err
	}
	_, err = c.s.Read(allocator.GetBytes())
	if err != nil {
		return err
	}

	// fmt.Println(allocator.GetByteSize(n))
	n, err = c.encryptor.Encrypt([]byte(masquerade), allocator.GetBytes())
	if err != nil {
		return err
	}

	_, err = c.s.Write(allocator.GetByteSize(n))
	if err != nil {
		return err
	}

	(*ProxyRequest)(allocator.GetPointer()).SetRemoteInfo(host, port)
	_, err = c.s.Write(allocator.GetByteSize(256))
	if err != nil {
		return err
	}

	_, err = c.s.Read(allocator.GetBytes())
	if err != nil {
		return err
	}

	if (*GeneralResponse)(allocator.GetPointer()).code != Code_Success {
		return errors.New((*GeneralResponse)(allocator.GetPointer()).Get())
	} else {
		// fmt.Println((*GeneralResponse)(allocator.GetPointer()).Get())
	}

	return nil
}

func (c *Client) Conn() net.Conn {
	return c.s
}

func (c *Client) Proxy(client net.Conn) {
	var allocator *util.Allocator = util.GetAlloctor(stream.PacketSize)
	defer util.FreeAllocator(allocator)
	var packet *stream.Packet = stream.NewPacket()
	defer stream.FreePacket(packet)

	defer closeConn(c.s)
	defer closeConn(client)
	var n int
	var err error

	client = &debugconn{client, "socks5_client"}
	c.s = &debugconn{c.s, "vpn_server"}

	var clientBuffer = allocator.GetByteSize(stream.PacketSize - (1 << 8))
	var serverBuffer = allocator.GetBytes()

	go func() {
		defer closeConn(client)
		defer closeConn(c.s)
		for {
			n, err = client.Read(clientBuffer)
			if n == 0 && err == nil {
				err = io.EOF
			}
			if err != nil {
				c.onError(err)
				return
			}
			n, err = c.encryptor.Encrypt(clientBuffer[:n], serverBuffer)
			if err != nil {
				c.onError(err)
				return
			}
			// _, err = c.s.Write(serverBuffer[:n])
			err = packet.WriteStream(c.s, serverBuffer[:n])
			if err != nil {
				c.onError(err)
				return
			}
		}

	}()

	func() {
		defer closeConn(client)
		defer closeConn(c.s)
		for {
			err = packet.Next(c.s)
			if err != nil {
				c.onError(err)
				return
			}
			n, err = c.encryptor.Decrypt(packet.Data(), serverBuffer)
			if err != nil {
				c.onError(err)
				return
			}
			_, err = client.Write(serverBuffer[:n])
			if err != nil {
				c.onError(err)
				return
			}
		}
	}()
}

func (c *Client) onError(err error) {
	if c.ErrorHandle != nil {
		c.ErrorHandle(err)
	}
}
