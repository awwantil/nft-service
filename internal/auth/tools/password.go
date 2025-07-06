package tools

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
)

// SaltSize Define salt size
const SaltSize = 16

// GenerateRandomSalt Generate 16 bytes randomly and securely using the
// Cryptographically secure pseudorandom number generator (CSPRNG)
// in the crypto.rand package
func GenerateRandomSalt(saltSize int) []byte {
	var salt = make([]byte, saltSize)

	_, err := rand.Read(salt)

	if err != nil {
		panic(err)
	}

	return salt
}

// HashPassword Combine password and salt then hash them using the SHA-512
// hashing algorithm and then return the hashed password
// as a hex string
func HashPassword(password, secret string, salt []byte) string {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Create sha-512 hasher
	var sha512Hasher = sha512.New()

	// Append salt to password
	passwordBytes = append(passwordBytes, salt...)
	// Append secret to password
	secretBytes := []byte(secret)
	passwordBytes = append(passwordBytes, secretBytes...)

	// Write password bytes to the hasher
	sha512Hasher.Write(passwordBytes)

	// Get the SHA-512 hashed password
	var hashedPasswordBytes = sha512Hasher.Sum(nil)

	// Convert the hashed password to a hex string
	var hashedPasswordHex = hex.EncodeToString(hashedPasswordBytes)

	return hashedPasswordHex
}

// PasswordsMatch Check if two passwords match
func PasswordsMatch(hashedPassword, currPassword, secret string, salt []byte) bool {
	var currPasswordHash = HashPassword(currPassword, secret, salt)

	return hashedPassword == currPasswordHash
}
