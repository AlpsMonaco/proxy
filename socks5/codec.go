package socks5

import (
	"unsafe"

	"github.com/AlpsMonaco/proxy/util"
)

func EncodeVersionMessage(msg *Socks5_VersionMessage) *[]byte {
	size := 2 + msg.NumMethod
	return util.ToBinary(msg, int(size))
}

func DecodeVersionMessage(b *[]byte) *Socks5_VersionMessage {
	return (*Socks5_VersionMessage)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(b))))
}

func EncodeSelectionMessage(msg *Socks5_SelectionMessage) *[]byte {
	return util.ToBinary(msg, 2)
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

func DecodeRequestMessage(b *[]byte) *Socks5_RequestMessage {
	return (*Socks5_RequestMessage)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(b))))
}

func EncodeResponseMessage(msg *Socks5_ResponseMessage) *[]byte {
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
