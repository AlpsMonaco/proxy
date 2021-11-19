package vpn

import (
	"fmt"
	"net"
)

func GetEncryptor(enum CipherEnum, key []byte) Encryptor {
	switch enum {
	case Cipher_Plain:
		return new(plain).Key(key)
	case Cipher_AES256GCM:
		return nil
	case Cipher_Chacha20poly1305:
		return nil
	default:
		fmt.Println("encryptor not found,using plain")
		return nil
	}
}

type Encryptor interface {
	Key([]byte) Encryptor
	Encrypt(plaintext, dst []byte) (int, error)
	Decrypt(ciphertext, dst []byte) (int, error)
}

type secureConn struct {
	net.Conn
	encryptor Encryptor
	buffer    []byte
}

func NewSecureConn(conn net.Conn, encryptor Encryptor, buffer []byte) *secureConn {
	return &secureConn{
		Conn:      conn,
		encryptor: encryptor,
		buffer:    buffer,
	}
}

func (sc *secureConn) Read(b []byte) (n int, err error) {
	n, err = sc.Conn.Read(sc.buffer)
	if err != nil {
		return
	}
	return sc.encryptor.Decrypt(sc.buffer[:n], b)
}

func (sc *secureConn) Write(b []byte) (n int, err error) {
	if err != nil {
		return
	}
	n, err = sc.encryptor.Encrypt(b, sc.buffer)
	if err != nil {
		return
	}
	return sc.Conn.Write(sc.buffer[:n])
}

type plain struct{}

func (p *plain) Encrypt(plaintext, dst []byte) (int, error) {
	return len(plaintext), nil
}

func (p *plain) Decrypt(ciphertext, dst []byte) (int, error) {
	return len(ciphertext), nil
}

func (p *plain) Key(b []byte) Encryptor {
	return p
}
