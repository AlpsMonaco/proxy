package vpn

import (
	"crypto/aes"
	"crypto/cipher"
	"math/rand"
	"time"

	"github.com/AlpsMonaco/proxy/util"
	"golang.org/x/crypto/chacha20poly1305"
)

// Random Number Generator
type RNG interface {
	Size(int)
	Get() []byte
}

var r RNG

type prng struct {
	size int
	buf  [12]byte
}

func (p *prng) Size(i int) {
	p.size = i
}

func (p *prng) Get() []byte {
	for i := 0; i < p.size; i++ {
		p.buf[i] = byte(rand.Intn(255))
	}
	return p.buf[:]
}

func init() {
	if r == nil {
		rand.Seed(time.Now().Unix())
		r = new(prng)
		r.Size(12)
	}
}

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
}

type cipherPlain struct{}

func (p *cipherPlain) Encrypt(plaintext, dst []byte) ([]byte, error) {
	return dst, nil
}

func (p *cipherPlain) Decrypt(ciphertext, dst []byte) ([]byte, error) {
	return dst, nil
}

func (p *cipherPlain) Key(b []byte) Encryptor {
	return p
}

type cipherAes256GCM struct{ aead cipher.AEAD }

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
	nonce := r.Get()
	copy(dst, nonce)
	return c.aead.Seal(dst[:12], nonce, plaintext, nil), nil
}

func (c *cipherAes256GCM) Decrypt(ciphertext, dst []byte) ([]byte, error) {
	nonce := ciphertext[:12]
	encryptedData := ciphertext[12:]
	return c.aead.Open(dst[:0], nonce, encryptedData, nil)
}

type cipherChaCha20Poly1305 struct{ aead cipher.AEAD }

func (c *cipherChaCha20Poly1305) Key(key []byte) Encryptor {
	var err error
	c.aead, err = chacha20poly1305.New([]byte(util.GetMD5String([]byte(key))))
	if err != nil {
		panic(err)
	}
	return c
}

func (c *cipherChaCha20Poly1305) Encrypt(plaintext, dst []byte) ([]byte, error) {
	nonce := r.Get()
	copy(dst, nonce)
	return c.aead.Seal(dst[:12], nonce, plaintext, nil), nil
}

func (c *cipherChaCha20Poly1305) Decrypt(ciphertext, dst []byte) ([]byte, error) {
	nonce := ciphertext[:12]
	encryptedData := ciphertext[12:]
	return c.aead.Open(dst[:0], nonce, encryptedData, nil)
}
