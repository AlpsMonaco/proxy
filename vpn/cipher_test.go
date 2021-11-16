package vpn

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"testing"

	"github.com/AlpsMonaco/proxy/util"
	"golang.org/x/crypto/chacha20poly1305"
)

func TestMD5(t *testing.T) {
	b := util.GetMD5([]byte("123456"))
	t.Log(b)
	t.Logf("%x", b)
}

func TestCipher(t *testing.T) {
	pass := "Hello"
	msg := "Pass"

	key := sha256.Sum256([]byte(pass))
	//aead, _ := chacha20poly1305.NewX(key[:])
	aead, _ := chacha20poly1305.New(key[:])

	//nonce := make([]byte, chacha20poly1305.NonceSizeX)
	nonce := make([]byte, chacha20poly1305.NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	ciphertext := aead.Seal(nil, nonce, []byte(msg), nil)
	plaintext, _ := aead.Open(nil, nonce, ciphertext, nil)

	fmt.Printf("Message:\t%s\n", msg)
	fmt.Printf("Passphrase:\t%s\n", pass)
	fmt.Printf("Key:\t%x\n", key)
	fmt.Printf("Nonce:\t%x\n", nonce)
	fmt.Printf("Cipher stream:\t%x\n", ciphertext)
	fmt.Printf("Plain text:\t%s\n", plaintext)
}
