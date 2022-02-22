package vpn

import (
	"bytes"
	"net"
	"unsafe"

	"github.com/AlpsMonaco/proxy/util"
)

type Client struct {
	ServerIP   string
	ServerPort int
	Key        string
	conn       net.Conn
	proxyConn  ProxyConn
}

func (c *Client) Connect(remoteAddr string, remotePort int) (err error) {
	c.conn, err = net.Dial("tcp", util.SprintfAddress(c.ServerIP, c.ServerPort))
	if err != nil {
		return
	}
	_, err = c.conn.Write(clientPrefixReqBytes)
	if err != nil {
		return
	}
	var buffer []byte = make([]byte, allocMemSize)
	n, err := c.conn.Read(buffer)
	if err != nil {
		return
	}
	if !bytes.Equal(buffer[:n], serverPrefixRetBytes) {
		err = ErrPrefixNotMacth
		return
	}
	var encryptor Encryptor = &ChaCha20Poly1305{}
	encryptor.Key([]byte(c.Key))
	var packet Packet = &SizePacket{}
	var recvbuffer []byte = make([]byte, allocMemSize)
	c.proxyConn = CreateProxyConn(encryptor, packet, c.conn, buffer, recvbuffer)
	return c.secureConnect(remoteAddr, remotePort)
}

func (c *Client) secureConnect(remoteAddr string, remotePort int) (err error) {
	err = c.proxyConn.Send(versionBytes)
	if err != nil {
		return
	}
	// block here
	var b []byte
	b, err = c.proxyConn.Recv()
	if err != nil {
		return
	}
	if !bytes.Equal(b, versionBytes) {
		err = ErrVersionDismatch
		return
	}
	ci := (*ConnectInfo)(unsafe.Pointer(&c.proxyConn.recvbuffer[0]))
	ci.SetConnection(remoteAddr, uint16(remotePort))
	err = c.proxyConn.Send(ci.info[:ci.Size()])
	if err != nil {
		return
	}
	b, err = c.proxyConn.Recv()
	if err != nil {
		return
	}
	if !bytes.Equal(b, SuccessBytes) {
		err = ErrServerRejected
		return
	}
	return
}

func (c *Client) Proxy(conn net.Conn) {
	BeginProxy(conn, &c.proxyConn)
}
