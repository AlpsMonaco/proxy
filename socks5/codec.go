package socks5

import (
	"unsafe"

	"github.com/AlpsMonaco/proxy/util"
)

func EncodeVersionMessage(msg *Socks5_VersionMessage) *[]byte {
	size := 2 + msg.NumMethod
	return util.ToBinary(msg, int(size))
}

func DecodeSelectionMessage(b *[]byte) *Socks5_SelectionMessage {
	return (*Socks5_SelectionMessage)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(b))))
}

func EncodeRequestMessage(msg *Socks5_RequestMessage) *[]byte {
	var size int = 6

	switch msg.Atype {
	case SOCKS5_ATYPE_IPV4:
		size += 4

	case SOCKS5_ATYPE_DOMAIN:
		// first byte is size
		size += int(msg.va[0]) + 1

	case SOCKS5_ATYPE_IPV6:
		size += 16

	default:
		return nil
	}

	return util.ToBinary(msg, size)
}

func DecodeResponseMessage(b *[]byte) *Socks5_ResponseMessage {
	return (*Socks5_ResponseMessage)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(b))))
}
