package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// generateRandomHex generates a secure random string of `n` bytes and returns its hex representation.
func GenerateRandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateAppCredentials() (string, string, error) {
	appToken, err := GenerateRandomHex(16) // 32 hex characters
	if err != nil {
		return "", "", err
	}

	appSecret, err := GenerateRandomHex(32) // 64 hex characters
	if err != nil {
		return "", "", err
	}

	return appToken, appSecret, nil
}
