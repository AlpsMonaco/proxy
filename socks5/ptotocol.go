package socks5

import (
	"github.com/AlpsMonaco/proxy/util"
)

const (
	SOCKS5_VERSION byte = 0x05
)

const (
	SOCKS5_METHOD_NO_AUTH byte = iota
	SOCKS5_METHOD_GSSAPI
	SOCKS5_METHOD_USER_PASSWORD
	SOCKS5_METHOD_NOT_SUPPORT
)

type Socks5_VersionMessage struct {
	Ver       byte
	NumMethod byte
	va        [256]byte
}

func (vm *Socks5_VersionMessage) GetSize() int {
	return 2 + int(vm.NumMethod)
}

type Socks5_SelectionMessage struct {
	Ver    byte
	Method byte
}

const (
	_ byte = iota
	SOCKS5_ATYPE_IPV4
	_
	SOCKS5_ATYPE_DOMAIN
	SOCKS5_ATYPE_IPV6
)

const (
	_ byte = iota
	SOCKS5_CMD_CONNECT
	SOCKS5_CMD_BIND
	SOCKS5_CMD_UDP_FORWARD
)

type Socks5_RequestMessage struct {
	Ver   byte
	Cmd   byte
	Rsv   byte
	Atype byte
	va    [256]byte
}

func (rm *Socks5_RequestMessage) GetSize() int {
	var i int = 4 + 2
	switch rm.Atype {
	case SOCKS5_ATYPE_IPV4:
		return i + 4
	case SOCKS5_ATYPE_DOMAIN:
		return int(rm.va[0]+1) + i
	case SOCKS5_ATYPE_IPV6:
		return i + 16
	default:
		return 0
	}

}

func (rm *Socks5_RequestMessage) GetHost() string {
	var s string
	switch rm.Atype {
	case SOCKS5_ATYPE_IPV4:
		for i := 0; i < 4; i++ {
			s = s + util.ByteToString(rm.va[i]) + "."
		}
		return s[:len(s)-1]
	case SOCKS5_ATYPE_DOMAIN:
		s = string(rm.va[1 : rm.va[0]+1])
		return s
	case SOCKS5_ATYPE_IPV6:
		return ""
	default:
		return ""
	}
}

func (rm *Socks5_RequestMessage) GetPort() int {
	var offset byte
	switch rm.Atype {
	case SOCKS5_ATYPE_IPV4:
		offset = 4
	case SOCKS5_ATYPE_DOMAIN:
		offset = rm.va[0] + 1
	case SOCKS5_ATYPE_IPV6:
		offset = 16
	default:
		return 0
	}

	return (int(rm.va[offset]) << 8) + int(rm.va[offset+1])
}

const (
	SOCKS5_REP_SUCCESS byte = iota
	SOCKS5_REP_CONNECTION_FAILED
	SOCKS5_REP_NOT_ALLOWED
	SOCKS5_REP_NETWORK_UNREACHABLE
	SOCKS5_REP_HOST_UNREACHABLE
	SOCKS5_REP_CONNECTION_REFUSED
	SOCKS5_REP_TTL_TIMEOUT
	SOCKS5_REP_UNSUPPORTED_COMMAND
	SOCKS5_REP_UNSUPPORTED_ATYPE
)

type Socks5_ResponseMessage struct {
	Ver   byte
	Rep   byte
	Rsv   byte
	Atype byte
	va    [256]byte
}

type Socks5_UserPassVerify struct {
	ProtocolVersion byte
}
