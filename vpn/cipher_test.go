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

func TestMD5(t *testing.T) {
	b := util.GetMD5([]byte("123456"))
	fmt.Println(b)
	t.Logf("%x", b)
}

func TestAes256GCM(t *testing.T) {
	const key = "123456"
	encryptor := GetEncryptor(CipherAes256GCM, []byte(key))
	dst := make([]byte, 2048)
	result, err := encryptor.Encrypt([]byte("HelloWorld"), dst)
	assert(err)
	t.Log(result)
	result, err = encryptor.Decrypt(result, dst)
	assert(err)
	t.Log(result)
}

func TestChaCha20Poly(t *testing.T) {
	const key = "123456"
	encryptor := GetEncryptor(CipherChaCha20poly1305, []byte(key))
	dst := make([]byte, 2048)
	result, err := encryptor.Encrypt([]byte("HelloWorld"), dst)
	assert(err)
	t.Log(result)
	result, err = encryptor.Decrypt(result, dst)
	assert(err)
	t.Log(result)
}
