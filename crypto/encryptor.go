package crypto

type Encryptor interface {
	Key(key []byte) Encryptor
	Encrypt(plainText []byte) ([]byte, error)
	Decrypt(cipherText []byte) ([]byte, error)
}

type PlainEncryptor struct{}

func (pe *PlainEncryptor) Key(key []byte) Encryptor {
	return pe
}

func (pe *PlainEncryptor) Encrypt(plainText []byte) ([]byte, error) {
	return plainText, nil
}

func (pe *PlainEncryptor) Decrypt(cipherText []byte) ([]byte, error) {
	return cipherText, nil
}
