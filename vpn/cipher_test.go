package vpn

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/AlpsMonaco/proxy/util"
	"golang.org/x/crypto/chacha20poly1305"
)

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

func TestMD5(t *testing.T) {
	b := util.GetMD5([]byte("123456"))
	fmt.Println(b)
	t.Logf("%x", b)
}

func TestCipher(t *testing.T) {
	aes256Algo()
}

func aes256Algo() {
	const key = "123456"
	const plainText = "helloworld"
	fmt.Println("==========START==========")
	nonce := make([]byte, 12)
	for i := range nonce {
		nonce[i] = byte(rand.Int31n(255))
	}
	fmt.Println("nonce", len(nonce), nonce)
	block, err := aes.NewCipher(util.GetMD5([]byte(key)))
	assert(err)
	gcm, err := cipher.NewGCM(block)
	assert(err)
	fmt.Println("plainText", len([]byte(plainText)), []byte(plainText))
	result := gcm.Seal(nil, nonce, []byte(plainText), nil)
	fmt.Println("result", len(result), result)
	fmt.Println("==========END==========")
}

func TestChacha20(t *testing.T) {
	chacha20poly1305Algo()
}

func chacha20poly1305Algo() {
	const key = "123456"
	plainText := make([]byte, 10+16)
	copy(plainText, "helloworld")

	fmt.Println("==========START==========")
	aead, err := chacha20poly1305.New([]byte(util.GetMD5String([]byte(key))))
	assert(err)
	nonce := make([]byte, 12)
	for i := range nonce {
		nonce[i] = byte(rand.Int31n(255))
	}
	fmt.Println("nonce", len(nonce), nonce)
	fmt.Println("plainText", len([]byte(plainText)), []byte(plainText))
	result := aead.Seal(plainText[:0], nonce, plainText[:len("helloworld")], nil)
	fmt.Println("result", len(result), result)
	fmt.Println(plainText)
	fmt.Println("==========END==========")
}

func TestBatchAlgo(t *testing.T) {
	rand.Seed(time.Now().Unix())
	aes256Algo()
	aes256Algo()
	aes256Algo()
	aes256Algo()
	aes256Algo()
	aes256Algo()
	aes256Algo()
	chacha20poly1305Algo()
	chacha20poly1305Algo()
	chacha20poly1305Algo()
	chacha20poly1305Algo()
	chacha20poly1305Algo()
	chacha20poly1305Algo()
	chacha20poly1305Algo()
}
