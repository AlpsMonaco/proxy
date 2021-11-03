package vpn

import (
	"errors"
	"fmt"

	"github.com/AlpsMonaco/proxy/util"
)

const ServerBufSize = ClientBufSize + 16
const ClientBufSize = 264

type RequestMessage struct {
	Ver   byte
	Atype byte
	VA    [256]byte
}

type ResponseMessage struct {
	Ver          byte
	ResponseCode byte
}

const (
	_ byte = iota
	AType_IPV4
	Atype_Domain
)

const VER = 0x01

var (
	ErrVersionDismatch = errors.New("version dismatch")
	ErrAtypeUnknown    = errors.New("Addr type unknown")
)

func (rm *RequestMessage) Parse() (string, error) {
	if rm.Ver != VER {
		return "", ErrVersionDismatch
	}

	var ip string
	var port int
	switch rm.Atype {
	case AType_IPV4:
		for i := 0; i < 4; i++ {
			ip = ip + util.ByteToString(rm.VA[i]) + "."
		}
		ip = ip[:len(ip)-1]
		port = int(rm.VA[4])<<8 + int(rm.VA[5])
	case Atype_Domain:
		ip = string(rm.VA[1:rm.VA[0]])
		port = int(rm.VA[rm.VA[0]])<<8 + int(rm.VA[rm.VA[0]+1])
	default:
		return "", ErrAtypeUnknown
	}

	return fmt.Sprintf("%s:%d", ip, port), nil
}
