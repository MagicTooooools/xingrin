package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const (
	// Django default PBKDF2 iterations
	defaultIterations = 600000
	// Key length for PBKDF2
	keyLength = 32
)

// VerifyDjangoPassword verifies a password against Django's pbkdf2_sha256 hash
// Django format: pbkdf2_sha256$iterations$salt$hash
func VerifyDjangoPassword(password, encoded string) bool {
	parts := strings.Split(encoded, "$")
	if len(parts) != 4 {
		return false
	}

	algorithm := parts[0]
	if algorithm != "pbkdf2_sha256" {
		return false
	}

	iterations, err := strconv.Atoi(parts[1])
	if err != nil {
		return false
	}

	salt := parts[2]
	expectedHash := parts[3]

	// Decode the expected hash from base64
	expectedBytes, err := base64.StdEncoding.DecodeString(expectedHash)
	if err != nil {
		return false
	}

	// Compute PBKDF2 hash
	computedHash := pbkdf2.Key([]byte(password), []byte(salt), iterations, len(expectedBytes), sha256.New)

	// Constant-time comparison to prevent timing attacks
	return subtle.ConstantTimeCompare(computedHash, expectedBytes) == 1
}

// HashPassword creates a Django-compatible password hash
func HashPassword(password string) (string, error) {
	// Generate random salt (12 bytes, base64 encoded = 16 chars)
	saltBytes := make([]byte, 12)
	if _, err := rand.Read(saltBytes); err != nil {
		return "", err
	}
	salt := base64.StdEncoding.EncodeToString(saltBytes)

	// Compute PBKDF2 hash
	hash := pbkdf2.Key([]byte(password), []byte(salt), defaultIterations, keyLength, sha256.New)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)

	// Format: pbkdf2_sha256$iterations$salt$hash
	return fmt.Sprintf("pbkdf2_sha256$%d$%s$%s", defaultIterations, salt, hashBase64), nil
}
