package encryption_test

import (
	"crypto/rand"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tupyy/gophoto/utils/encryption"
)

func TestEncryption(t *testing.T) {
	key := make([]byte, 24)

	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		t.Fatal(err)
	}

	g := encryption.NewGenerator(string(key))

	data := "hey"
	encryptedData, err := g.EncryptData(data)
	assert.Nil(t, err)
	assert.NotEmpty(t, encryptedData)

	decrypted, err := g.DecryptData(encryptedData)
	assert.Nil(t, err)
	assert.Equal(t, data, decrypted)
}
