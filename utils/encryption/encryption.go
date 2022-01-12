package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const (
	nonceSize = 12
)

type Generator struct {
	key []byte
}

func NewGenerator(key string) *Generator {
	return &Generator{key: []byte(key)}
}

// EncryptData encrypts the data with the provided key. The result is not determinist.
func (g *Generator) EncryptData(data string) (string, error) {
	hash := sha256.Sum256(bytes.NewBufferString(data).Bytes())

	nonce := make([]byte, nonceSize)
	copy(nonce, hash[:nonceSize])

	return g.encryptDataWithNonce(data, nonce)
}

func (g *Generator) encryptDataWithNonce(data string, nonce []byte) (string, error) {
	aesgcm, err := g.getAEAD()
	if err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, []byte(data), nil)

	ret := fmt.Sprintf("%x%x", nonce, ciphertext)

	return ret, nil
}

// DecryptData decrypts the data with the provided key.
func (g *Generator) DecryptData(data string) (string, error) {
	aesgcm, err := g.getAEAD()
	if err != nil {
		return "", err
	}

	nonce, ciphertext, err := parseData(data)
	if err != nil {
		return "", err
	}

	hexEncodedText, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	// b, err := hex.DecodeString(string(hexEncodedText))
	// if err != nil {
	// 	return "", err
	// }

	ret := string(hexEncodedText)

	return ret, nil
}

func (g *Generator) getAEAD() (cipher.AEAD, error) {
	block, err := aes.NewCipher(g.key)
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(block)
}

func parseData(data string) ([]byte, []byte, error) {
	databytes, err := hex.DecodeString(data)
	if err != nil {
		return nil, nil, err
	}

	return databytes[:nonceSize], databytes[nonceSize:], nil
}
