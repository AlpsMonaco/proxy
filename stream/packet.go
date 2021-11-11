package stream

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
