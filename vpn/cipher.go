package vpn

import (
	"fmt"
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

type plain struct{}

func (p *plain) Encrypt(plaintext, dst []byte) (int, error) {
	return len(plaintext), nil
}

func (p *plain) Decrypt(ciphertext, dst []byte) (int, error) {
	copy(dst, ciphertext)
	return len(ciphertext), nil
}

func (p *plain) Key(b []byte) Encryptor {
	return p
}
