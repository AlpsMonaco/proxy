package socks5

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"runtime"
	"testing"
	"time"
	"unsafe"

	"github.com/AlpsMonaco/proxy/forward"
	"github.com/AlpsMonaco/proxy/util"

	_ "net/http/pprof"
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
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.LittleEndian, uint16(80))
	t.Log(buffer.Bytes())
	assertPointer(t)

	var conn net.Conn
	var err error
	// var buf = make([]byte, 256)

	conn, err = net.Dial("tcp", "127.0.0.1:7890")
	assert(err)

	var a util.Allocator
	a.Alloc(264)

	vMsg := (*Socks5_VersionMessage)(a.GetPointer())
	vMsg.Ver = SOCKS5_VERSION
	vMsg.NumMethod = 0x01
	vMsg.va[0] = SOCKS5_METHOD_NO_AUTH
	_, err = conn.Write(a.GetByteSize(3))
	assert(err)
	_, err = conn.Read(a.GetBytes())
	assert(err)

	sMsg := (*Socks5_SelectionMessage)(a.GetPointer())
	if sMsg.Ver != SOCKS5_VERSION || sMsg.Method != SOCKS5_METHOD_NO_AUTH {
		assert(errors.New("sMsg.Ver != SOCKS5_VERSION || sMsg.Method != SOCKS5_METHOD_NO_AUTH{"))
	}

	reqMsg := (*Socks5_RequestMessage)(a.GetPointer())
	reqMsg.Ver = SOCKS5_VERSION
	reqMsg.Cmd = SOCKS5_CMD_CONNECT
	reqMsg.Rsv = 0
	reqMsg.Atype = SOCKS5_ATYPE_DOMAIN

	domain := "www.google.com"
	domainSize := len(domain)
	reqMsg.va[0] = byte(domainSize)
	for i := 1; i <= domainSize; i++ {
		reqMsg.va[i] = domain[i-1]
	}

	reqMsg.va[domainSize+1] = 0
	reqMsg.va[domainSize+2] = 80

	_, err = conn.Write(a.GetByteSize(4 + domainSize + 1 + 2))
	assert(err)

	_, err = conn.Read(a.GetBytes())
	assert(err)

	// 7890
	// ox1E oxD2
	// 30 210

	respMsg := (*Socks5_ResponseMessage)(a.GetPointer())
	t.Log(respMsg)

	conn.Write([]byte(`GET / HTTP/1.1
HOST: www.google.com

`))

	for {
		n, err := conn.Read(a.GetBytes())
		fmt.Print(string(a.GetBytes()))
		if n < 256 {
			break
		}
		assert(err)
	}

	t.Log(a.GetBytes())
}

func gc() {
	for {
		time.Sleep(10 * time.Second)
		runtime.GC()
	}
}

func socks5ServerSide(t *testing.T) {
	go http.ListenAndServe(":8888", nil)

	const port = 7899
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	assert(err)

	go gc()

	for {
		conn, err := l.Accept()
		assert(err)

		go func(conn net.Conn) {
			var a util.Allocator
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
		// fmt.Println(err)
		return
	}

	go forward.NewForward(dst, conn, func(e error) { fmt.Println(e) }).Start()

}

func TestReceive(t *testing.T) {
	assertPointer(t)

	socks5ServerSide(t)
}

func TestByte(t *testing.T) {
	a := make([]byte, 10)
	a[0] = 1
	a[1] = 1
	a[2] = 2
	a[3] = 3

	b := a[0:2]
	t.Logf("0x%08x\n", uintptr(unsafe.Pointer(&a[0])))
	t.Logf("0x%08x\n", uintptr(unsafe.Pointer(&b[0])))

	t.Log(*(*reflect.SliceHeader)(unsafe.Pointer(&a)))
	t.Log(*(*reflect.SliceHeader)(unsafe.Pointer(&b)))

}
