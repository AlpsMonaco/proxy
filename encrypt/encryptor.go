package encrypt

type Encryptor interface {
	Encrypt([]byte) []byte
	Decrypt([]byte) []byte
}

type Plain struct{}

func (p *Plain) Encrypt(b []byte) []byte {
	return b
}

func (p *Plain) Decrypt(b []byte) []byte {
	return b
}
