package vpn

/*
protocol part of vpn.
*/

const Version byte = 0x01

const (
	Code_Success byte = iota
	Code_Error
)

type GeneralResponse struct {
	Code    byte
	MsgSize byte
	Msg     [64]byte
}

func (gr *GeneralResponse) Set(code byte, msg string) {
	gr.MsgSize = 0
	gr.Code = code
	for i := 0; i < len(msg); i++ {
		gr.MsgSize++
		gr.Msg[i] = msg[i]
	}
}

func (gr *GeneralResponse) Get() string {
	if gr.MsgSize == 0 {
		return ""
	}
	b := make([]byte, gr.MsgSize)
	var i byte
	for i = 0; i < gr.MsgSize; i++ {
		b[i] = gr.Msg[i]
	}
	return string(b)
}

type Verify struct {
	Va [256]byte
}

func (v *Verify) SetData(size byte, b []byte) {
	v.Va[12] = byte(size)
	var i byte
	for i = 0; i < size; i++ {
		v.Va[1+i] = b[i]
	}
}

func (v *Verify) GetData() (size byte, b []byte) {
	size = v.Va[0]
	return size, v.Va[1 : 1+size]
}

type ProxyRequest struct {
	Va [256]byte
}

func (pr *ProxyRequest) SetRemoteInfo(ip string, port int) {
	pr.Va[0] = byte(len(ip))
	copy(pr.Va[1:], []byte(ip))
	pr.Va[pr.Va[0]+1] = byte(port & 0x00FF)
	pr.Va[pr.Va[0]+2] = byte((port & 0xFF00) >> 8)
}

func (pr *ProxyRequest) GetRemoteInfo() (ip string, port int) {
	ip = string(pr.Va[1 : pr.Va[0]+1])
	port = int(pr.Va[pr.Va[0]+1]) + int(pr.Va[pr.Va[0]+2])<<8
	return
}
