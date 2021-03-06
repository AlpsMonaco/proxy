package util

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"unsafe"
)

func Assert(err error) {
	if err != nil {
		panic(err)
	}
}

func AssertEqual(expected, actual int) {
	if expected != actual {
		Assert(errors.New(fmt.Sprint("expected != actual", expected, actual)))
	}
}

func AssertEqualByteSlice(expected, actual *[]byte) {
	var err error = fmt.Errorf("expcet:%v\nactual:%v\nnot equal", expected, actual)
	if len(*expected) != len(*actual) {
		Assert(err)
	}

	for i := 0; i < len(*expected); i++ {
		if (*expected)[i] != (*actual)[i] {
			Assert(err)
		}
	}
}

type TestStruct1 struct {
	ui32 uint32
	i32  int32
	i64  int64
}

type TestStruct2 struct {
	b    byte
	ui64 uint64
	i    int
	i8   int8
}

func TestBinaryEncode(t *testing.T) {
	var b []byte = []byte{32, 0, 0, 0, 64, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0}
	t1 := TestStruct1{
		ui32: 32,
		i32:  64,
		i64:  128,
	}

	AssertEqualByteSlice(&b, StructToBinary(&t1))

	t2 := TestStruct2{
		b:    10,
		ui64: 2 << 32,
		i:    2 << 20,
		i8:   -10,
	}

	b = []byte{10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 246, 0, 0, 0, 0, 0, 0, 0}

	AssertEqualByteSlice(&b, StructToBinary(&t2))
}

func TestSizeInfo(t *testing.T) {
	t1 := TestStruct1{
		ui32: 32,
		i32:  64,
		i64:  128,
	}

	AssertEqual(GetSize(&t1), SizeOf(t1))

	t2 := TestStruct2{
		b:    1,
		ui64: 65535,
		i:    65535 << 1,
		i8:   2 << 3,
	}

	AssertEqual(GetSize(&t2), SizeOf(t2))
}

func TestBinaryToStruct(t *testing.T) {
	var t1, t2 TestStruct1
	t1.i32 = 32
	t1.i64 = 64
	t1.ui32 = 128

	BinaryToStruct(&t2, StructToBinary(&t1))
	if t1 != t2 {
		t.Fatal("t1 != t2")
	}

	var t3, t4 TestStruct2
	t3.b = 1
	t3.i = 12
	t3.i8 = 8
	t3.ui64 = 765
	BinaryToStruct(&t4, StructToBinary(&t3))
	if t3 != t4 {
		t.Fatal("t3 != t4")
	}
}

func TestBytesToString(t *testing.T) {
	b := []byte("你好，世界")
	var bytes [256]byte

	for i := 0; i < len(b); i++ {
		bytes[i] = b[i]
	}

	fmt.Println(BytesToString(&bytes))
}

func TestSetBytes(t *testing.T) {
	// var b []byte = []byte{32, 0, 0, 0, 64, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0}
	t1 := TestStruct1{
		ui32: 32,
		i32:  64,
		i64:  128,
	}

	SetBytes(&t1, 4, []byte{32, 0, 0, 0, 125, 1})

	t.Log(t1)
}

func TestEndian(t *testing.T) {
	var port uint16 = 443
	fmt.Printf("%x\n", port)
	fmt.Printf("0x%02x %08b\n", port>>8, port>>8)
	fmt.Printf("0x%02x %08b\n", port<<8, port<<8)
	fmt.Printf("0x%02x %08b\n", port&0xFF00, port&0xFF00)
	fmt.Printf("0x%02x %08b\n", port&0x00FF, port&0x00FF)

}

func TestBinary(t *testing.T) {
	var a uint32 = 62537885
	var size int = 4

	fmt.Printf("0x%08x\n", a)
	fmt.Printf("0x%08x\n", a>>8)
	fmt.Printf("0x%08x\n", a<<8)
	fmt.Printf("0x%08x\n", a&0x00FF)
	fmt.Printf("0x%08x\n", a&0xFF00)

	b := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr((unsafe.Pointer(&a))),
		Len:  size,
		Cap:  size,
	}))

	for i := 0; i < size; i++ {
		fmt.Printf("0x%02x\n", b[i])
	}
}

func runtimeCaller(a ...interface{}) {
	fmt.Println(runtime.Caller(1))
	fmt.Println(a...)
}

func TestCaller(t *testing.T) {
	runtimeCaller(1)
}

func TestHttpRequest(t *testing.T) {
	var data []byte = []byte("GET /api/request HTTP/1.1\r\nHost: www.gogames.com\r\n\r\n")
	t.Log(data)
	var b []byte
	t.Log(b == nil)
}

func TestSlice(t *testing.T) {
	buffer := make([]byte, 64)
	bufferPart1 := buffer[:64]
	bufferPart2 := buffer[:32]

	bufferPart1[1] = 1
	bufferPart1[2] = 2
	bufferPart1[3] = 3
	bufferPart1[4] = 4

	t.Log(*(*reflect.SliceHeader)(unsafe.Pointer(&buffer)))
	t.Log(*(*reflect.SliceHeader)(unsafe.Pointer(&bufferPart1)))
	t.Log(*(*reflect.SliceHeader)(unsafe.Pointer(&bufferPart2)))

	t.Log(buffer)
	t.Log(bufferPart1)
	t.Log(bufferPart2[:5])

}
