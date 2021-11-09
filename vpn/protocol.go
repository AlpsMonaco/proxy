package vpn

const (
	Command_Response byte = iota
	Command_End
)

const (
	Success byte = iota
	Failed
)

type Protocol_Response struct {
	Code    byte
	MsgSize byte
	Msg     [62]byte
}

func (r *Protocol_Response) FillMsg(s string) {
	for i, v := range []byte(s) {
		r.Msg[i] = v
		r.MsgSize++
	}
}

func (r *Protocol_Response) GetSize() int {
	return int(r.MsgSize + 1)
}
