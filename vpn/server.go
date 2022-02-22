package vpn

import (
	"bytes"
	"fmt"
	"net"
	"time"
	"unsafe"

	"github.com/AlpsMonaco/proxy/util"
)

func handleError(err error) {
	fmt.Println(err)

}

func closeConn(conn net.Conn) {
	if err := conn.Close(); err != nil {
		conn.Close()
	}
}

type Server struct {
	IP        string
	Port      int
	Key       string
	listener  net.Listener
	encryptor Encryptor
}

func (s *Server) Serve() (err error) {
	s.encryptor = &ChaCha20Poly1305{}
	s.encryptor.Key([]byte(s.Key))
	if err != nil {
		return
	}
	s.listener, err = net.Listen("tcp", util.SprintfAddress(s.IP, s.Port))
	if err != nil {
		return
	}
	for {
		var clientConn net.Conn
		clientConn, err = s.listener.Accept()
		if err != nil {
			return
		}
		go s.newConn(clientConn)
	}
}

func (s *Server) newConn(conn net.Conn) {
	var buffer []byte = make([]byte, allocMemSize)
	n, err := conn.Read(buffer)
	if err != nil {
		handleError(err)
		return
	}
	if !bytes.Equal(clientPrefixReqBytes, buffer[:n]) {
		handleError(ErrPrefixNotMacth)
		return
	}
	_, err = conn.Write(serverPrefixRetBytes)
	if err != nil {
		handleError(err)
		return
	}
	var packet Packet = &SizePacket{}
	var recvbuffer []byte = make([]byte, allocMemSize)
	proxyConn := CreateProxyConn(s.encryptor, packet, conn, buffer, recvbuffer)
	var b []byte

	// block here
	b, err = proxyConn.Recv()
	if err != nil {
		closeConn(conn)
		handleError(err)
		return
	}
	if !bytes.Equal(b, versionBytes) {
		closeConn(conn)
		handleError(ErrVersionDismatch)
		return
	}
	err = proxyConn.Send(versionBytes)
	if err != nil {
		closeConn(conn)
		handleError(err)
		return
	}
	s.acceptConn(&proxyConn)
}

func (s *Server) acceptConn(proxyConn *ProxyConn) {
	b, err := proxyConn.Recv()
	if err != nil {
		handleError(err)
		proxyConn.Close()
		return
	}
	addr, port := GetConnectInfo(b)
	var remoteConn net.Conn
	remoteConn, err = net.DialTimeout("tcp", util.SprintfAddress(addr, port), 10*time.Second)
	if err != nil {
		handleError(err)
		proxyConn.Close()
		return
	}
	err = proxyConn.Send(SuccessBytes)
	if err != nil {
		handleError(err)
		proxyConn.Close()
		return
	}
	BeginProxy(remoteConn, proxyConn)
}

func GetConnectInfo(b []byte) (addr string, port int) {
	ci := (*ConnectInfo)(unsafe.Pointer(&b[0]))
	return ci.GetConnection()
}
