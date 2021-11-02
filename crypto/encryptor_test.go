package crypto

import (
	"crypto/md5"
	"fmt"
	"testing"
)

func TestMD5(t *testing.T) {
	b := md5.Sum([]byte("123456"))
	fmt.Println(b)
}
