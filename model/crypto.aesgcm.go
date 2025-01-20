package model

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// AESGCMService implements CryptoService using AES-GCM encryption
type AESGCMService struct {
	KeySize int
}

// NewAESGCMService creates a new AES-GCM service with 256-bit key size
func NewAESGCMService() CryptoService {
	return &AESGCMService{
		KeySize: 32, // 256 bits
	}
}

// GenerateKeys generates a new AES key and returns it as base64 encoded string
// For AES, we only need one key for both encryption and decryption
func (s *AESGCMService) GenerateKeys() ([]byte, []byte, error) {
	// Generate random key
	key := make([]byte, s.KeySize)
	_, err := rand.Read(key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate random key: %v", err)
	}

	// Encode key to base64
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(key)))
	base64.StdEncoding.Encode(encoded, key)

	// Return the same key for both public and private since AES is symmetric
	return encoded, nil, nil
}

// Encrypt encrypts data using AES-GCM
func (s *AESGCMService) Encrypt(keyStr string, data []byte) ([]byte, []byte, error) {
	key, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid key format: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, 12) // GCM standard nonce size
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	ciphertext := aesGCM.Seal(nil, nonce, data, nil)
	return ciphertext, nonce, nil
}

// Decrypt decrypts data using AES-GCM
func (s *AESGCMService) Decrypt(keyStr string, ciphertext, nonce []byte) ([]byte, error) {
	panic("no any method of aes gcm decryption")
}
