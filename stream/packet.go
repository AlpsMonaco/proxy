package stream

import (
	"io"
	"net"
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
	net.Conn

	Header  Header
	Body    []byte
	bufSize uint16

	i      int
	status byte
	h      [2]byte

	data uintptr
	len  int
	cap  int
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
	p.Body = *(*[]byte)(unsafe.Pointer(&p.data))

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

func (p *Packet) Read(b []byte) (n int, err error) {
	p.i = 0
	if p.bufSize > p.Header.Size {
		old := p.ExtraPacket()
		copy(b, old)
		p.i = len(old)
		p.status = p.Parse(b[:len(old)])
		if p.status != PacketShort {
			return p.i, nil
		}
	} else {
		p.bufSize = 0
	}
	for {
		n, err = p.Conn.Read(b[p.bufSize:])
		if n == 0 {
			err = io.EOF
		}
		if err != nil {
			return
		}
		p.i += n
		p.status = p.Parse(b[:p.i])
		if p.status == PacketShort {
			continue
		}
		return p.i, nil
	}
}

func (p *Packet) Write(b []byte) (n int, err error) {
	p.i = len(b)
	p.h[0] = byte(p.i&0x00FF) + 2
	p.h[1] = byte((p.i & 0xFF00) >> 8)

	n, err = p.Conn.Write(*(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&p.h[0])),
		Len:  2,
		Cap:  2,
	})))
	if err != nil {
		return
	}
	return p.Conn.Write(b)
}
