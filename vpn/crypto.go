package vpn

import (
	"crypto/cipher"

	"github.com/AlpsMonaco/proxy/util"
	"golang.org/x/crypto/chacha20poly1305"
)

func initAead(pAead *cipher.AEAD, key string) (err error) {
	*pAead, err = chacha20poly1305.New([]byte(util.GetMD5String([]byte(key))))
	return err
}

const nonce = "jKBBqFWvYZi2"

var nonceBytes = []byte(nonce)

func GetNonce(dst []byte) {
	copy(dst, nonceBytes)
}

type Encryptor interface {
	Key([]byte)
	Encrypt(plainText []byte, buffer []byte) []byte
	Decrypt(cipherText []byte, buffer []byte) ([]byte, error)
}

type ChaCha20Poly1305 struct {
	aead cipher.AEAD
}

func (c *ChaCha20Poly1305) Key(key []byte) {
	var err error
	c.aead, err = chacha20poly1305.New([]byte(util.GetMD5String([]byte(key))))
	if err != nil {
		panic(err)
	}
}

func (c *ChaCha20Poly1305) Encrypt(plainText []byte, buffer []byte) (b []byte) {
	GetNonce(buffer)
	b = c.aead.Seal(buffer[nonceSize:][:0], nonceBytes, plainText, nil)
	return buffer[:len(b)+nonceSize]
}

func (c *ChaCha20Poly1305) Decrypt(cipherText []byte, buffer []byte) ([]byte, error) {
	nonce := cipherText[:nonceSize]
	return c.aead.Open(buffer[:0], nonce, cipherText[nonceSize:], nil)
}

func GetEncrypt(key []byte) Encryptor {
	var c ChaCha20Poly1305
	c.Key(key)
	return &c
}
