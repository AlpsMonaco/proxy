package vpn

import (
	"io"
	"unsafe"
)

type Packet interface {
	PacketSender
	PacketReceiver
}

type PacketSender interface {
	Send(io.Writer, []byte) error
}

type PacketReceiver interface {
	Next(io.Reader, []byte) error
	Data() []byte
}

type Raw struct{ b []byte }

func (r *Raw) Next(reader io.Reader, b []byte) error {
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

func (p *Raw) Data() []byte { return p.b }

func (p *Raw) Send(writer io.Writer, b []byte) (err error) {
	_, err = writer.Write(b)
	return err
}

const PacketBytes = 2
const PacketSize = 1<<(PacketBytes*8) - 1

const (
	packetShort = iota
	packetEqual
	packetExtra
)

type SizePacket struct {
	data     []byte
	bodysize int
	fullsize int
	bufsize  int
	cursor   int
}

func (sp *SizePacket) Next(r io.Reader, buf []byte) error {
	var status byte
	sp.cursor = sp.cursor + sp.fullsize
	if sp.cursor < sp.bufsize {
		status = sp.parse(buf[sp.cursor:sp.bufsize])
		if status != packetShort {
			sp.data = buf[sp.cursor+2 : sp.cursor+sp.fullsize]
			return nil
		}
		copy(buf, buf[sp.cursor:sp.bufsize])
		sp.bufsize = sp.bufsize - sp.cursor
		sp.cursor = 0
	} else {
		sp.cursor = 0
		sp.bufsize = 0
	}

	var n int
	var err error
	for {
		n, err = r.Read(buf[sp.bufsize:])
		if n == 0 && err == nil {
			err = io.EOF
		}
		if err != nil {
			return err
		}
		sp.bufsize += n
		status = sp.parse(buf[:sp.bufsize])
		if status == packetShort {
			continue
		}
		sp.data = buf[sp.cursor+2 : sp.cursor+sp.fullsize]
		return nil
	}
}

func (p *SizePacket) Data() []byte { return p.data }

func (p *SizePacket) BodySize() int { return p.bodysize }

func (p *SizePacket) FullSize() int { return p.fullsize }

func (p *SizePacket) parse(b []byte) byte {
	if len(b) < PacketBytes {
		return packetShort
	}
	p.bodysize = int(b[0]) + int(b[1])<<8
	p.fullsize = p.bodysize + PacketBytes
	if len(b) < p.fullsize {
		return packetShort
	}
	if len(b) > p.fullsize {
		return packetExtra
	}
	return packetEqual
}

func (p *SizePacket) Send(writer io.Writer, b []byte) error {
	size := len(b)
	var err error
	_, err = writer.Write((*(*[2]byte)(unsafe.Pointer(&size)))[:2])
	if err != nil {
		return err
	}
	_, err = writer.Write(b)
	return err
}
