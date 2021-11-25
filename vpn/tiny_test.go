package vpn

import (
	"errors"
	"testing"
)

func TestConnectionInfo(t *testing.T) {
	var remoteinfo RemoteConnectionInfo
	remoteinfo.SetInfo("www.baidu.com", 80)
	t.Log(remoteinfo.GetInfo())
	remoteinfo.SetInfo("www.googlemap.com", 65523)
	t.Log(remoteinfo.GetInfo())
}

func TestLog(t *testing.T) {
	log("log demo test")
	onerror(errors.New("error"))
}
