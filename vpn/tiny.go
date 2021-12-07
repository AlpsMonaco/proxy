package vpn

import (
	"fmt"
	"net"
	"strconv"
	"time"
	"unsafe"

	"github.com/AlpsMonaco/proxy/util"
)

var onerror func(err error)
var log func(s string)

func init() {
	if onerror == nil {
		onerror = func(err error) {
			fmt.Print(format("ERROR", err))
		}
	}
	if log == nil {
		log = func(s string) {
			fmt.Print(format("INFO", s))
		}
	}
}

func format(logtype string, content interface{}) string {
	return fmt.Sprintf("[%s]\t[%s] %v\n", logtype, time.Now().Format("2006-01-02 15:04:05"), content)
}

func SetErrorHandle(errhandle func(error)) {
	onerror = errhandle
}

func StartServer(listenaddr string, listenport int) error {
	listener, err := net.Listen("tcp", util.SprintfAddress(listenaddr, listenport))
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go serve(conn)
	}
}

const (
	code_success byte = iota
	code_error
)

func serve(client net.Conn) {
	buf := make([]byte, 0xFF)
	pointer := unsafe.Pointer(&buf[0])
	_, err := client.Read(buf)
	if err != nil {
		onerror(err)
		return
	}

	remoteinfo := (*RemoteConnectionInfo)(pointer)
	ip, port := remoteinfo.GetInfo()
	log("remote info:" + ip + ":" + strconv.Itoa(port))
	remoteconn, err := net.Dial("tcp", util.SprintfAddress(ip, port))
	if err != nil {
		buf[0] = code_error
		_, _ = client.Write(buf[:1])
		onerror(err)
		return
	}
	buf[0] = code_success
	_, err = client.Write(buf[:1])
	if err != nil {
		onerror(err)
		return
	}
	proxy(remoteconn, client)
}

func proxy(remoteconn, vpnconn net.Conn) {
	// remoteconn = &debugconn{remoteconn, "[remoteconn]" + remoteconn.RemoteAddr().String()}
	// vpnconn = &debugconn{vpnconn, "[vpnconn]" + vpnconn.LocalAddr().String()}

	defer closeConn(remoteconn)
	defer closeConn(vpnconn)
	remotebuf := make([]byte, 0xFFFF-(1<<7))
	clientbuf := make([]byte, 0xFFFF)

	var remotepacket Packet = &Raw{}
	var vpnpacket Packet = &Raw{}

	go func() {
		defer closeConn(remoteconn)
		defer closeConn(vpnconn)
		for {
			err := remotepacket.Next(remoteconn, remotebuf)
			if err != nil {
				onerror(err)
				return
			}

			err = vpnpacket.Send(vpnconn, remotepacket.Data())
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
			err := vpnpacket.Next(vpnconn, clientbuf)
			if err != nil {
				onerror(err)
				return
			}

			err = remotepacket.Send(remoteconn, vpnpacket.Data())
			if err != nil {
				onerror(err)
				return
			}
		}
	}()
}

type client struct {
	serverconn net.Conn
}

func (c *client) Proxy(clientconn net.Conn) {
	proxy(clientconn, c.serverconn)
}

func NewClient(serverconn net.Conn, remoteip string, remoteport int) *client {
	buf := make([]byte, 0xFF)
	pointer := unsafe.Pointer(&buf[0])
	(*RemoteConnectionInfo)(pointer).SetInfo(remoteip, remoteport)

	_, err := serverconn.Write(buf)
	if err != nil {
		onerror(err)
		return nil
	}

	_, err = serverconn.Read(buf)
	if err != nil {
		onerror(err)
		return nil
	}

	if buf[0] != code_success {
		return nil
	}

	return &client{serverconn}
}

type RemoteConnectionInfo struct {
	va [256]byte
}

func (info *RemoteConnectionInfo) SetInfo(remoteip string, remoteport int) {
	info.va[0] = byte(len(remoteip))
	copy(info.va[1:], []byte(remoteip))

	info.va[info.va[0]+1] = byte(remoteport & 0x00FF)
	info.va[info.va[0]+2] = byte((remoteport & 0xFF00) >> 8)
}

func (info *RemoteConnectionInfo) GetInfo() (remoteip string, remoteport int) {
	if info.va[0] == 0 {
		return
	}
	remoteip = string(info.va[1 : info.va[0]+1])
	remoteport = int(info.va[info.va[0]+1]) + int(info.va[info.va[0]+2])<<8
	return
}

func (info *RemoteConnectionInfo) GetSize() int {
	return 1 + 2 + int(info.va[0])
}
