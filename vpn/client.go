package vpn

import (
	"net"

	"github.com/AlpsMonaco/proxy/util"
)

type Client struct {
	ServerIP   string
	ServerPort int
	Cipher     CipherEnum
	Key        []byte
	encryptor  Encryptor
}

func (c *Client) Connect() error {
	var conn net.Conn
	var err error
	conn, err = net.Dial("tcp", util.SprintfAddress(c.ServerIP, c.ServerPort))
	if err != nil {
		return err
	}

	var allocator *util.Allocator = util.GetAlloctor(256)
	// var n int
	defer util.FreeAllocator(allocator)

	(*HelloMessage)(allocator.GetPointer()).SetMessage(masquerade)
	_, err = conn.Write((*HelloMessage)(allocator.GetPointer()).GetBytes())
	if err != nil {
		closeConn(conn)
		return err
	}

	_, err = conn.Read(allocator.GetBytes())
	if err != nil {
		closeConn(conn)
		return err
	}

	if (*Ack)(allocator.GetPointer()).GetCode() != Code_Success {
		closeConn(conn)
		return ErrServerRejected
	}

	c.encryptor = GetEncryptor(c.Cipher, c.Key)
	return nil
}
