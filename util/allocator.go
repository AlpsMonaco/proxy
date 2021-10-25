package util

import (
	"unsafe"
)

type Alloctor struct {
	b []byte
}

func (a *Alloctor) Alloc(size int) {
	a.b = make([]byte, size)
}

func (a *Alloctor) GetPointer() unsafe.Pointer {
	return unsafe.Pointer(&a.b[0])
}

func (a *Alloctor) GetBytes() []byte {
	return a.b
}

func (a *Alloctor) GetByteSize(size int) []byte {
	return a.b[0:size]
}
