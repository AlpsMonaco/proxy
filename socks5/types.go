package socks5

import (
	"reflect"
	"unsafe"

	"github.com/AlpsMonaco/proxy/util"
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

type Socks5_RequestMessage struct {
	Ver   byte
	Cmd   byte
	Rsv   byte
	Atype byte
	va    [256]byte
}

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
