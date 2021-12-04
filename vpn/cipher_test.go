package vpn

import (
	"fmt"
	"testing"

	"github.com/AlpsMonaco/proxy/util"
)

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

func printAddress(a interface{}) {
	fmt.Printf("%p\n", a)
}

func TestMD5(t *testing.T) {
	b := util.GetMD5([]byte("123456"))
	fmt.Println(b)
	t.Logf("%x", b)
}

func TestCipher(t *testing.T) {
	const PlainText = "HelloWorld"
	const PlainTextSize = len(PlainText)

	var key = []byte("key")
	var buf = make([]byte, 0xFF)
	var encryptor Encryptor
	var cipherText, plainText []byte
	var err error

	encryptor = GetEncryptor(CipherPlain, key)
	copy(buf, []byte(PlainText))
	cipherText, err = encryptor.Encrypt(buf[:PlainTextSize], buf)
	assert(err)
	t.Log(cipherText)

	plainText, err = encryptor.Decrypt(cipherText, buf)
	assert(err)
	t.Log(plainText)

	encryptor = GetEncryptor(CipherAes256GCM, key)
	copy(buf[12:], []byte(PlainText))
	cipherText, err = encryptor.Encrypt(buf[encryptor.NonceSize():encryptor.NonceSize()+PlainTextSize], buf)
	assert(err)
	t.Log(cipherText)

	plainText, err = encryptor.Decrypt(cipherText, buf)
	assert(err)
	t.Log(plainText)

	encryptor = GetEncryptor(CipherChaCha20poly1305, key)
	copy(buf[12:], []byte(PlainText))
	cipherText, err = encryptor.Encrypt(buf[encryptor.NonceSize():encryptor.NonceSize()+PlainTextSize], buf)
	assert(err)
	t.Log(cipherText)

	plainText, err = encryptor.Decrypt(cipherText, buf)
	assert(err)
	t.Log(plainText)

}
