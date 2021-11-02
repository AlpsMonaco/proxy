package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
)

type AESEncryptor struct {
	a cipher.AEAD
}

func (ae *AESEncryptor) Key(key []byte) Encryptor {
	hash := md5.Sum(key)
	key = make([]byte, 16)
	for i := 0; i < 16; i++ {
		key[i] = hash[i]
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ae.a, err = cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	return ae
}

func (ae *AESEncryptor) Encrypt(plainText []byte) ([]byte, error) {
	fmt.Println("Encrypt", "plainText", plainText)
	result := make([]byte, 14)
	nonce := result[2:14]
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	fmt.Println("Encrypt", "nonce", nonce)

	cipher := ae.a.Seal(nil, nonce, plainText, nil)
	fmt.Println("Encrypt", "cipher", cipher)
	size := 14 + len(cipher)
	fmt.Println("Encrypt", "size", size)
	result[0] = byte(size >> 8)
	result[1] = byte(size & 0x00FF)
	result = append(result, cipher...)
	fmt.Println("Encrypt", "result", result)
	return result, nil
}

func (ae *AESEncryptor) Decrypt(cipherText []byte) ([]byte, error) {
	var totalSize int = int(cipherText[0])<<8 + int(cipherText[1])
	fmt.Println("Decrypt ", "cipherText", cipherText)
	nonce := cipherText[2:14]
	fmt.Println("Decrypt ", "nonce", nonce)

	cipher := cipherText[14:totalSize]
	fmt.Println("Decrypt ", "cipher", cipher)

	return ae.a.Open(nil, nonce, cipher, nil)
}
