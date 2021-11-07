package stream

import (
	"reflect"
	"sync"
	"unsafe"
)

const PacketSize = 1 << 15

var (
	headerSize uint16 = uint16(unsafe.Sizeof(Header{}))
	bodySize          = PacketSize - headerSize
)

func GetHeaderSize() uint16 {
	return headerSize
}

const (
	PacketEqual byte = 0x00
	PacketShort byte = 0x01
	PacketExtra byte = 0x02
)

type Header struct {
	Size uint16
}

type Packet struct {
	Header  Header
	Body    *[]byte
	bufSize uint16
	data    uintptr
	len     int
	cap     int
}

var packetPool sync.Pool = sync.Pool{
	New: func() interface{} {
		return &Packet{}
	},
}

func GetPacket() *Packet {
	return (packetPool.Get()).(*Packet)
}

func FreePacket(p *Packet) {
	packetPool.Put(p)
}

func (p *Packet) Parse(b []byte) byte {
	p.bufSize = uint16(len(b))
	if p.bufSize < headerSize {
		return PacketShort
	}
	p.Header = *(*Header)(unsafe.Pointer(&b[0]))
	if p.bufSize < p.Header.Size {
		return PacketShort
	}
	p.data = uintptr(unsafe.Pointer(&b[headerSize]))
	p.len = int(p.Header.Size - headerSize)
	p.cap = p.len
	p.Body = (*[]byte)(unsafe.Pointer(&p.data))

	if p.bufSize > p.Header.Size {
		return PacketExtra
	}
	return PacketEqual
}

func (p *Packet) ExtraPacket() []byte {
	if p.bufSize <= p.Header.Size {
		return nil
	}
	var newBufSize int = int(p.bufSize - p.Header.Size)

	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: p.data + uintptr(p.Header.Size-headerSize),
		Len:  newBufSize,
		Cap:  newBufSize,
	}))
}
