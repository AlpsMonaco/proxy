package vpn

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/AlpsMonaco/proxy/util"
	"golang.org/x/crypto/chacha20poly1305"
)

const (
	CipherPlain CipherEnum = iota
	CipherAes256GCM
	CipherChaCha20poly1305
)

func GetEncryptor(enum CipherEnum, key []byte) Encryptor {
	switch enum {
	case CipherPlain:
		return new(cipherPlain).Key(key)
	case CipherAes256GCM:
		return new(cipherAes256GCM).Key(key)
	case CipherChaCha20poly1305:
		return new(cipherChaCha20Poly1305).Key(key)
	default:
		return new(cipherPlain).Key(key)
	}
}

type Encryptor interface {
	Key([]byte) Encryptor
	Encrypt(plaintext, dst []byte) ([]byte, error)
	Decrypt(ciphertext, dst []byte) ([]byte, error)
	NonceSize() int
}

type cipherPlain struct{}

func (p *cipherPlain) Encrypt(plaintext, dst []byte) ([]byte, error) {
	return plaintext, nil
}

func (p *cipherPlain) Decrypt(ciphertext, dst []byte) ([]byte, error) {
	return ciphertext, nil
}

func (p *cipherPlain) Key(b []byte) Encryptor {
	return p
}

func (p *cipherPlain) NonceSize() int {
	return 0
}

type cipherAes256GCM struct {
	aead cipher.AEAD
}

func (c *cipherAes256GCM) Key(key []byte) Encryptor {
	block, err := aes.NewCipher(util.GetMD5([]byte(key)))
	if err != nil {
		panic(err)
	}
	c.aead, err = cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
	return c
}

func (c *cipherAes256GCM) Encrypt(plaintext, dst []byte) ([]byte, error) {
	nonce := dst[:c.NonceSize()]
	globalRandomNumberGenerator.Read(nonce)
	return c.aead.Seal(dst[:c.NonceSize()], nonce, plaintext, nil), nil
}

func (c *cipherAes256GCM) Decrypt(ciphertext, dst []byte) ([]byte, error) {
	nonce := ciphertext[:c.NonceSize()]
	encryptedData := ciphertext[c.NonceSize():]
	return c.aead.Open(dst[:0], nonce, encryptedData, nil)
}

func (c *cipherAes256GCM) NonceSize() int {
	return c.aead.NonceSize()
}

type cipherChaCha20Poly1305 struct {
	aead cipher.AEAD
}

func (c *cipherChaCha20Poly1305) Key(key []byte) Encryptor {
	var err error
	c.aead, err = chacha20poly1305.New([]byte(util.GetMD5String([]byte(key))))
	if err != nil {
		panic(err)
	}
	return c
}

func (c *cipherChaCha20Poly1305) Encrypt(plaintext, dst []byte) ([]byte, error) {
	nonce := dst[:c.NonceSize()]
	globalRandomNumberGenerator.Read(nonce)
	return c.aead.Seal(dst[:c.NonceSize()], nonce, plaintext, nil), nil
}

func (c *cipherChaCha20Poly1305) Decrypt(ciphertext, dst []byte) ([]byte, error) {
	nonce := ciphertext[:c.NonceSize()]
	encryptedData := ciphertext[c.NonceSize():]
	return c.aead.Open(dst[:0], nonce, encryptedData, nil)
}

func (c *cipherChaCha20Poly1305) NonceSize() int {
	return c.aead.NonceSize()
}
