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

type RequestProxy struct {
	Version byte
	Rsv     byte
}

// import "github.com/AlpsMonaco/proxy/util"

// type Verify struct {
// 	Key [16]byte
// }

// var key = [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

// func (v *Verify) IsKeyMatch() bool {
// 	for i := 0; i < 16; i++ {
// 		if v.Key[i] != key[i] {
// 			return false
// 		}
// 	}
// 	return true
// }

// func (v *Verify) SetKey() {
// 	for i := 0; i < 16; i++ {
// 		v.Key[i] = key[i]
// 	}
// }

// type GeneralResponse struct {
// 	Code    byte
// 	MsgSize byte
// 	Msg     [253]byte
// }

// func (gr *GeneralResponse) SetMsg(s string) {
// 	gr.MsgSize = byte(copy(gr.Msg[0:], []byte(s)))
// }

// func (gr *GeneralResponse) GetMsg() string {
// 	return string(gr.Msg[:gr.MsgSize])
// }

// func (gr *GeneralResponse) GetSize() int {
// 	return int(gr.MsgSize + 2)
// }

// type ProxyRequest struct {
// 	Atype byte
// 	Va    [255]byte
// }

// func isDomain(s string) bool {
// 	for _, v := range []byte(s) {
// 		if v >= 58 {
// 			return true
// 		}
// 	}

// 	return false
// }

// func (pr *ProxyRequest) SetInfo(ip string, port int) {
// 	offset := len(ip) + 1
// 	if isDomain(ip) {
// 		pr.Atype = ATYPE_DOMAIN
// 		pr.Va[0] = byte(len(ip))
// 		copy(pr.Va[1:], ip)
// 	} else {
// 		pr.Atype = ATYPE_IPV4
// 		b := make([]byte, 4)
// 		util.IPV4AddrToByte(ip, &b)
// 		for i := 0; i < 4; i++ {
// 			pr.Va[i] = b[i]
// 		}
// 		offset = 4
// 	}

// 	pr.Va[offset] = byte(port & 0xFF00 >> 8)
// 	pr.Va[offset+1] = byte(port & 0x00FF)
// }

// func (pr *ProxyRequest) GetSize() int {
// 	var i int = 1 + 2
// 	switch pr.Atype {
// 	case ATYPE_IPV4:
// 		return i + 4
// 	case ATYPE_DOMAIN:
// 		return int(pr.Va[0]+1) + i
// 	case ATYPE_IPV6:
// 		return i + 16
// 	default:
// 		return 0
// 	}
// }

// func (pr *ProxyRequest) GetHost() string {
// 	var s string
// 	switch pr.Atype {
// 	case ATYPE_IPV4:
// 		for i := 0; i < 4; i++ {
// 			s = s + util.ByteToString(pr.Va[i]) + "."
// 		}
// 		return s[:len(s)-1]
// 	case ATYPE_DOMAIN:
// 		s = string(pr.Va[1 : pr.Va[0]+1])
// 		return s
// 	case ATYPE_IPV6:
// 		return ""
// 	default:
// 		return ""
// 	}
// }

// func (pr *ProxyRequest) GetPort() int {
// 	var offset byte
// 	switch pr.Atype {
// 	case ATYPE_IPV4:
// 		offset = 4
// 	case ATYPE_DOMAIN:
// 		offset = pr.Va[0] + 1
// 	case ATYPE_IPV6:
// 		offset = 16
// 	default:
// 		return 0
// 	}

// 	return (int(pr.Va[offset]) << 8) + int(pr.Va[offset+1])
// }

// const (
// 	Success byte = iota
// 	Error
// )
