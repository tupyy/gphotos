package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
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
	nonce := make([]byte, nonceSize)

	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

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

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	ret := string(plaintext)

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
