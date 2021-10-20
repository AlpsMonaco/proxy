package socks5

const (
	SOCKS5_METHOD_NO_AUTH uint8 = iota
	SOCKS5_METHOD_GSSAPI
	SOCKS5_METHOD_USER_PASSWORD
)

type Socks5_VersionMessage struct {
	Ver       uint8
	NumMethod uint8
	Methods   [256]uint8
}

type Socks5_SelectionMessage struct {
	Ver    uint8
	Method uint8
}

type Socks5_RequestMessage struct {
}
