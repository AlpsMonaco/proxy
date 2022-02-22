package vpn

import "testing"

func TestHelloMessage(t *testing.T) {
	var hm HelloMessage
	hm.SetMessage("hello")
	t.Log(hm.GetMessage())

	hm.SetBytes([]byte("123"))
	t.Log(hm.GetBytes())
}
