package stream

import (
	"fmt"
	"io"
	"reflect"
	"testing"
	"unsafe"

	"github.com/AlpsMonaco/proxy/util"
)

func TestPacketParser(t *testing.T) {
	var msg MessageTest
	msg.Header.Size = 0x0008
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

	var p PacketTest
	status, extra := p.Parse(b)
	fmt.Println(status)
	fmt.Println(extra)
	fmt.Println(p.Header)
	fmt.Println(p.Body)

}

func TestPacketTransport(t *testing.T) {
	var result = []byte{}
	var msg MessageTest
	var size int = 10

	b := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&msg)),
		Len:  size,
		Cap:  size,
	}))

	msg.Header.Size = uint16(size)
	for i := 0; i < int(msg.Header.Size)-1; i++ {
		msg.Body[i] = byte(size)
	}
	fmt.Println(b)
	result = append(result, b...)

	size = 6
	b = *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&msg)),
		Len:  size,
		Cap:  size,
	}))

	msg.Header.Size = uint16(size)
	for i := 0; i < int(msg.Header.Size)-1; i++ {
		msg.Body[i] = byte(size)
	}
	fmt.Println(b)
	result = append(result, b...)

	size = 5
	b = *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&msg)),
		Len:  size,
		Cap:  size,
	}))

	msg.Header.Size = uint16(size)
	for i := 0; i < int(msg.Header.Size)-1; i++ {
		msg.Body[i] = byte(size)
	}
	result = append(result, b...)

	size = 17
	b = *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&msg)),
		Len:  size,
		Cap:  size,
	}))

	msg.Header.Size = uint16(size)
	for i := 0; i < int(msg.Header.Size)-1; i++ {
		msg.Body[i] = byte(size)
	}
	fmt.Println(b)
	result = append(result, b...)

	fmt.Println(len(result))
	var s string = "{"
	for i := range result {
		s = s + fmt.Sprint(result[i]) + ","
	}
	s = s[:len(s)-1] + "}"
	fmt.Println(s)
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
		copy(b, fn.b[0:3])
		return 3, nil
	case 2:
		copy(b, fn.b[3:12])
		return 12 - 3, nil
	case 3:
		copy(b, fn.b[12:21])
		return 21 - 12, nil
	case 4:
		copy(b, fn.b[21:38])
		return 38 - 21, nil
	default:
		return 0, nil
	}
}

func getPacket(p *PacketTest) {
	fmt.Println("----------Packet Start----------")
	fmt.Println(p.Header.Size)
	fmt.Println(p.Body)
	fmt.Println("----------Packet End----------")
}

type PacketTest struct {
	Header *Header
	Body   *[maxMessageSize - 4]byte
}

func (p *PacketTest) Parse(b []byte) (byte, []byte) {
	if len(b) < 4 {
		return PacketShort, nil
	}

	p.Header = (*Header)(unsafe.Pointer(&b[0]))
	if len(b) < int(p.Header.Size) {
		return PacketShort, nil
	}

	p.Body = (*[maxMessageSize - 4]byte)(unsafe.Pointer(&b[4]))
	if len(b) > int(p.Header.Size) {
		return PacketExtra, b[p.Header.Size:]
	}

	return PacketEqual, nil
}

const maxMessageSize uint16 = 1<<6 - 1

type MessageTest struct {
	Header Header
	Body   [maxMessageSize - 4]byte
}

func TestPacketSplit(t *testing.T) {
	var buf []byte = make([]byte, 1024)
	for i := range buf {
		buf[i] = 99
	}

	var b = []byte{10, 0, 0, 0, 10, 10, 10, 10, 10, 10, 6, 0, 0, 0, 6, 6, 5, 0, 0, 0, 5, 5, 0, 0, 0, 5, 6, 0, 1, 0, 17, 17, 6, 0, 2, 0, 17, 1}
	var fn FakeNet
	fn.b = b

	var p PacketTest
	var n, i int

	for {
		n, _ = fn.Read(buf[i:])
		if n == 0 {
			break
		}
		status, extra := p.Parse(buf[:i+n])
		switch status {
		case PacketShort:
			i += n
		case PacketExtra:
			getPacket(&p)
			copy(buf, extra)
			i = len(extra)
			isBreak := false
			for {
				status, extra = p.Parse(buf[:i])
				switch status {
				case PacketShort:
					isBreak = true
				case PacketExtra:
					getPacket(&p)
					copy(buf, extra)
					i = len(extra)
				case PacketEqual:
					getPacket(&p)
					isBreak = true
					i = 0
				}
				if isBreak {
					break
				}
			}
		case PacketEqual:
			getPacket(&p)
			i = 0
		}
	}

}

func PrintPacket(p *Packet) {
	fmt.Println(p.Header)
	fmt.Println((*p.Body)[:p.Header.Size-HeaderSize])
	fmt.Println(p.bufSize)
}

func TestPacketParse(t *testing.T) {
	a := util.GetAlloctor(PacketSize)
	var n, i int
	var status byte
	var err error
	var p Packet
	var b = []byte{10, 0, 10, 10, 10, 10, 10, 10, 10, 10, 6, 0, 6, 6, 6, 6, 5, 0, 5, 5, 5, 17, 0, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17}
	var fn FakeNet
	fn.b = b

	for {
		n, err = fn.Read(a.Shift(i))
		if err != nil {
			panic(err)
		}
		if n == 0 {
			break
		}
		status = p.Parse(a.GetByteSize(i + n))
		if status == PacketShort {
			i += n
			continue
		}
		PrintPacket(&p)
		if status == PacketEqual {
			i = 0
			continue
		}

		// case PacketExtra
		for {
			b := p.ExtraPacket()
			status = p.Parse(b)
			if status == PacketShort {
				copy(a.GetBytes(), b)
				i = len(b)
				break
			}
			PrintPacket(&p)
			if status == PacketEqual {
				i = 0
				break
			}
		}
	}
}
