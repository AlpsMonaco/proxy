package stream

import (
	"net"

	"github.com/AlpsMonaco/proxy/util"
)

type Stream interface {
	Upstream
	Downstream
}

type Upstream interface {
	Accept() net.Conn
	OnError(err error)
}

type Downstream interface {
	Handle([]byte, net.Conn)
}

func Serve(s Stream) {
	for {
		conn := s.Accept()
		if conn == nil {
			break
		}
		go handle(conn, s)
	}
}

func handle(conn net.Conn, s Stream) {
	a := util.GetAlloctor(PacketSize)
	defer util.FreeAllocator(a)
	var n, i int
	var status byte
	var err error
	var p = GetPacket()

	for {
		n, err = conn.Read(a.Shift(i))
		if err != nil {
			s.OnError(err)
			return
		}
		status = p.Parse(a.GetByteSize(i + n))
		if status == PacketShort {
			i += n
			continue
		}
		s.Handle(*p.Body, conn)
		if status == PacketEqual {
			i = 0
			continue
		}

		// case PacketExtra
		for {
			b := p.ExtraPacket()
			status = p.Parse(b)
			if status == PacketShort {
				copy(a.GetBytes(), b)
				i = len(b)
				break
			}
			s.Handle(*p.Body, conn)
			if status == PacketEqual {
				i = 0
				break
			}
		}
	}
}
