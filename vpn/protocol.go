package vpn

const DataSize int = 2 ^ 10
const PacketSize = DataSize + 1

type Packet struct {
	Size byte
	Data [DataSize]byte
}
