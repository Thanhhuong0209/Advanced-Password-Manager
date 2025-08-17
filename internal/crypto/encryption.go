package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
)

const (
	// Key derivation parameters
	SaltLength     = 32
	KeyLength      = 32 // 256 bits for AES-256
	Iterations     = 100000
	NonceLength    = 12 // GCM nonce size
)

// EncryptedData represents encrypted data with metadata
type EncryptedData struct {
	Salt      []byte `json:"salt"`
	Nonce     []byte `json:"nonce"`
	Ciphertext []byte `json:"ciphertext"`
	Tag       []byte `json:"tag"`
}

// DeriveKey derives a cryptographic key from password using PBKDF2
func DeriveKey(password string, salt []byte) ([]byte, error) {
	if len(salt) != SaltLength {
		return nil, fmt.Errorf("invalid salt length: expected %d, got %d", SaltLength, len(salt))
	}
	
	// Derive key using PBKDF2-SHA256
	key := pbkdf2.Key([]byte(password), salt, Iterations, KeyLength, sha256.New)
	return key, nil
}

// Encrypt encrypts plaintext using AES-256-GCM
func Encrypt(plaintext string, password string) (*EncryptedData, error) {
	// Generate random salt
	salt := make([]byte, SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive key from password
	key, err := DeriveKey(password, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %w", err)
	}

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM mode: %w", err)
	}

	// Generate random nonce
	nonce := make([]byte, NonceLength)
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt and authenticate
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	// Split ciphertext and tag
	tagStart := len(ciphertext) - gcm.Overhead()
	encryptedData := &EncryptedData{
		Salt:      salt,
		Nonce:     nonce,
		Ciphertext: ciphertext[:tagStart],
		Tag:       ciphertext[tagStart:],
	}

	// Zero out sensitive data
	zeroBytes(key)
	
	return encryptedData, nil
}

// Decrypt decrypts ciphertext using AES-256-GCM
func Decrypt(encryptedData *EncryptedData, password string) (string, error) {
	// Validate input
	if encryptedData == nil {
		return "", fmt.Errorf("encrypted data is nil")
	}
	if len(encryptedData.Salt) != SaltLength {
		return "", fmt.Errorf("invalid salt length")
	}
	if len(encryptedData.Nonce) != NonceLength {
		return "", fmt.Errorf("invalid nonce length")
	}

	// Derive key from password
	key, err := DeriveKey(password, encryptedData.Salt)
	if err != nil {
		return "", fmt.Errorf("failed to derive key: %w", err)
	}

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM mode: %w", err)
	}

	// Combine ciphertext and tag
	ciphertext := append(encryptedData.Ciphertext, encryptedData.Tag...)

	// Decrypt and authenticate
	plaintext, err := gcm.Open(nil, encryptedData.Nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	// Zero out sensitive data
	zeroBytes(key)
	
	return string(plaintext), nil
}

// GenerateRandomBytes generates cryptographically secure random bytes
func GenerateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return bytes, nil
}

// HashPassword creates a hash of the password for verification
func HashPassword(password string) (string, error) {
	salt := make([]byte, SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	hash := pbkdf2.Key([]byte(password), salt, Iterations, KeyLength, sha256.New)
	
	// Combine salt and hash
	combined := append(salt, hash...)
	return base64.StdEncoding.EncodeToString(combined), nil
}

// VerifyPassword verifies a password against its hash
func VerifyPassword(password, hash string) (bool, error) {
	// Decode the combined salt+hash
	combined, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}

	if len(combined) != SaltLength+KeyLength {
		return false, fmt.Errorf("invalid hash format")
	}

	// Extract salt and hash
	salt := combined[:SaltLength]
	storedHash := combined[SaltLength:]

	// Derive key from password
	derivedKey := pbkdf2.Key([]byte(password), salt, Iterations, KeyLength, sha256.New)

	// Compare hashes
	return constantTimeCompare(derivedKey, storedHash), nil
}

// constantTimeCompare performs constant-time comparison to prevent timing attacks
func constantTimeCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	
	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}
	return result == 0
}

// zeroBytes overwrites the given slice with zeros to clear sensitive data
func zeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
