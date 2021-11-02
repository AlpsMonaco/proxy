package vpn

import (
	"net"

	"github.com/AlpsMonaco/proxy/crypto"
)

type SecureConn struct {
	net.Conn
	e crypto.Encryptor
}

func (sc *SecureConn) Read(b []byte) (n int, err error) {
	n, err = sc.Conn.Read(b)
	if err != nil {
		return n, err
	}
	plain, err := sc.e.Decrypt(b)
	if err != nil {
		return n, err
	}
	copy(b, plain)
	return len(plain), nil
}

func (sc *SecureConn) Write(b []byte) (n int, err error) {
	cipherText, err := sc.e.Encrypt(b)
	if err != nil {
		return 0, err
	}
	_, err = sc.Conn.Write(cipherText)
	if err != nil {
		return 0, err
	}

	return len(b), err
}
