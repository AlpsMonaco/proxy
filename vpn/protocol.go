package vpn

const version byte = 0x01

type HelloMessage struct {
	msgSize byte
	msg     [255]byte
}

func (hm *HelloMessage) SetMessage(s string) {
	hm.SetBytes([]byte(s))
}

func (hm *HelloMessage) SetBytes(b []byte) {
	if len(b) > 255 {
		hm.msgSize = 255
	} else {
		hm.msgSize = byte(len(b))
	}
	copy(hm.msg[:], b)
}

func (hm *HelloMessage) GetBytes() []byte {
	return hm.msg[:hm.msgSize]
}

func (hm *HelloMessage) GetMessage() string {
	return string(hm.GetBytes())
}

const (
	Code_Error byte = iota
	Code_Success
)

type Ack struct {
	code byte
}

func (a *Ack) SetCode(code byte) {
	a.code = code
}

func (a *Ack) GetCode() byte {
	return a.code
}
