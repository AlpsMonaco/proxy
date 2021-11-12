package vpn

import (
	"io"
	"unsafe"

	"github.com/AlpsMonaco/proxy/util"
)

const headerSize = 2

type Header struct {
	Size uint16
}

type Packet struct {
	a      *util.Allocator
	header Header
	offset int
}

func (p *Packet) Init() {
	InitPacket(p)
}

func (p *Packet) Free() {
	FreePacket(p)
}

func InitPacket(p *Packet) {
	p.a = util.GetAlloctor(1 << 16)
}

func FreePacket(p *Packet) {
	util.FreeAllocator(p.a)
}

func (p *Packet) Next(r io.Reader) error {
	return p.readFull(r)
}

func (p *Packet) GetData() []byte {
	if p.header.Size < headerSize {
		return nil
	}
	return p.a.GetBytes()[headerSize:p.header.Size]
}

func (p *Packet) GetPointer() unsafe.Pointer {
	return p.a.GetPointerN(headerSize)
}

func (p *Packet) WriteBuffer(w io.Writer, size int) error {
	p.a.GetBytes()[0] = byte(size&0x00FF) + 2
	p.a.GetBytes()[1] = byte((size & 0xFF00) >> 8)
	_, err := w.Write(p.a.GetByteSize(size + headerSize))
	return err
}

func (p *Packet) SetBuffer(b []byte) {
	copy(p.a.Shift(headerSize), b)
}

func (p *Packet) WriteSize(w io.Writer, size int) error {
	size += 2
	p.a.GetBytes()[0] = byte(size & 0x00FF)
	p.a.GetBytes()[1] = byte((size & 0xFF00) >> 8)
	_, err := w.Write(p.a.GetByteSize(headerSize))
	return err
}

// read a full packet.
func (p *Packet) readFull(r io.Reader) error {
	var err error
	if err = p.readHeader(r); err != nil {
		return err
	}
	if err = p.readBody(r); err != nil {
		return err
	}

	return nil
}

// provide full buffer everytime.
func (p *Packet) readHeader(r io.Reader) error {
	var n int
	n = p.offset - int(p.header.Size)
	if n > 0 {
		copy(p.a.GetBytes(), p.a.GetBytes()[p.header.Size:p.offset])
		p.offset = n
		if n >= headerSize {
			p.header.Size = uint16(p.a.GetBytes()[1]) << 8
			p.header.Size += uint16(p.a.GetBytes()[0])
			return nil
		}
	} else {
		p.offset = 0
	}

	var err error
	for {
		n, err = r.Read(p.a.Shift(p.offset))
		if n == 0 {
			err = io.EOF
		}
		if err != nil {
			return err
		}
		p.offset += n
		if p.offset < headerSize {
			continue
		}
		p.header.Size = uint16(p.a.GetBytes()[1]) << 8
		p.header.Size += uint16(p.a.GetBytes()[0])
		return nil
	}
}

func (p *Packet) readBody(r io.Reader) error {
	if p.offset >= int(p.header.Size) {
		return nil
	}

	var n int
	var err error
	for {
		n, err = r.Read(p.a.Shift(p.offset))
		if n == 0 {
			err = io.EOF
		}
		if err != nil {
			return err
		}
		p.offset += n
		if p.offset < int(p.header.Size) {
			continue
		}
		return nil
	}
}
