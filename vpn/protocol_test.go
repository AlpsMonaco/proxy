package vpn

import (
	"fmt"
	"io"
	"reflect"
	"testing"
	"unsafe"
)

func TestPacketParser(t *testing.T) {
	var msg Message
	msg.Header.Size = 0x0008
	msg.Header.Cmd = 0x0000
	var size int = 2

	msg.Body[0] = 0
	msg.Body[1] = 1
	msg.Body[2] = 2
	msg.Body[3] = 3
	msg.Body[4] = 4

	b := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&msg)),
		Len:  size,
		Cap:  size,
	}))

	var p Packet
	status, extra := p.Parse(b)
	fmt.Println(status)
	fmt.Println(extra)
	fmt.Println(p.Header)
	fmt.Println(p.Body)

}

func TestPacketTransport(t *testing.T) {
	var msg Message
	var size int = 10

	b := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&msg)),
		Len:  size,
		Cap:  size,
	}))

	msg.Header.Size = uint16(size)
	msg.Header.Cmd = 0
	msg.Body[0] = 0x0A
	msg.Body[1] = 0x0A
	msg.Body[2] = 0x0A
	msg.Body[3] = 0x0A
	msg.Body[4] = 0x0A
	msg.Body[5] = 0x0A
	fmt.Println(b)

	size = 6
	b = *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&msg)),
		Len:  size,
		Cap:  size,
	}))

	msg.Header.Size = uint16(size)
	msg.Header.Cmd = 0
	msg.Body[0] = 0x06
	msg.Body[1] = 0x06
	fmt.Println(b)

	size = 5
	b = *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&msg)),
		Len:  size,
		Cap:  size,
	}))

	msg.Header.Size = uint16(size)
	msg.Header.Cmd = 0
	msg.Body[0] = 0x05
	fmt.Println(b)

	size = 17
	b = *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&msg)),
		Len:  size,
		Cap:  size,
	}))

	msg.Header.Size = 17
	msg.Header.Cmd = 0
	msg.Body[0] = 17
	msg.Body[1] = 17
	msg.Body[2] = 17
	msg.Body[3] = 17
	msg.Body[4] = 17
	msg.Body[5] = 17
	msg.Body[6] = 17
	msg.Body[7] = 17
	msg.Body[8] = 17
	msg.Body[9] = 17
	msg.Body[10] = 17
	msg.Body[11] = 17
	msg.Body[12] = 17
	fmt.Println(b)

}

type FakeNet struct {
	io.ReadWriter
	b []byte
	t int
}

func (fn *FakeNet) Write(b []byte) (int, error) {
	return 0, nil
}

func (fn *FakeNet) Read(b []byte) (int, error) {
	fn.t = fn.t + 1
	switch fn.t {
	case 1:
		copy(b, fn.b[0:5])
		return 5, nil
	case 2:
		copy(b, fn.b[5:12])
		return 12 - 5, nil
	case 3:
		copy(b, fn.b[12:25])
		return 25 - 12, nil
	case 4:
		copy(b, fn.b[25:38])
		return 38 - 25, nil
	default:
		panic("end")
	}
}

func TestPacketSplit(t *testing.T) {
	var buf []byte = make([]byte, 1024)
	var b = []byte{10, 0, 0, 0, 10, 10, 10, 10, 10, 10, 6, 0, 0, 0, 6, 6, 5, 0, 0, 0, 5, 17, 0, 0, 0, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17}
	var fn FakeNet
	fn.b = b
	var beginIndex int = 0

	for {
		var p Packet
		n, _ := fn.Read(buf[beginIndex:])
		status, extra := p.Parse(buf[:beginIndex+n])
		if status == PacketShort {
			beginIndex += n
			continue
		} else if status == PacketExtra {
			copy(buf, extra)
			beginIndex = len(extra)
		} else {
			beginIndex = 0
		}

		fmt.Println(p.Header)
		fmt.Println(p.Body)
	}
}