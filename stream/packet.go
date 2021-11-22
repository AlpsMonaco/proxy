package stream

import (
	"io"
	"sync"
	"unsafe"
)

const PacketBytes = 2
const PacketSize = 1<<(PacketBytes*8) - 1

type Stream = io.ReadWriter

type Packet struct {
	b        []byte
	bodySize int
	fullSize int
	bufSize  int
	cursor   int
}

const (
	packetShort = iota
	packetEqual
	packetExtra
)

var packetPool sync.Pool = sync.Pool{
	New: func() interface{} {
		return &Packet{
			b: make([]byte, PacketSize),
		}
	},
}

func NewPacket() *Packet {
	return packetPool.Get().(*Packet)
}

func FreePacket(packet *Packet) {
	packetPool.Put(packet)
}

func (p *Packet) Next(stream Stream) error {
	if p.b == nil {
		p.b = make([]byte, PacketSize)
	}

	var status byte
	p.cursor = p.cursor + p.fullSize
	if p.cursor < p.bufSize {
		status = p.parse(p.b[p.cursor:p.bufSize])
		if status != packetShort {
			return nil
		}
		copy(p.b, p.b[p.cursor:p.bufSize])
		p.bufSize = p.bufSize - p.cursor
		p.cursor = 0
	} else {
		p.cursor = 0
		p.bufSize = 0
	}

	var n int
	var err error
	for {
		n, err = stream.Read(p.b[p.bufSize:])
		if err != nil {
			return err
		}
		p.bufSize += n
		status = p.parse(p.b[:p.bufSize])
		if status == packetShort {
			continue
		}
		return nil
	}
}

func (p *Packet) Data() []byte {
	return p.b[p.cursor+2 : p.cursor+p.fullSize]
}

func (p *Packet) BodySize() int {
	return p.bodySize
}

func (p *Packet) FullSize() int {
	return p.fullSize
}

func (p *Packet) parse(b []byte) byte {
	if len(b) < PacketBytes {
		return packetShort
	}
	p.bodySize = int(b[0]) + int(b[1])<<8
	p.fullSize = p.bodySize + 2
	if len(b) < p.fullSize {
		return packetShort
	}
	if len(b) > p.fullSize {
		return packetExtra
	}
	return packetEqual
}

func (p *Packet) WriteStream(stream Stream, b []byte) error {
	size := len(b)
	var err error
	_, err = stream.Write((*(*[2]byte)(unsafe.Pointer(&size)))[:2])
	if err != nil {
		return err
	}
	_, err = stream.Write(b)
	return err
}
