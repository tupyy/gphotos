package v1

type EncryptionService interface {
	// Decrypt decrypt the data.
	Decrypt(data string) (string, error)
}
