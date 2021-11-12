package vpn

import (
	"fmt"
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

func TestPacket(t *testing.T) {
	var p Packet
	p.Init()

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
			for {
				err = p.Next(conn)
				assert(err)
				fmt.Println(p.GetData())
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

	WritePacketBuffer(conn, 0xFF, &p, buf)
	WritePacketBuffer(conn, 0xFF, &p, buf)
	WritePacketBuffer(conn, 0xFF, &p, buf)
	WritePacketBuffer(conn, 0xFF, &p, buf)
	WritePacketBuffer(conn, 0xFF, &p, buf)
	WritePacketBuffer(conn, 0xFF, &p, buf)
	WritePacketBuffer(conn, 0xFF, &p, buf)
	WritePacketBuffer(conn, 0xFF, &p, buf)
	WritePacketBuffer(conn, 0xFF, &p, buf)

}

func WritePacketBuffer(conn net.Conn, size int, p *Packet, buf []byte) {
	p.WriteSize(conn, size)
	conn.Write(buf[:size])
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

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
