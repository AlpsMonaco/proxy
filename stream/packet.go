package stream

import (
	"io"

	"github.com/AlpsMonaco/proxy/util"
)

const PacketBytes = 2
const PacketSize = 1<<(PacketBytes*8) - 1

type stream = io.ReadWriteCloser

type Packet struct {
	stream
	a       *util.Allocator
	bufSize int
}

const (
	packetShort byte = iota
	packetEqual
	packetExtra
)

func (p *Packet) Next() error {
	if p.a == nil {
		util.GetAlloctor(PacketSize)
	}
	var n int
	var err error
	var status byte
	for {
		n, err = p.stream.Read(p.a.GetBytes())
		if err != nil {
			return err
		}
		if n == 0 {
			return io.EOF
		}
		p.bufSize += n
		status = p.parse(p.a.GetByteSize(p.bufSize))
		if status == packetShort {
			continue
		}
	}

	return nil
}

func (p *Packet) accept(b []byte) {

}

func (p *Packet) parse(b []byte) byte {
	if len(b) < PacketBytes {
		return packetShort
	}
	size := int(b[0]) + int(b[1])<<8
	if len(b) < size {
		return packetShort
	}
	if len(b) > size {
		return packetExtra
	}
	return packetEqual
}

// const (
// 	PacketShort byte = iota
// 	PacketEqual
// 	PacketExtra
// )

// type Packet struct {
// 	buf        []byte
// 	offset     int
// 	size       int
// 	laststatus byte
// }

// func (p *Packet) Encode(size int, b []byte) []byte {
// 	if len(p.buf) < len(b)+PacketBytes {
// 		p.Extend(len(b) + PacketBytes)
// 	}
// 	p.buf[1] = byte((size & 0xFF00) >> 8)
// 	p.buf[0] = byte(size & 0x00FF)
// 	copy(p.buf[PacketBytes:], b)
// 	return p.buf[PacketBytes : size+PacketBytes]
// }

// func (p *Packet) Parse(b []byte) byte {
// 	totalSize := len(b) + p.offset
// 	if len(p.buf) < totalSize {
// 		p.Extend(totalSize)
// 	}
// 	copy(p.buf[p.offset:], b)
// 	p.offset = totalSize
// 	if p.offset < PacketBytes {
// 		return PacketShort
// 	}
// 	p.size = int(p.buf[1])<<8 + int(p.buf[0])
// 	totalSize = p.size + PacketBytes
// 	if p.offset < totalSize {
// 		p.laststatus = PacketShort
// 	} else if p.offset == totalSize {
// 		p.laststatus = PacketEqual
// 	} else {
// 		p.laststatus = PacketExtra
// 	}
// 	return p.laststatus
// }

// func (p *Packet) Data() []byte {
// 	return p.buf[PacketBytes : p.size+PacketBytes]
// }

// func (p *Packet) Sort() {
// 	if p.laststatus == PacketShort {
// 		return
// 	} else if p.laststatus == PacketEqual {
// 		p.offset = 0
// 		return
// 	} else {
// 		copy(p.buf, p.buf[p.size+PacketBytes:p.offset])
// 		p.offset = p.offset - p.size - PacketBytes
// 	}
// }

// func (p *Packet) Extend(size int) {
// 	newBuf := make([]byte, size)
// 	copy(newBuf, p.buf)
// 	p.buf = newBuf
// }

// type TinyBuffer struct {
// 	net.Conn

// 	b   []byte
// 	old []byte
// }

// func (tb *TinyBuffer) Init() {
// 	tb.b = make([]byte, 0xFF)
// 	tb.old = make([]byte, 0xFF)
// }

// func (tb *TinyBuffer) GetBuffer() []byte {
// 	return tb.b[1:]
// }

// func (tb *TinyBuffer) GetSize() byte {
// 	return tb.b[0]
// }

// func (tb *TinyBuffer) SetSize(size byte) {
// 	tb.b[0] = size
// }

// func (tb *TinyBuffer) Write(b []byte) (int, error) {
// 	if len(b) > 0xFF {
// 		panic("Stack Overflow")
// 	}
// 	tb.SetSize(byte(len(b)))
// 	copy(tb.GetBuffer(), b)
// 	_, err := tb.Write(tb.b[:tb.GetSize()])
// 	return len(b), err
// }

// func (tb *TinyBuffer) Read(b []byte) (n int, err error) {
// 	if len(b) > 0xFF {
// 		panic("Stack Overflow")
// 	}
// 	for {
// 		n, err = tb.Conn.Read(tb.GetBuffer())
// 		if err != nil {
// 			return n, err
// 		}
// 	}
// }
