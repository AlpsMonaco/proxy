package vpn

import "unsafe"

type Header struct {
	Size uint16
	Cmd  uint16
}

const maxMessageSize uint16 = 1<<6 - 1

type Message struct {
	Header Header
	Body   [maxMessageSize - 4]byte
}

type Packet struct {
	Header *Header
	Body   *[maxMessageSize - 4]byte
}

func (p *Packet) Parse(b []byte) (byte, []byte) {
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

/*
Expect n byte
Got(buf,m byte)
Split()
*/

const (
	PacketEqual byte = iota
	PacketShort
	PacketExtra
)

// func ParseHeader(b []byte) (byte, *Header) {
// 	if len(b) < 4 {
// 		return PacketShort, nil
// 	}

// 	return PacketEqual, (*Header)(unsafe.Pointer(&b[0]))
// }

// func ParseBody(header *Header, b []byte) byte {
// 	var status byte = PacketEqual
// 	if len(b) < int(header.Size) {
// 		status = PacketShort
// 		return status
// 	}

// 	if len(b) > int(header.Size) {
// 		status = PacketExtra
// 	}

// }
