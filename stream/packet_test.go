package stream

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"strconv"
	"testing"
	"time"
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

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

func TestPacket(t *testing.T) {
	var p Packet

	const port = 7777
	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	assert(err)

	go func() {
		conn, err := l.Accept()
		assert(err)
		conn = &connWithLog{
			conn,
		}
		// p.Stream = conn

		go func(conn net.Conn) {

			for {
				err = p.Next(conn)
				assert(err)
				fmt.Println("body is", p.Data(), len(p.Data()))
			}
		}(conn)
	}()

	conn, err := net.Dial("tcp", ":"+strconv.Itoa(port))
	assert(err)
	var buf = make([]byte, 0xFFFF)

	for i := 0; i < 0xFFFF; i++ {
		buf[i] = byte(i % 0xFFFF)
	}
	conn = &connWithLog{
		Conn: conn,
	}

	p.WriteStream(conn, buf[:0xF-1])
	p.WriteStream(conn, buf[:0xF-1])
	p.WriteStream(conn, buf[:0xF-1])
	p.WriteStream(conn, buf[:0xF-2])
	p.WriteStream(conn, buf[:0xF-2])
	p.WriteStream(conn, buf[:0xF-3])
	p.WriteStream(conn, buf[:0xF-3])
	p.WriteStream(conn, buf[:0xF-4])
	p.WriteStream(conn, buf[:0xF-4])
	p.WriteStream(conn, buf[:0xF-5])
	p.WriteStream(conn, buf[:0xF-5])
	p.WriteStream(conn, buf[:0xF-6])
	p.WriteStream(conn, buf[:0xF-6])

	time.Sleep(3 * time.Second)
}

func WritePacketBuffer(conn net.Conn, size int, p *Packet, buf []byte) {
	b := make([]byte, 2)
	b[0] = byte(size & 0x00FF)
	b[1] = byte(size >> 8)
	conn.Write(b)
	conn.Write(buf[:size])
}

// func assert(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }

type connWithLog struct {
	net.Conn
}

func (cl *connWithLog) Write(b []byte) (int, error) {
	fmt.Println("send", b)
	return cl.Conn.Write(b)
}

func (cl *connWithLog) Read(b []byte) (int, error) {
	n, err := cl.Conn.Read(b)
	fmt.Println("read", b[:n])
	return n, err
}

func TestVPNServer(t *testing.T) {
	const port = 7777
	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	assert(err)

	go func() {
		conn, err := l.Accept()
		assert(err)
		go func(conn net.Conn) {
			conn = &connWithLog{
				conn,
			}
			var buf = make([]byte, 0xFF)
			for {
				_, err := conn.Read(buf)
				assert(err)
			}
		}(conn)
	}()

	conn, err := net.Dial("tcp", ":"+strconv.Itoa(port))
	assert(err)
	var buf = make([]byte, 0xFFFF)

	for i := 0; i < 0xFFFF; i++ {
		buf[i] = byte(i % 0xFFFF)
	}
	conn = &connWithLog{
		Conn: conn,
	}

	conn.Write(buf[:0xFF>>2-1])
	conn.Write(buf[:0xFF>>2-2])
	conn.Write(buf[:0xFF>>2-3])
	conn.Write(buf[:0xFF>>2-4])
	conn.Write(buf[:0xFF>>2-5])
	conn.Write(buf[:0xFF>>2-6])
	conn.Write(buf[:0xFF>>2-7])
	conn.Write(buf[:0xFF>>2-8])
	conn.Write(buf[:0xFF>>2-9])
	conn.Write(buf[:0xFF>>2-10])
	conn.Write(buf[:0xFF>>2-11])

	time.Sleep(5 * time.Second)
}

type templatePacket struct {
	net.Conn
	b      []byte
	cursor int
}

func (tp *templatePacket) init() {
	tp.b = PacketGenerate(32321)
	tp.b = append(tp.b, PacketGenerate(10)...)
	tp.b = append(tp.b, PacketGenerate(20)...)
	tp.b = append(tp.b, PacketGenerate(33320)...)
	tp.b = append(tp.b, PacketGenerate(3320)...)
	tp.b = append(tp.b, PacketGenerate(3320)...)
	tp.b = append(tp.b, PacketGenerate(256)...)
}

func (tp *templatePacket) len() int {
	return len(tp.b)
}

func (tp *templatePacket) Read(b []byte) (n int, err error) {
	var forward int = rand.Intn(99) + 40000
	var nextcursor int = tp.cursor + forward
	if nextcursor >= tp.len() {
		nextcursor = tp.len()
	}

	var copysize int = copy(b, tp.b[tp.cursor:nextcursor])
	tp.cursor = nextcursor
	if copysize == 0 {
		err = io.EOF
	}
	return copysize, err
}

func (tp *templatePacket) Write(b []byte) (n int, err error) {
	return 0, nil
}

func PacketGenerate(size int) []byte {
	var result []byte = make([]byte, size+2)
	result[0] = byte(size & 0x00FF)
	result[1] = byte((size & 0xFF00) >> 8)
	for i := 2; i < size+2; i++ {
		result[i] = result[i%2]
	}
	return result
}

func TestPacketParse(t *testing.T) {
	var template templatePacket
	template.init()
	fmt.Println(template.len())
	rand.Seed(time.Now().Unix())
	var packet *Packet = NewPacket()

	c := &connWithLog{&template}
	for {
		err := packet.Next(c)
		assert(err)
		fmt.Println("got", len(packet.Data()), packet.Data())
	}

}
