package vpn

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/AlpsMonaco/proxy/socks5"
	"github.com/AlpsMonaco/proxy/stream"
	"github.com/AlpsMonaco/proxy/util"
)

func onerror(err error) {
	util.ErrorCatch(err)
}

func StartServer(listenaddr string, port int) error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", listenaddr, port))
	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go newVpnClient(conn)
	}
}

func newVpnClient(clientconn net.Conn) {
	a := util.GetAlloctor(264)
	_, err := clientconn.Read(a.GetBytes())
	if err != nil {
		onerror(err)
		return
	}

	requestmsg := (*socks5.Socks5_RequestMessage)(a.GetPointer())
	host := requestmsg.GetHost()
	port := requestmsg.GetPort()
	remote, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 10*time.Second)
	if err != nil {
		a.GetBytes()[0] = Code_Error
		_, _ = clientconn.Write(a.GetByteSize(1))
		onerror(err)
		return
	}

	a.GetBytes()[0] = Code_Success
	_, err = clientconn.Write(a.GetByteSize(1))
	if err != nil {
		onerror(err)
		return
	}

	// beginProxy
	proxy(remote, clientconn)
}

type directproxy struct {
	serverconn net.Conn
}

func (dp *directproxy) Proxy(client net.Conn) {
	// buf := make([]byte, 65535)
	// clientBuf := make([]byte, 65535)
	// // var packet = stream.NewPacket()

	// go func() {
	// 	for {
	// 		n, err := dp.serverconn.Read(buf)
	// 		if err != nil {
	// 			onerror(err)
	// 			return
	// 		}

	// 		_, err = client.Write(buf[:n])
	// 		if err != nil {
	// 			onerror(err)
	// 			return
	// 		}
	// 	}
	// }()

	// func() {
	// 	for {
	// 		n, err := client.Read(clientBuf)
	// 		if err != nil {
	// 			onerror(err)
	// 			return
	// 		}

	// 		_, err = dp.serverconn.Write(clientBuf[:n])
	// 		if err != nil {
	// 			onerror(err)
	// 			return
	// 		}
	// 	}
	// }()

	// remote

	proxy(client, dp.serverconn)
}

func ConnectVPN(serverip string, serverport int, remoteip string, remoteport int) socks5.ProxyConn {
	// remoteconn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", remoteip, remoteport), 10*time.Second)
	// if err != nil {
	// 	onerror(err)
	// 	return nil
	// }

	// return &directproxy{remoteconn}

	serverconn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", serverip, serverport), 10*time.Second)
	if err != nil {
		onerror(err)
		return nil
	}
	a := util.GetAlloctor(264)
	requestmsg := (*socks5.Socks5_RequestMessage)(a.GetPointer())
	socks5.FillRequestMessage(requestmsg, socks5.SOCKS5_CMD_CONNECT, remoteip, remoteport)
	_, err = serverconn.Write(a.GetByteSize(requestmsg.GetSize()))
	if err != nil {
		onerror(err)
		return nil
	}

	_, err = serverconn.Read(a.GetBytes())
	if err != nil {
		onerror(err)
		return nil
	}
	if a.GetBytes()[0] != Code_Success {
		return nil
	}

	return &directproxy{serverconn}
}

func proxy(remoteconn, vpnconn net.Conn) {
	defer closeConn(remoteconn)
	defer closeConn(vpnconn)
	buf := make([]byte, stream.PacketSize)
	packet := stream.NewPacket()

	go func() {
		defer closeConn(remoteconn)
		defer closeConn(vpnconn)
		for {
			n, err := remoteconn.Read(buf)
			if n == 0 && err == io.EOF {
				err = io.EOF
			}
			if err != nil {
				onerror(err)
				return
			}
			err = packet.WriteStream(vpnconn, buf[:n])
			if err != nil {
				onerror(err)
				return
			}
		}
	}()

	func() {
		defer closeConn(remoteconn)
		defer closeConn(vpnconn)
		for {

			err := packet.Next(vpnconn)
			if err != nil {
				onerror(err)
				return
			}
			_, err = remoteconn.Write(packet.Data())
			if err != nil {
				onerror(err)
				return
			}
		}
	}()
}
