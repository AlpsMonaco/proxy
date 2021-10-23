package forward

import (
	"fmt"
	"net"
	"testing"
)

func TestForward(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:7899")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		dst, err := net.Dial("tcp", "localhost:7890")
		if err != nil {
			t.Log(err)
		} else {
			go NewForward(dst, conn, func(e error) { fmt.Println(e) }).Start()
		}

	}
}
