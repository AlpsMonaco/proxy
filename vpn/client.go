package vpn

import (
	"errors"
	"fmt"
	"net"
	"unsafe"

	"github.com/AlpsMonaco/proxy/crypto"
	"github.com/AlpsMonaco/proxy/util"
)

type Client struct {
	IP       string
	Port     int
	Password string
	e        crypto.Encryptor
	SecureConn
}

var (
	ErrConnectFailed = errors.New("Connect to server failed.")
)

func (c *Client) Connect(host string, port uint16) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.IP, c.Port))
	if err != nil {
		return err
	}

	c.e = new(crypto.AESEncryptor).Key([]byte(c.Password))
	c.SecureConn = SecureConn{
		Conn: conn,
		e:    c.e,
	}

	IP := make([]byte, 4)
	util.IPV4AddrToByte(host, &IP)
	buf := make([]byte, ClientBufSize)
	// rm := (*RequestMessage)(unsafe.Pointer(&buf))
	fillRequestMsg((*RequestMessage)(unsafe.Pointer(&buf[0])), AType_IPV4, IP, port)
	_, err = c.SecureConn.Write(buf[:8])
	if err != nil {
		return err
	}

	return nil
}

// current ipv4 onlt
func fillRequestMsg(rm *RequestMessage, addrType byte, IP []byte, port uint16) {
	rm.Ver = VER
	rm.Atype = addrType
	rm.VA[0] = IP[0]
	rm.VA[1] = IP[1]
	rm.VA[2] = IP[2]
	rm.VA[3] = IP[3]
	rm.VA[4] = byte(port >> 8)
	rm.VA[5] = byte(port & 0xF0)
}
