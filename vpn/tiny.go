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

type RemoteConnectionInfo struct {
	va [256]byte
}

func (info *RemoteConnectionInfo) SetInfo(remoteip string, remoteport int) {
	var sizepointer *byte = &info.va[0]
	*sizepointer = 0

	for i := range []byte(remoteip) {
		info.va[i+1] = []byte(remoteip)[i]
		*sizepointer++
	}
	info.va[*sizepointer+1] = byte(remoteport & 0x00FF)
	info.va[*sizepointer+2] = byte((remoteport & 0xFF00) >> 8)
}

func (info *RemoteConnectionInfo) GetInfo() (remoteip string, remoteport int) {
	if info.va[0] == 0 {
		return
	}
	remoteip = string(info.va[1 : info.va[0]+1])
	remoteport = int(info.va[info.va[0]+1]) + int(info.va[info.va[0]+2])<<8
	return
}

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
		buf[0] = 0
		_, _ = client.Write(buf[:1])
		onerror(err)
		return
	}
	buf[0] = 1
	_, err = client.Write(buf[:1])
	if err != nil {
		onerror(err)
		return
	}

	proxy(remoteconn, client)
}

func proxy(remoteconn, clientinfo net.Conn) {

}

type parser struct {
}
