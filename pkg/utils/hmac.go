// pkg/utils/hmac.go
package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// GenerateHMAC creates an HMAC signature for the given data using the secret
func GenerateHMAC(secret, data string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// ValidateHMAC verifies that the provided signature matches the expected HMAC
func ValidateHMAC(secret, data, signature string) bool {
	expectedMAC := GenerateHMAC(secret, data)
	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}

// GenerateTimestampedHMAC creates an HMAC with timestamp to prevent replay attacks
// Format: HMAC-SHA256(secret, userID:timestamp)
func GenerateTimestampedHMAC(secret, userID string, timestamp int64) string {
	data := fmt.Sprintf("%s:%d", userID, timestamp)
	return GenerateHMAC(secret, data)
}

// ValidateTimestampedHMAC validates HMAC with timestamp and checks if it's within the time window
func ValidateTimestampedHMAC(secret, userID, signature string, timestamp int64, maxAgeSeconds int64) error {
	// Check timestamp is not too old (prevent replay attacks)
	now := time.Now().Unix()
	if now-timestamp > maxAgeSeconds {
		return fmt.Errorf("signature expired: timestamp too old")
	}

	// Check timestamp is not in the future
	if timestamp > now+60 { // Allow 60 second clock skew
		return fmt.Errorf("signature invalid: timestamp in future")
	}

	// Validate HMAC
	data := fmt.Sprintf("%s:%d", userID, timestamp)
	if !ValidateHMAC(secret, data, signature) {
		return fmt.Errorf("signature invalid: HMAC mismatch")
	}

	return nil
}
