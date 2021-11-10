package vpn

import (
	"fmt"
	"io"
	"net"
	"time"
)

type Server struct {
	IP          string
	Port        int
	Key         []byte
	ErrorHandle func(err error)
}

func (s *Server) Listen() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		return err
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go s.newConn(conn)
	}
}

func (s *Server) onError(err error) {
	if s.ErrorHandle != nil {
		s.ErrorHandle(err)
	}
}

func (s *Server) newConn(client net.Conn) {
	var p Packet
	p.Init()
	defer p.Free()
	defer closeConn(client)

	err := p.Next(client)
	if err != nil {
		s.onError(err)
		return
	}

	v := (*Verify)(p.GetPointer())
	if !v.IsKeyMatch() {
		return
	}

	if err = p.Next(client); err != nil {
		s.onError(err)
		return
	}

	pr := (*ProxyRequest)(p.GetPointer())
	host := pr.GetHost()
	port := pr.GetPort()
	remote, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 10*time.Second)
	gr := (*GeneralResponse)(p.GetPointer())
	if err != nil {
		s.onError(err)
		gr.Code = Error
		gr.SetMsg("连接失败")
		err = p.WriteBuffer(client, gr.GetSize())
		if err != nil {
			s.onError(err)
		}
		return
	}
	defer closeConn(remote)
	gr.Code = Success
	gr.SetMsg("可以开始传输数据")
	err = p.WriteBuffer(client, gr.GetSize())
	if err != nil {
		s.onError(err)
		return
	}

	go func() {
		var buf = make([]byte, 64)
		for {
			n, err := remote.Read(buf)
			fmt.Println("remote says", buf)
			if n == 0 {
				err = io.EOF
			}
			if err != nil {
				s.onError(err)
				return
			}
			err = p.WriteSize(client, n)
			if err != nil {
				s.onError(err)
				return
			}
			_, err = client.Write(buf[:n])
			if err != nil {
				s.onError(err)
				return
			}
		}
	}()

	for {
		err = p.Next(client)
		if err != nil {
			s.onError(err)
			return
		}
		_, err = remote.Write(p.GetData())
		if err != nil {
			s.onError(err)
			return
		}
	}
}

func closeConn(conn net.Conn) {
	if err := conn.Close(); err != nil {
		conn.Close()
	}
}
