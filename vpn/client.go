package vpn

import (
	"errors"
	"fmt"
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

func (c *Client) dial() (err error) {
	c.s, err = net.Dial("tcp", fmt.Sprintf("%s:%d", c.IP, c.Port))
	return err
}

func (c *Client) Conn() net.Conn {
	return c.s
}

func (c *Client) Connect(host string, port int) error {
	// var n int
	var err error
	if c.s == nil {
		err = c.dial()
		if err != nil {
			return err
		}
	}
	c.encryptor = GetEncryptor(c.Cipher, c.Key)
	a := util.GetAlloctor(stream.PacketSize)
	defer util.FreeAllocator(a)
	sc := NewSecureConn(c.s, c.encryptor, a.GetBytes())

	(*Verify)(a.GetPointer()).SetData(10, []byte("0123456789"))
	_, err = sc.Write(a.GetBytes()[:11])
	if err != nil {
		return err
	}

	_, err = sc.Read(a.GetBytes())
	if err != nil {
		return err
	}
	gr := (*GeneralResponse)(a.GetPointer())
	if gr.Code != Code_Success {
		return errors.New(gr.Get())
	}

	(*ProxyRequest)(a.GetPointer()).SetRemoteInfo(host, port)
	_, err = sc.Write(a.GetByteSize(len(host) + 2))
	if err != nil {
		return err
	}

	_, err = sc.Read(a.GetBytes())
	if err != nil {
		return err
	}
	gr = (*GeneralResponse)(a.GetPointer())
	if gr.Code != Code_Success {
		return errors.New(gr.Get())
	}

	return nil
}

// import (
// 	"errors"
// 	"fmt"
// 	"net"
// )

// type Client struct {
// 	ServerIP    string
// 	ServerPort  int
// 	Key         []byte
// 	ErrorHandle func(err error)
// 	conn        net.Conn
// 	p           *Packet
// }

// func (c *Client) dial() (err error) {
// 	c.conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", c.ServerIP, c.ServerPort))
// 	if err != nil {
// 		return err
// 	}

// 	var p Packet
// 	c.p = &p
// 	p.Init()
// 	return nil
// }

func (c *Client) onError(err error) {
	if c.ErrorHandle != nil {
		c.ErrorHandle(err)
	}
}

// func (c *Client) Read() (b []byte, err error) {
// 	err = c.p.Next(c.conn)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return c.p.GetData(), nil
// }

// func (c *Client) GetConn() net.Conn {
// 	return c.conn
// }

// func (c *Client) Write(b []byte) (err error) {
// 	if err = c.p.WriteSize(c.conn, len(b)); err != nil {
// 		return err
// 	}

// 	_, err = c.conn.Write(b)
// 	return err
// 	// var buffer = make([]byte, len(b)+2)
// 	// buffer[0] = byte(len(b)&0x00FF) + 2
// 	// buffer[1] = byte((len(b) & 0xFF00) >> 8)
// 	// copy(buffer[2:], b)

// 	// _, err = c.conn.Write(buffer)
// 	// return err
// }

// func (c *Client) Connect(ip string, port int) error {
// 	var err error
// 	if err = c.dial(); err != nil {
// 		return err
// 	}

// 	v := (*Verify)(c.p.GetPointer())
// 	v.SetKey()
// 	if err = c.p.WriteBuffer(c.conn, 16); err != nil {
// 		return err
// 	}

// 	pr := (*ProxyRequest)(c.p.GetPointer())
// 	pr.SetInfo(ip, port)
// 	err = c.p.WriteBuffer(c.conn, pr.GetSize())
// 	if err != nil {
// 		return err
// 	}

// 	if err = c.p.Next(c.conn); err != nil {
// 		return err
// 	}

// 	gr := (*GeneralResponse)(c.p.GetPointer())
// 	if gr.Code != Success {
// 		return errors.New(string(gr.Msg[:gr.MsgSize]))
// 	}

// 	return nil
// }
