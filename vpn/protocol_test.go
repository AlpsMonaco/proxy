package vpn

import (
	"testing"

	"github.com/AlpsMonaco/proxy/util"
	"golang.org/x/crypto/chacha20poly1305"
)

func assertError(err error) {
	if err != nil {
		panic(err)
	}
}

func TestCrypto(t *testing.T) {
	const key = "This is test key"
	var b []byte = make([]byte, 32)
	t.Logf("%p\n", &b[0])
	nonce := []byte{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4}
	aead, err := chacha20poly1305.New([]byte(util.GetMD5String([]byte(key))))
	assertError(err)
	sealResult := aead.Seal(b[:0], nonce, []byte("1"), nil)
	t.Logf("%p\n", &sealResult[0])
	t.Log("seal result", sealResult)
	openResult, err := aead.Open(sealResult[:0], nonce, sealResult, nil)
	assertError(err)
	t.Log("open result", openResult)
	t.Logf("%p\n", &openResult[0])
}

func TestEncrypt(t *testing.T) {
	var c ChaCha20Poly1305
	c.Key([]byte("123456"))
	data := []byte("123456")
	buffer := make([]byte, 256)
	cipherText := c.Encrypt(data, buffer)
	t.Logf("%p\n", &buffer[0])
	t.Logf("%p\n", &cipherText[0])
	t.Log(cipherText)
	b, _ := c.Decrypt(cipherText, buffer)
	t.Logf("%p\n", &b[0])
}

func TestConnection(t *testing.T) {
	var c ConnectInfo
	c.SetConnection("127.0.0.1", 33333)
	t.Log(c)
	t.Log(c.GetConnection())
	c.SetConnection("www.baidu.com", 443)
	t.Log(c)
	t.Log(c.GetConnection())
	c.SetConnection("www.zhihu.com", 888)
	t.Log(c)
	t.Log(c.GetConnection())
}
