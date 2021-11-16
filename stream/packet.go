package stream

import (
	"io"

	"github.com/AlpsMonaco/proxy/util"
)

const PacketBytes = 2
const PacketSize = 1<<(PacketBytes*8) - 1

type Stream = io.ReadWriter

type Packet struct {
	Stream
	a *util.Allocator

	BodySize int
	FullSize int
	bufSize  int
	cursor   int
}

const (
	packetShort = iota
	packetEqual
	packetExtra
)

func (p *Packet) Free() {
	if p.a != nil {
		util.FreeAllocator(p.a)
	}
}

func (p *Packet) Next() error {
	if p.a == nil {
		p.a = util.GetAlloctor(PacketSize)
	}
	var status byte
	p.cursor = p.cursor + p.FullSize
	if p.cursor < p.bufSize {
		status = p.parse(p.a.GetBytes()[p.cursor:p.bufSize])
		if status != packetShort {
			return nil
		}
		copy(p.a.GetBytes(), p.a.GetBytes()[p.cursor:p.bufSize])
		p.bufSize = p.bufSize - p.cursor
		p.cursor = 0
	} else {
		p.cursor = 0
		p.bufSize = 0
	}

	var n int
	var err error
	for {
		n, err = p.Read(p.a.Shift(p.bufSize))
		if err != nil {
			return err
		}
		p.bufSize += n
		status = p.parse(p.a.GetByteSize(p.bufSize))
		if status == packetShort {
			continue
		}
		return nil
	}
}

func (p *Packet) Data() []byte {
	return p.a.GetBytes()[p.cursor+2 : p.cursor+p.FullSize]
}

func (p *Packet) parse(b []byte) byte {
	if len(b) < PacketBytes {
		return packetShort
	}
	p.BodySize = int(b[0]) + int(b[1])<<8
	p.FullSize = p.BodySize + 2
	if len(b) < p.FullSize {
		return packetShort
	}
	if len(b) > p.FullSize {
		return packetExtra
	}
	return packetEqual
}
