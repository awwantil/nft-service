package tools

import (
	"crypto/sha512"
	"encoding/hex"

	"github.com/google/uuid"
)

// GenerateRefreshToken generates a new refresh token using SHA-512 hashing.
// It returns the generated refresh token as a hexadecimal string.
func GenerateRefreshToken() string {
	// Generate a new UUID (Universally Unique Identifier)
	input := uuid.New()

	// Create a new SHA-512 hash instance
	hash := sha512.New()

	// Write the UUID string representation to the hash
	hash.Write([]byte(input.String()))

	// Get the hashed bytes
	hashBytes := hash.Sum(nil)

	// Convert the hashed bytes to a hexadecimal string
	return hex.EncodeToString(hashBytes)
}
