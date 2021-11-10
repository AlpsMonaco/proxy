package vpn

import (
	"fmt"
	"testing"
)

var TestStreamData = []byte{10, 0, 10, 10, 10, 10, 10, 10, 10, 10, 6, 0, 6, 6, 6, 6, 5, 0, 5, 5, 5, 17, 0, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17}

type TestStream struct {
	count int
}

func (ts *TestStream) Read(b []byte) (int, error) {
	ts.count = ts.count + 1
	switch ts.count {
	case 1:
		fmt.Println("pop", TestStreamData[0:3])
		copy(b, TestStreamData[0:3])
		return 3, nil
	case 2:
		fmt.Println("pop", TestStreamData[3:12])
		copy(b, TestStreamData[3:12])
		return 12 - 3, nil
	case 3:
		fmt.Println("pop", TestStreamData[12:21])
		copy(b, TestStreamData[12:21])
		return 21 - 12, nil
	case 4:
		fmt.Println("pop", TestStreamData[21:38])
		copy(b, TestStreamData[21:38])
		return 38 - 21, nil
	default:
		return 0, nil
	}
}

func (ts *TestStream) Write(b []byte) (int, error) {
	fmt.Println("got", b)
	return len(b), nil
}

func TestPacket(t *testing.T) {
	var p Packet
	p.Init()

	var ts TestStream
	var err error

	err = p.Next(&ts)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(p.header.Size)
	fmt.Println(p.GetData())

	p.SetBuffer([]byte{1, 2, 3})
	p.WriteBuffer(&ts, 3)
	err = p.Next(&ts)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(p.header.Size)
	fmt.Println(p.GetData())
	err = p.Next(&ts)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(p.header.Size)
	fmt.Println(p.GetData())
	err = p.Next(&ts)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(p.header.Size)
	fmt.Println(p.GetData())
}
