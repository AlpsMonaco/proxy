package socks5

import (
	"testing"
)

func TestClient(t *testing.T) {
	b := []byte{1, 5, 0, 0}
	t.Log(RecvSelectionMessage(b))
}
