package util

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestByteToString(t *testing.T) {
	t.Log(ByteToString(111))
}

func TestStringToByte(t *testing.T) {
	var b byte
	StringToByte("123", &b)
	t.Log(b)

}

func TestIPV4AddrToByte(t *testing.T) {
	var b []byte = make([]byte, 4)

	IPV4AddrToByte("127.0.0.1", (*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&(b[0]))),
		Len:  4,
		Cap:  4,
	})))
	t.Log(b)
}
