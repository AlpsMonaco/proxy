package socks5

import (
	"fmt"
	"net"
	"testing"
)

var testptr *testing.T

func assert(err error) {
	if err != nil {
		testptr.Fatal(err)
	}
}

func SetTestptr(t *testing.T) {
	testptr = t
}

func TestConnection(t *testing.T) {
	SetTestptr(t)
	var conn net.Conn
	var err error
	var buf = make([]byte, 256)

	conn, err = net.Dial("tcp", "127.0.0.1:7890")
	assert(err)

	var vMsg Socks5_VersionMessage
	vMsg.Ver = 0x05
	vMsg.SetMethod(0x00)

	bPtr := EncodeVersionMessage(&vMsg)
	_, err = conn.Write(*bPtr)
	assert(err)

	_, err = conn.Read(buf)
	assert(err)

	sMsg := DecodeSelectionMessage(&buf)
	t.Log(sMsg)

	var reqMsg Socks5_RequestMessage = Socks5_RequestMessage{
		Ver:   SOCKS5_VERSION,
		Cmd:   SOCKS5_CMD_CONNECT,
		Rsv:   0,
		Atype: SOCKS5_ATYPE_DOMAIN,
	}

	domain := "www.google.com"
	domainSize := len(domain)
	reqMsg.va[0] = byte(domainSize)
	for i := 1; i <= domainSize; i++ {
		reqMsg.va[i] = domain[i-1]
	}

	reqMsg.va[domainSize+1] = 0
	reqMsg.va[domainSize+2] = 80

	bPtr = EncodeRequestMessage(&reqMsg)
	_, err = conn.Write(*bPtr)
	assert(err)

	_, err = conn.Read(buf)
	assert(err)

	respMsg := DecodeResponseMessage(&buf)
	t.Log(respMsg)

	conn.Write([]byte(`GET / HTTP/1.1
HOST: www.google.com

`))

	for {
		n, err := conn.Read(buf)
		fmt.Print(string(buf))
		if n < 256 {
			break
		}
		assert(err)
	}

	t.Log(string(buf))

}
