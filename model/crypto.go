package model

// CryptoService defines the interface for encryption/decryption operations
type CryptoService interface {
	// GenerateKeys generates a public and private key pair.
	// in RSA service, public key, private key and error
	// in AES-GCM service, only public key and error
	GenerateKeys() ([]byte, []byte, error)

	// Encrypt encrypts data using the provided key.
	// in RSA service, encrypted data, nil and error
	// in AES-GCM service, encrypted data, nonce and error
	Encrypt(key string, data []byte) ([]byte, []byte, error)

	// Decrypt decrypts data using the provided key and nonce.
	// in RSA service, key, data, nil
	// in AES-GCM service, key, data, nonce
	Decrypt(key string, data []byte, nonce []byte) ([]byte, error)
}
