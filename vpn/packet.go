package vpn

import (
	"io"
	"unsafe"
)

type PacketProtocol interface {
	PacketSpliter
	Package
}

type Package interface {
	Pack(io.Writer, []byte) error
}

type PacketSpliter interface {
	Next(io.Reader, []byte) error
	Data() []byte
}

type raw struct {
	b []byte
}

func (r *raw) Next(reader io.Reader, b []byte) error {
	n, err := reader.Read(b)
	if n == 0 && err == nil {
		err = io.EOF
	}
	if err != nil {
		return err
	}
	r.b = b[:n]
	return nil
}

func (r *raw) Data() []byte {
	return r.b
}

func (r *raw) Pack(writer io.Writer, b []byte) (err error) {
	_, err = writer.Write(b)
	return err
}

const PacketBytes = 2
const PacketSize = 1<<(PacketBytes*8) - 1

type Packet struct {
	data     []byte
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

func (p *Packet) Next(r io.Reader, buf []byte) error {
	var status byte
	p.cursor = p.cursor + p.fullSize
	if p.cursor < p.bufSize {
		status = p.parse(buf[p.cursor:p.bufSize])
		if status != packetShort {
			p.data = buf[p.cursor+2 : p.cursor+p.fullSize]
			return nil
		}
		copy(buf, buf[p.cursor:p.bufSize])
		p.bufSize = p.bufSize - p.cursor
		p.cursor = 0
	} else {
		p.cursor = 0
		p.bufSize = 0
	}

	var n int
	var err error
	for {
		n, err = r.Read(buf[p.bufSize:])
		if n == 0 && err == nil {
			err = io.EOF
		}
		if err != nil {
			return err
		}
		p.bufSize += n
		status = p.parse(buf[:p.bufSize])
		if status == packetShort {
			continue
		}
		p.data = buf[p.cursor+2 : p.cursor+p.fullSize]
		return nil
	}
}

func (p *Packet) Data() []byte {
	return p.data
	// return buf[p.cursor+2 : p.cursor+p.fullSize]
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
	p.fullSize = p.bodySize + PacketBytes
	if len(b) < p.fullSize {
		return packetShort
	}
	if len(b) > p.fullSize {
		return packetExtra
	}
	return packetEqual
}

func (p *Packet) Pack(writer io.Writer, b []byte) error {
	size := len(b)
	var err error
	_, err = writer.Write((*(*[2]byte)(unsafe.Pointer(&size)))[:2])
	if err != nil {
		return err
	}
	_, err = writer.Write(b)
	return err
}
