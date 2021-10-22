package socks5

import (
	"reflect"
	"unsafe"

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

func (vm *Socks5_VersionMessage) GetMethods() []byte {
	if vm.NumMethod == 0 {
		return nil
	}

	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(vm)) + unsafe.Offsetof(vm.va),
		Len:  int(vm.NumMethod),
		Cap:  int(vm.NumMethod),
	}))
}

func (vm *Socks5_VersionMessage) SetMethod(methods ...byte) {
	vm.NumMethod = byte(len(methods))
	util.SetBytes(vm, int(unsafe.Offsetof(vm.va)), methods)
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
