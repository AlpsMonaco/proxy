package util

import (
	"crypto/md5"
	"errors"
	"reflect"
	"unsafe"
)

var ErrNeedPointer = errors.New("p must be a pointer")
var PtrSize uintptr

func init() {
	var i *byte
	PtrSize = unsafe.Sizeof(&i)
}

// p must be a pointer to a struct
// StructToBinary() will return a pointer to byte slice p covers in the memory.
// return a pointer to a byte slice whose undelying pointer to point to p and size is real size of p.
// so it is safe.
// The underlying byte slice is a struct,so the data of slice field in a struct is not correct.
func StructToBinary(p interface{}) *[]byte {
	return ToBinary(p, GetSize(p))
}

// p must be a pointer to a struct
// size could be calculated by GetSize() or SizeOf().
// ToBinary is much faster,But is not safe.
// Only use when you can make sure arguments are safe.
func ToBinary(p interface{}, size int) *[]byte {
	return (*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: *(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(&p)) + PtrSize)),
		Len:  size,
		Cap:  size,
	}))
}

// GetSize uses reflect,will be a little slow,but it is safe.
// panic when arg p is not a pointer to a struct.
func GetSize(p interface{}) int {
	val := reflect.ValueOf(p)
	if val.Kind() != reflect.Ptr {
		panic(ErrNeedPointer)
	}

	return int(val.Elem().Type().Size())
}

// SizeOf is much faster,get underlying size.
// p shoule be a value of a struct to get the size of a struct instance.
func SizeOf(p interface{}) int {
	return int(*(*uintptr)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&p)))))
}

// p must be a pointer to a struct.
// set p to whatever bPtr holds.
// it's a little bit slower than ToStruct.
func BinaryToStruct(p interface{}, bPtr *[]byte) {
	size := GetSize(p)
	ToStruct(p, bPtr, size)
}

// p must be a pointer to a struct.
// set p to whatever bPtr holds.
// size to define how many bytes bPtr copys to p.
// much faster than BinaryToStruct
func ToStruct(p interface{}, bPtr *[]byte, size int) {
	b := *bPtr
	ptr := (*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(&p)) + PtrSize)))
	for i := 0; i < size; i++ {
		*(*byte)(unsafe.Pointer(uintptr(ptr) + uintptr(i))) = b[i]
	}
}

// p is a pointer to a byte array,not a slice.
func BytesToString(p interface{}) string {
	ptr := (*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(&p)) + PtrSize)))
	var size uintptr

	for {
		if *(*byte)(unsafe.Pointer(ptr + size)) == 0 {
			break
		}
		size++
	}

	return string(*(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: ptr,
		Len:  int(size),
		Cap:  int(size),
	})))
}

// SetBytes set data from address of p to address + len(new) to data new holds.
func SetBytes(p interface{}, offset int, new []byte) {
	ptr := (*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(&p)) + PtrSize))) + uintptr(offset)
	var size uintptr = uintptr(len(new))
	var i uintptr
	for i = 0; i < size; i++ {
		*(*byte)(unsafe.Pointer(ptr + i)) = new[i]
	}
}

// Get any  address of received arg p.
//
func GetAddr(p interface{}) uintptr {
	return (*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(&p)) + PtrSize)))
}

func GetMD5(b []byte) []byte {
	h := md5.New()
	h.Write(b)
	return h.Sum(nil)
}
