package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	keyVersion = "v1"
)

// EncryptAPIKey encrypts an API key using AES-256-GCM
func EncryptAPIKey(plaintext string) (string, error) {
	key, err := getEncryptionKey()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	// Format: version:nonce:ciphertext (all base64 encoded)
	result := fmt.Sprintf("%s:%s:%s",
		keyVersion,
		base64.StdEncoding.EncodeToString(nonce),
		base64.StdEncoding.EncodeToString(ciphertext),
	)

	return result, nil
}

// DecryptAPIKey decrypts an encrypted API key
// If the input is not in encrypted format (e.g., plain text), returns it as-is
func DecryptAPIKey(encrypted string) (string, error) {
	// If empty, return empty
	if encrypted == "" {
		return "", nil
	}

	key, err := getEncryptionKey()
	if err != nil {
		return "", err
	}

	parts := strings.Split(encrypted, ":")
	// If not in encrypted format (v1:nonce:ciphertext), assume it's plain text
	if len(parts) != 3 {
		return encrypted, nil
	}

	version := parts[0]
	if version != keyVersion {
		// If version doesn't match, might be plain text with colons
		return encrypted, nil
	}

	nonce, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		// Not valid base64, assume plain text
		return encrypted, nil
	}

	ciphertext, err := base64.StdEncoding.DecodeString(parts[2])
	if err != nil {
		// Not valid base64, assume plain text
		return encrypted, nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		// Decryption failed, might be plain text that looks like encrypted format
		return encrypted, nil
	}

	return string(plaintext), nil
}

// getEncryptionKey retrieves and validates the encryption key from environment
func getEncryptionKey() ([]byte, error) {
	keyStr := os.Getenv("ENCRYPTION_KEY")
	if keyStr == "" {
		return nil, errors.New("ENCRYPTION_KEY environment variable not set")
	}

	key, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		return nil, fmt.Errorf("invalid ENCRYPTION_KEY format (must be base64): %w", err)
	}

	if len(key) != 32 {
		return nil, fmt.Errorf("invalid ENCRYPTION_KEY length: expected 32 bytes, got %d", len(key))
	}

	return key, nil
}

// ValidateEncryptionKey validates the encryption key on startup
func ValidateEncryptionKey() error {
	_, err := getEncryptionKey()
	return err
}
