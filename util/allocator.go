package util

import (
	"sync"
	"unsafe"
)

type Allocator struct {
	b []byte
}

func (a *Allocator) Alloc(size int) {
	a.b = make([]byte, size)
}

func (a *Allocator) GetPointer() unsafe.Pointer {
	return unsafe.Pointer(&a.b[0])
}

func (a *Allocator) GetBytes() []byte {
	return a.b
}

func (a *Allocator) GetByteSize(size int) []byte {
	return a.b[0:size]
}

var p sync.Pool = sync.Pool{
	New: func() interface{} {
		return &Allocator{}
	},
}

func GetAlloctor(size int) *Allocator {
	a := p.Get().(*Allocator)
	if len(a.b) < size {
		a.Alloc(size)
	}
	return a
}

func FreeAllocator(a *Allocator) {
	p.Put(a)
}
