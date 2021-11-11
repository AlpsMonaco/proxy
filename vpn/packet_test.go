package vpn

import (
	"fmt"
	"testing"
)

var TestStreamData = []byte{10, 0, 10, 10, 10, 10, 10, 10, 10, 10, 6, 0, 6, 6, 6, 6, 5, 0, 5, 5, 5, 17, 0, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 3, 0, 1, 250, 250}

type TestStream struct {
	count int
}

func (ts *TestStream) Read(b []byte) (int, error) {
	ts.count = ts.count + 1
	switch ts.count {
	case 1:
		fmt.Println("pop", TestStreamData[0:20])
		copy(b, TestStreamData[0:20])
		return 20 - 0, nil
	case 2:
		fmt.Println("pop", TestStreamData[20:21])
		copy(b, TestStreamData[20:21])
		return 1, nil
	case 3:
		fmt.Println("pop", TestStreamData[21:22])
		copy(b, TestStreamData[21:22])
		return 1, nil
	case 4:
		fmt.Println("pop", TestStreamData[22:38])
		copy(b, TestStreamData[22:38])
		return 38 - 22, nil
	case 5:
		fmt.Println("pop", TestStreamData[38:41])
		copy(b, TestStreamData[38:41])
		return 41 - 38, nil
	case 6:
		fmt.Println("pop", TestStreamData[41:43])
		copy(b, TestStreamData[41:43])
		return 2, nil
	case 7:
		nb := make([]byte, 0xFAFA)
		for i := 0; i < 0xFAFA-3; i++ {
			nb[i] = byte(i)
		}
		fmt.Println("pop", "final")
		copy(b, nb)
		return 0xFAFA - 3, nil
	case 8:
		fmt.Println("pop", "final")
		nb := []byte{0x8, 0x8, 0x8}
		copy(b, nb)
		return 3, nil
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
	fmt.Println(len(p.GetData()))
	fmt.Println(p.GetData())
}

func TestVPNServer(t *testing.T) {

}
