package vpn

import (
	"io"
	"net"
)

type Proxy struct{}

type ProxyConn struct {
	encryptor  Encryptor
	packet     Packet
	conn       net.Conn
	buffer     []byte
	recvbuffer []byte
}

func CreateProxyConn(encryptor Encryptor, packet Packet, conn net.Conn, buffer []byte, recvbuffer []byte) ProxyConn {
	return ProxyConn{
		encryptor, packet, conn, buffer, recvbuffer,
	}
}

func (pc *ProxyConn) Send(b []byte) error {
	return pc.packet.Send(pc.conn, pc.encryptor.Encrypt(b, pc.buffer))
}

func (pc *ProxyConn) Recv() ([]byte, error) {
	err := pc.packet.Next(pc.conn, pc.recvbuffer)
	if err != nil {
		return nil, err
	}
	return pc.encryptor.Decrypt(pc.packet.Data(), pc.buffer)
}

func (pc *ProxyConn) Close() {
	closeConn(pc.conn)
}

func BeginProxy(remoteConn net.Conn, proxyConn *ProxyConn) {
	var remoteBuffer = make([]byte, remoteBufferSize)
	go func() {
		var n int
		var err error
		for {
			n, err = remoteConn.Read(remoteBuffer)
			if n == 0 {
				err = io.EOF
			}
			if err != nil {
				handleError(err)
				closeConn(remoteConn)
				proxyConn.Close()
				return
			}
			err = proxyConn.Send(remoteBuffer[:n])
			if err != nil {
				handleError(err)
				closeConn(remoteConn)
				proxyConn.Close()
				return
			}
		}
	}()

	var b []byte
	var err error
	for {
		b, err = proxyConn.Recv()
		if err != nil {
			handleError(err)
			closeConn(remoteConn)
			proxyConn.Close()
			return
		}
		_, err = remoteConn.Write(b)
		if err != nil {
			handleError(err)
			closeConn(remoteConn)
			proxyConn.Close()
			return
		}
	}
}
