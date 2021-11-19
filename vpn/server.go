package vpn

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/AlpsMonaco/proxy/stream"
	"github.com/AlpsMonaco/proxy/util"
)

type CipherEnum byte

const (
	Cipher_Plain CipherEnum = iota
	Cipher_AES256GCM
	Cipher_Chacha20poly1305
)

type Server struct {
	Addr        string
	Port        int
	Key         []byte
	Cipher      CipherEnum
	ErrorHandle func(err error)
	encryptor   Encryptor
}

func (s *Server) Listen() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Addr, s.Port))
	if err != nil {
		return err
	}
	s.encryptor = GetEncryptor(s.Cipher, s.Key)
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go s.newConn(conn)
	}
}

func (s *Server) newConn(client net.Conn) {
	defer closeConn(client)
	a := util.GetAlloctor(stream.PacketSize)
	defer util.FreeAllocator(a)
	sc := secureConn{
		Conn:      client,
		encryptor: s.encryptor,
		buffer:    a.GetBytes(),
	}
	var err error

	_, err = sc.Read(a.GetBytes())
	if err != nil {
		s.onError(err)
		return
	}
	fmt.Println((*Verify)(a.GetPointer()).GetData())

	(*GeneralResponse)(a.GetPointer()).Set(Code_Success, "success")
	_, err = sc.Write(a.GetByteSize((*GeneralResponse)(a.GetPointer()).GetSize()))
	if err != nil {
		s.onError(err)
		return
	}

	_, err = sc.Read(a.GetBytes())
	if err != nil {
		s.onError(err)
		return
	}
	pr := (*ProxyRequest)(a.GetPointer())
	ip, port := pr.GetRemoteInfo()

	var remote net.Conn
	remote, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), 10*time.Second)
	if err != nil {
		(*GeneralResponse)(a.GetPointer()).Set(Code_Error, "dial to remote failed")
		s.onError(err)
		return
	}
	defer closeConn(remote)

	(*GeneralResponse)(a.GetPointer()).Set(Code_Success, "success")
	_, err = sc.Write(a.GetByteSize((*GeneralResponse)(a.GetPointer()).GetSize()))
	if err != nil {
		s.onError(err)
		return
	}
	s.beginProxy(client, remote, a)
}

func (s *Server) beginProxy(client, remote net.Conn, a *util.Allocator) {
	remoteBuffer := a.GetByteSize(stream.PacketSize - (1 << 8))
	packet := stream.Packet{
		Stream: client,
	}
	packet.SetAllocator(a)

	var n int
	var err error
	go func() {
		defer closeConn(remote)
		defer closeConn(client)
		for {
			n, err = remote.Read(remoteBuffer)
			if err == nil && n == 0 {
				err = io.EOF
			}
			if err != nil {
				s.onError(err)
				return
			}
			// err = packet.WriteData(remoteBuffer[:n])
			n, err = s.encryptor.Encrypt(remoteBuffer[:n], a.GetBytes())
			_, err = client.Write(a.GetByteSize(n))
			if err != nil {
				s.onError(err)
				return
			}
		}
	}()

	func() {
		defer closeConn(remote)
		defer closeConn(client)
		for {
			err = packet.Next()
			if err != nil {
				s.onError(err)
				return
			}
			// _, err = remote.Write(packet.Data())
			n, err = s.encryptor.Decrypt(packet.Data(), a.GetBytes())
			if err != nil {
				s.onError(err)
				return
			}
			_, err = remote.Write(a.GetByteSize(n))
			if err != nil {
				s.onError(err)
				return
			}
		}
	}()
}

func (s *Server) onError(err error) {
	if s.ErrorHandle != nil {
		s.ErrorHandle(err)
	}
}

// import (
// 	"fmt"
// 	"io"
// 	"net"
// 	"time"
// )

// type Server struct {
// 	IP          string
// 	Port        int
// 	Key         []byte
// 	ErrorHandle func(err error)
// }

// func (s *Server) Listen() error {
// 	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
// 	if err != nil {
// 		return err
// 	}
// 	for {
// 		conn, err := l.Accept()
// 		if err != nil {
// 			return err
// 		}
// 		go s.newConn(conn)
// 	}
// }

// func (s *Server) onError(err error) {
// 	if s.ErrorHandle != nil {
// 		s.ErrorHandle(err)
// 	}
// }

// func (s *Server) newConn(client net.Conn) {
// 	var p Packet
// 	p.Init()
// 	defer p.Free()
// 	defer closeConn(client)

// 	err := p.Next(client)
// 	if err != nil {
// 		s.onError(err)
// 		return
// 	}

// 	v := (*Verify)(p.GetPointer())
// 	if !v.IsKeyMatch() {
// 		return
// 	}

// 	if err = p.Next(client); err != nil {
// 		s.onError(err)
// 		return
// 	}

// 	pr := (*ProxyRequest)(p.GetPointer())
// 	host := pr.GetHost()
// 	port := pr.GetPort()
// 	remote, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 10*time.Second)
// 	gr := (*GeneralResponse)(p.GetPointer())
// 	if err != nil {
// 		s.onError(err)
// 		gr.Code = Error
// 		gr.SetMsg("连接失败")
// 		err = p.WriteBuffer(client, gr.GetSize())
// 		if err != nil {
// 			s.onError(err)
// 		}
// 		return
// 	}
// 	defer closeConn(remote)
// 	gr.Code = Success
// 	gr.SetMsg("可以开始传输数据")
// 	err = p.WriteBuffer(client, gr.GetSize())
// 	if err != nil {
// 		s.onError(err)
// 		return
// 	}

// 	go func() {
// 		var buf = make([]byte, 512<<4)
// 		for {
// 			_, err := remote.Read(buf)
// 			if n == 0 && err == nil {
// 				err = io.EOF
// 			}
// 			if err != nil {
// 				s.onError(err)
// 				return
// 			}
// 			err = p.WriteSize(client, n)
// 			if err != nil {
// 				s.onError(err)
// 				return
// 			}
// 			_, err = client.Write(buf[:n])
// 			if err != nil {
// 				s.onError(err)
// 				return
// 			}
// 		}
// 	}()

// 	for {
// 		err = p.Next(client)
// 		if err != nil {
// 			s.onError(err)
// 			return
// 		}
// 		b := p.GetData()
// 		_, err = remote.Write(b)
// 		if err != nil {
// 			s.onError(err)
// 			return
// 		}
// 	}
// }

func closeConn(conn net.Conn) {
	if err := conn.Close(); err != nil {
		conn.Close()
	}
}
