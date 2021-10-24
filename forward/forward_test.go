package forward

import (
	"fmt"
	"testing"
)

func TestForward(t *testing.T) {
	f := &TCPForward{
		ListenAddr: "localhost:7899",
		HostAddr:   "127.0.0.1:7890",
		OnError:    func(e error) { fmt.Println(e) },
	}

	f.Start()
}
