package socks5

import (
	"errors"
	"fmt"
	"io"
	"net"
	"runtime"
	"testing"
	"time"

	"github.com/AlpsMonaco/proxy/util"
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

func gc() {
	for {
		time.Sleep(10 * time.Second)
		runtime.GC()
	}
}

func socks5ServerSide(t *testing.T) {
	const port = 7899
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	assert(err)

	go gc()

	for {
		conn, err := l.Accept()
		assert(err)

		go func(conn net.Conn) {
			var a util.Alloctor
			a.Alloc(264)

			_, err := conn.Read(a.GetBytes())
			assert(err)
			vMsg := (*Socks5_VersionMessage)(a.GetPointer())
			if vMsg.Ver != SOCKS5_VERSION {
				assert(errors.New("vMsg.Ver != SOCKS5_VERSION"))
			}
			// t.Log(vMsg)

			sMsg := (*Socks5_SelectionMessage)(a.GetPointer())
			sMsg.Method = 0x00
			sMsg.Ver = SOCKS5_VERSION
			_, err = conn.Write(a.GetByteSize(2))
			// t.Log(a.GetBytes())
			assert(err)

			_, err = conn.Read(a.GetBytes())
			assert(err)
			reqMsg := (*Socks5_RequestMessage)(a.GetPointer())
			addr := reqMsg.GetHost()
			port := reqMsg.GetPort()

			respMsg := (*Socks5_ResponseMessage)(a.GetPointer())
			respMsg.Ver = 0x05
			respMsg.Rep = 0x00
			respMsg.Rsv = 0x00
			respMsg.Atype = 0x01
			respMsg.va[0] = 127
			respMsg.va[1] = 0
			respMsg.va[2] = 0
			respMsg.va[3] = 1
			respMsg.va[4] = 0x1E
			respMsg.va[5] = 0xDB
			_, err = conn.Write(a.GetByteSize(10))
			// t.Log(a.GetBytes())
			assert(err)
			Proxy(addr, port, conn)

		}(conn)
	}
}

func Proxy(addr string, port int, conn net.Conn) {
	dst, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		defer func() {
			dst.Close()
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		for {
			if n, err := io.Copy(dst, conn); err != nil || n == 0 {
				fmt.Println(err)
				return
			}
		}
	}()

	go func() {
		defer func() {
			conn.Close()
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()

		for {
			if n, err := io.Copy(conn, dst); err != nil || n == 0 {
				fmt.Println(err)
				return
			}
		}
	}()
}

func TestReceive(t *testing.T) {
	assertPointer(t)

	socks5ServerSide(t)
}
