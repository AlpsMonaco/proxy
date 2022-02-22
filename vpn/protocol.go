package vpn

import "golang.org/x/crypto/chacha20poly1305"

const clientPrefixReq = "GET /api/v1/portal HTTP1.1\r\nHost: www.zhihu.com\r\n\r\n"
const serverPrefixRet = "HTTP/1.1 200 OK\r\nContent-Type: text/html;api=v1\r\n\r\n"
const version = "20220221"
const SUCCESS = "1"
const FAILED = "0"

const nonceSize = chacha20poly1305.NonceSize
const allocMemSize = PacketSize + nonceSize
const remoteBufferSize = PacketSize - 256

var clientPrefixReqBytes = []byte(clientPrefixReq)
var serverPrefixRetBytes = []byte(serverPrefixRet)
var versionBytes = []byte(version)
var SuccessBytes = []byte(SUCCESS)
var FailedBytes = []byte(FAILED)

type ConnectInfo struct {
	info [256]byte
}

func (ci *ConnectInfo) SetConnection(addr string, port uint16) {
	var totalSize byte = byte(len(addr)) + 3
	ci.info[0] = totalSize
	copy(ci.info[1:], []byte(addr))
	ci.info[totalSize-1] = byte(port & 0x00FF)
	ci.info[totalSize-2] = byte(port & 0xFF00 >> 8)
}

func (ci *ConnectInfo) GetConnection() (addr string, port int) {
	addr = string(ci.info[1 : ci.info[0]-2])
	port = int(ci.info[ci.info[0]-1]) + int(ci.info[ci.info[0]-2])<<8
	return
}

func (ci *ConnectInfo) Size() (size int) {
	size = int(ci.info[0])
	return
}
