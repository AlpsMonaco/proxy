package vpn

import "testing"

func TestProxyRequest(t *testing.T) {
	var pr ProxyRequest
	pr.SetInfo("120.92.17.85", 80)
	t.Log(pr)
	pr.SetInfo("www.baidu.com", 80)
	t.Log(pr)
}

func TestGeneralResponse(t *testing.T) {
	var gr GeneralResponse
	gr.SetMsg("你好")
	t.Log(gr.GetMsg())
}
