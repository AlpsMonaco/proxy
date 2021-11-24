package vpn

import (
	"fmt"
	"io"
	"net"

	"github.com/AlpsMonaco/proxy/stream"
)

/*
protocol part of vpn.
*/

const Version byte = 0x01

const (
	Code_Success byte = iota
	Code_Error
)

type HelloMessage struct {
	msgSize byte
	msg     [255]byte
}

func (hm *HelloMessage) SetMsg(msg string) {
	hm.msgSize = 0
	for i := range []byte(msg) {
		hm.msg[i] = msg[i]
		hm.msgSize++
	}
}

func (hm *HelloMessage) GetMsg() string {
	return string(hm.msg[:hm.msgSize])
}

type GeneralResponse struct {
	code    byte
	msgSize byte
	msg     [256]byte
}

func (gr *GeneralResponse) Set(code byte, msg string) {
	gr.msgSize = 0
	gr.code = code
	for i := 0; i < len(msg); i++ {
		gr.msgSize++
		gr.msg[i] = msg[i]
	}
}

func (gr *GeneralResponse) Get() string {
	if gr.msgSize == 0 {
		return ""
	}
	b := make([]byte, gr.msgSize)
	var i byte
	for i = 0; i < gr.msgSize; i++ {
		b[i] = gr.msg[i]
	}
	return string(b)
}

func (gr *GeneralResponse) GetSize() int {
	return int(1 + 1 + gr.msgSize)
}

type Verify struct {
	va [256]byte
}

func (v *Verify) SetData(size byte, b []byte) {
	v.va[12] = byte(size)
	var i byte
	for i = 0; i < size; i++ {
		v.va[1+i] = b[i]
	}
}

func (v *Verify) GetData() (size byte, b []byte) {
	size = v.va[0]
	return size, v.va[1 : 1+size]
}

type ProxyRequest struct {
	va [256]byte
}

func (pr *ProxyRequest) SetRemoteInfo(ip string, port int) {
	pr.va[0] = byte(len(ip))
	copy(pr.va[1:], []byte(ip))
	pr.va[pr.va[0]+1] = byte(port & 0x00FF)
	pr.va[pr.va[0]+2] = byte((port & 0xFF00) >> 8)
}

func (pr *ProxyRequest) GetRemoteInfo() (ip string, port int) {
	ip = string(pr.va[1 : pr.va[0]+1])
	port = int(pr.va[pr.va[0]+1]) + int(pr.va[pr.va[0]+2])<<8
	return
}

type debugconn struct {
	net.Conn
	Name string
}

func (dc *debugconn) Read(b []byte) (n int, err error) {
	n, err = dc.Conn.Read(b)
	fmt.Printf("[%s]Read %d %v\n", dc.Name, n, b[:n])
	return
}

func (dc *debugconn) Write(b []byte) (n int, err error) {
	n, err = dc.Conn.Write(b)
	fmt.Printf("[%s]Write %d %v\n", dc.Name, n, b[:n])
	return
}

func transport(remote net.Conn, vpnconn net.Conn, encryptor Encryptor, onError func(error)) {
	defer func() {
		closeConn(vpnconn)
		closeConn(remote)
	}()

	remote = &debugconn{remote, "remote"}
	vpnconn = &debugconn{vpnconn, "vpn"}

	var packet *stream.Packet = stream.NewPacket()
	defer stream.FreePacket(packet)

	var buffer []byte = make([]byte, stream.PacketSize)
	var remoteBuffer = buffer[:stream.PacketSize]
	// var vpnBuffer = buffer
	// var remoteBuffer = make([]byte, 128)
	// var vpnBuffer = make([]byte, 256)

	go func() {
		defer closeConn(vpnconn)
		defer closeConn(remote)
		var n int
		var err error
		for {
			n, err = remote.Read(remoteBuffer)
			if n == 0 && err == nil {
				err = io.EOF
			}
			if err != nil {
				onError(err)
				return
			}
			// n, err = encryptor.Encrypt(remoteBuffer[:n], vpnBuffer)
			// if err != nil {
			// onError(err)
			// return
			// }

			// err = packet.WriteStream(vpnconn, vpnBuffer[:n])
			err = packet.WriteStream(vpnconn, remoteBuffer[:n])
			if err != nil {
				onError(err)
				return
			}
		}
	}()

	func() {
		defer closeConn(vpnconn)
		defer closeConn(remote)
		// var n int
		var err error
		for {
			err = packet.Next(vpnconn)
			if err != nil {
				onError(err)
				return
			}
			// n, err = encryptor.Decrypt(packet.Data(), vpnBuffer)
			// if err != nil {
			// onError(err)
			// return
			// }
			_, err = remote.Write(packet.Data())
			// _, err = remote.Write(vpnBuffer[:n])
			if err != nil {
				onError(err)
				return
			}
		}
	}()

}
