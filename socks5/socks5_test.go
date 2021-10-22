package socks5

import (
	"fmt"
	"io"
	"net"
	"runtime"
	"strconv"
	"testing"
	"time"
)

var testptr *testing.T

func assert(err error) {
	if err != nil {
		testptr.Fatal(err)
	}
}

func assertPointer(t *testing.T) {
	testptr = t
}

func TestConnection(t *testing.T) {
	assertPointer(t)

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

func TestReceive(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			t.Log(err)
		}
	}()
	assertPointer(t)

	var listener net.Listener
	var err error
	listener, err = net.Listen("tcp", "127.0.0.1:25250")
	assert(err)
	t.Log(1)

	for {
		conn, err := listener.Accept()
		assert(err)

		func(conn net.Conn) {
			defer func() {
				err := recover()
				if err != nil {
					t.Log(err)
				}
			}()

			var buf []byte = make([]byte, 1024)
			_, err = conn.Read(buf)
			assert(err)
			t.Log(buf)
			runtime.GC()

			verMsg := DecodeVersionMessage(&buf)
			if verMsg.Ver != SOCKS5_VERSION {
				t.Fatal("verMsg.Ver != SOCKS5_VERSION")
			}
			t.Log(verMsg)

			var sMsg Socks5_SelectionMessage
			sMsg.Ver = SOCKS5_VERSION
			sMsg.Method = 0x00

			bPtr := EncodeSelectionMessage(&sMsg)
			_, err = conn.Write(*bPtr)
			assert(err)

			_, err = conn.Read(buf)
			assert(err)
			reqMsg := DecodeRequestMessage(&buf)
			t.Log(reqMsg.GetHost())
			t.Log(reqMsg.GetPort())

			var respMsg = Socks5_ResponseMessage{
				Ver:   SOCKS5_VERSION,
				Rep:   SOCKS5_REP_SUCCESS,
				Rsv:   0,
				Atype: SOCKS5_ATYPE_IPV4,
				va:    [256]byte{127, 0, 0, 1, 0x62, 0xa2},
			}

			bPtr = EncodeResponseMessage(&respMsg)
			_, err = conn.Write(*bPtr)
			assert(err)

			host, err := net.DialTimeout("tcp", reqMsg.GetHost()+":"+strconv.Itoa(reqMsg.GetPort()), 1*time.Second)
			if err != nil {
				return
			}

			go func() {
				for {
					_, err = io.Copy(conn, host)
					if err != nil {
						assert(err)
					}
				}
			}()

			go func() {
				for {
					_, err = io.Copy(host, conn)
					if err != nil {
						assert(err)
					}
				}
			}()
		}(conn)
	}

	time.Sleep(10 * time.Second)

}
