package vpn

const (
	Command_Echo uint16 = iota
	Command_End
)

type Protocol_Echo struct {
	Va [256]byte
}
