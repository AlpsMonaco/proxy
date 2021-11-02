package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func main() {

}

func demo() {
	// The key argument should be the AES key, either 16 or 32 bytes
	// to select AES-128 or AES-256.
	plainText := make([]byte, 64)
	plainText[0] = 1
	plainText[20] = 1
	key := []byte("0123456789ABCDEF")

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	nonce := make([]byte, 12)

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	fmt.Printf("%v %s\n", nonce, string(nonce))

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	ciphertext := gcm.Seal(nil, nonce, plainText, nil)
	gcm.Open(nil, nonce, ciphertext, nil)
	fmt.Println(ciphertext, len(ciphertext))

	var buf bytes.Buffer

	enc := base64.NewEncoder(base64.StdEncoding, &buf)
	_, err = enc.Write(ciphertext)
	if err != nil {
		panic(err)
	}
	enc.Close()
	fmt.Println(buf.String())
}
