package stream

import "net"

const PackageBytes = 2
const PackageSize = 1<<(PackageBytes*8) - 1
const extendSize = 256

const (
	PacketShort byte = iota
	PacketEqual
	PacketExtra
)

type Packet struct {
	buf        []byte
	offset     int
	size       int
	laststatus byte
}

func (p *Packet) Encode(size int, b []byte) []byte {
	if len(p.buf) < len(b)+PackageBytes {
		p.Extend(len(b) + PackageBytes)
	}
	p.buf[1] = byte((size & 0xFF00) >> 8)
	p.buf[0] = byte(size & 0x00FF)
	copy(p.buf[PackageBytes:], b)
	return p.buf[PackageBytes : size+PackageBytes]
}

func (p *Packet) Parse(b []byte) byte {
	totalSize := len(b) + p.offset
	if len(p.buf) < totalSize {
		p.Extend(totalSize)
	}
	copy(p.buf[p.offset:], b)
	p.offset = totalSize
	if p.offset < PackageBytes {
		return PacketShort
	}
	p.size = int(p.buf[1])<<8 + int(p.buf[0])
	totalSize = p.size + PackageBytes
	if p.offset < totalSize {
		p.laststatus = PacketShort
	} else if p.offset == totalSize {
		p.laststatus = PacketEqual
	} else {
		p.laststatus = PacketExtra
	}
	return p.laststatus
}

func (p *Packet) Data() []byte {
	return p.buf[PackageBytes : p.size+PackageBytes]
}

func (p *Packet) Sort() {
	if p.laststatus == PacketShort {
		return
	} else if p.laststatus == PacketEqual {
		p.offset = 0
		return
	} else {
		copy(p.buf, p.buf[p.size+PackageBytes:p.offset])
		p.offset = p.offset - p.size - PackageBytes
	}
}

func (p *Packet) Extend(size int) {
	newBuf := make([]byte, size)
	copy(newBuf, p.buf)
	p.buf = newBuf
}

type TinyBuffer struct {
	net.Conn

	b   []byte
	old []byte
}

func (tb *TinyBuffer) Init() {
	tb.b = make([]byte, 0xFF)
	tb.old = make([]byte, 0xFF)
}

func (tb *TinyBuffer) GetBuffer() []byte {
	return tb.b[1:]
}

func (tb *TinyBuffer) GetSize() byte {
	return tb.b[0]
}

func (tb *TinyBuffer) SetSize(size byte) {
	tb.b[0] = size
}

func (tb *TinyBuffer) Write(b []byte) (int, error) {
	if len(b) > 0xFF {
		panic("Stack Overflow")
	}
	tb.SetSize(byte(len(b)))
	copy(tb.GetBuffer(), b)
	_, err := tb.Write(tb.b[:tb.GetSize()])
	return len(b), err
}

func (tb *TinyBuffer) Read(b []byte) (n int, err error) {
	if len(b) > 0xFF {
		panic("Stack Overflow")
	}
	for {
		n, err = tb.Conn.Read(tb.GetBuffer())
		if err != nil {
			return n, err
		}
	}
}
