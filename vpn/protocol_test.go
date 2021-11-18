package vpn

import (
	"math/rand"
	"testing"
	"time"
)

func TestGeneralResponse(t *testing.T) {
	var gr GeneralResponse
	gr.Set(Code_Success, "成功")
	t.Log(gr.Code)
	t.Log(gr.Get())
	gr.Set(Code_Error, "存在错误")
	t.Log(gr.Code)
	t.Log(gr.Get())
}

func TestRandGenerate(t *testing.T) {
	rand.Seed(time.Now().Unix())
}

func TestProxyRequest(t *testing.T) {
	var pr ProxyRequest
	pr.SetRemoteInfo("118.98.90.87", 65534)
	t.Log(pr.GetRemoteInfo())
}
