package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
)

// deriveKey hashes secretKey with SHA-256 to produce a 32-byte AES-256 key.
func deriveKey(secretKey string) []byte {
	h := sha256.Sum256([]byte(secretKey))
	return h[:]
}

// Encrypt encrypts plaintext using AES-256-GCM and returns hex-encoded ciphertext.
// If secretKey or plaintext is empty, plaintext is returned unchanged.
func Encrypt(plaintext, secretKey string) (string, error) {
	if secretKey == "" || plaintext == "" {
		return plaintext, nil
	}

	block, err := aes.NewCipher(deriveKey(secretKey))
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherBytes := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return hex.EncodeToString(cipherBytes), nil
}

// Decrypt decrypts a hex-encoded AES-256-GCM ciphertext.
// If secretKey is empty, ciphertext is returned unchanged (no encryption mode).
// If the value is not valid hex it is treated as plaintext (migration safety).
func Decrypt(ciphertext, secretKey string) (string, error) {
	if secretKey == "" || ciphertext == "" {
		return ciphertext, nil
	}

	data, err := hex.DecodeString(ciphertext)
	if err != nil {
		// Value was stored before encryption was enabled — return as-is.
		return ciphertext, nil
	}

	block, err := aes.NewCipher(deriveKey(secretKey))
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("crypto: ciphertext too short")
	}

	nonce, cipherBytes := data[:nonceSize], data[nonceSize:]
	plain, err := aesGCM.Open(nil, nonce, cipherBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plain), nil
}
