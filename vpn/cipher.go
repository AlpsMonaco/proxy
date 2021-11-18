package vpn

import "net"

type Encryptor interface {
	Encrypt(plaintext, dst []byte) (int, error)
	Decrypt(ciphertext, dst []byte) (int, error)
}

type SecureConn struct {
	net.Conn
	encryptor Encryptor
	buffer    []byte
}

func (sc *SecureConn) Read(b []byte) (n int, err error) {
	n, err = sc.Conn.Read(sc.buffer)
	if err != nil {
		return
	}
	return sc.encryptor.Decrypt(sc.buffer[:n], b)
}

func (sc *SecureConn) Write(b []byte) (n int, err error) {
	if err != nil {
		return
	}
	n, err = sc.encryptor.Encrypt(b, sc.buffer)
	if err != nil {
		return
	}
	return sc.Conn.Write(sc.buffer[:n])
}
