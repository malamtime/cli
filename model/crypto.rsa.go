package model

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// RSAService implements CryptoService using RSA encryption
type RSAService struct {
	KeySize int
}

// NewRSAService creates a new RSA service with 2048 bit key size
func NewRSAService() CryptoService {
	return &RSAService{
		KeySize: 2048,
	}
}

// GenerateKeys generates a new RSA key pair and returns them as PEM-encoded strings
func (s *RSAService) GenerateKeys() ([]byte, []byte, error) {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, s.KeySize)
	if err != nil {
		return nil, nil, err
	}

	// Extract public key
	publicKey := &privateKey.PublicKey

	// Encode private key to PEM format
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	privateKeyBuf := pem.EncodeToMemory(privateKeyPEM)

	// Encode public key to PEM format
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, nil, err
	}
	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicKeyBuf := pem.EncodeToMemory(publicKeyPEM)

	return publicKeyBuf, privateKeyBuf, nil
}

// Encrypt encrypts data using the provided public key
func (s *RSAService) Encrypt(publicKeyStr string, data []byte) ([]byte, []byte, error) {

	// Decode PEM-encoded public key
	block, _ := pem.Decode([]byte(publicKeyStr))
	if block == nil {
		return nil, nil, errors.New("failed to parse PEM block containing public key")
	}

	// Parse public key
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, nil, err
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, nil, errors.New("not an RSA public key")
	}

	// Encrypt data
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, data)
	if err != nil {
		return nil, nil, err
	}

	return encrypted, nil, nil
}

// Decrypt decrypts data using the provided private key
func (s *RSAService) Decrypt(privateKeyStr string, data []byte, _ []byte) ([]byte, error) {
	panic("no any method of decryption")
}
