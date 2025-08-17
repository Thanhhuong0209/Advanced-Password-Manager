package crypto

import (
	"strings"
	"testing"
)

func TestDeriveKey(t *testing.T) {
	password := "testpassword"
	salt := make([]byte, SaltLength)
	
	key, err := DeriveKey(password, salt)
	if err != nil {
		t.Fatalf("DeriveKey failed: %v", err)
	}
	
	if len(key) != KeyLength {
		t.Errorf("Expected key length %d, got %d", KeyLength, len(key))
	}
	
	// Test with invalid salt length
	invalidSalt := make([]byte, 16)
	_, err = DeriveKey(password, invalidSalt)
	if err == nil {
		t.Error("Expected error for invalid salt length")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	plaintext := "Hello, World! This is a test message."
	password := "mypassword123"
	
	// Encrypt
	encrypted, err := Encrypt(plaintext, password)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}
	
	// Verify encrypted data structure
	if len(encrypted.Salt) != SaltLength {
		t.Errorf("Expected salt length %d, got %d", SaltLength, len(encrypted.Salt))
	}
	if len(encrypted.Nonce) != NonceLength {
		t.Errorf("Expected nonce length %d, got %d", NonceLength, len(encrypted.Nonce))
	}
	
	// Decrypt
	decrypted, err := Decrypt(encrypted, password)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}
	
	if decrypted != plaintext {
		t.Errorf("Expected decrypted text '%s', got '%s'", plaintext, decrypted)
	}
}

func TestEncryptDecryptEmptyString(t *testing.T) {
	plaintext := ""
	password := "testpass"
	
	encrypted, err := Encrypt(plaintext, password)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}
	
	decrypted, err := Decrypt(encrypted, password)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}
	
	if decrypted != plaintext {
		t.Errorf("Expected decrypted text '%s', got '%s'", plaintext, decrypted)
	}
}

func TestEncryptDecryptLongText(t *testing.T) {
	// Create a long text
	plaintext := strings.Repeat("This is a very long text that should be encrypted and decrypted properly. ", 100)
	password := "verylongpassword123"
	
	encrypted, err := Encrypt(plaintext, password)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}
	
	decrypted, err := Decrypt(encrypted, password)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}
	
	if decrypted != plaintext {
		t.Error("Long text encryption/decryption failed")
	}
}

func TestDecryptWithWrongPassword(t *testing.T) {
	plaintext := "secret message"
	password := "correctpassword"
	wrongPassword := "wrongpassword"
	
	encrypted, err := Encrypt(plaintext, password)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}
	
	// Try to decrypt with wrong password
	_, err = Decrypt(encrypted, wrongPassword)
	if err == nil {
		t.Error("Expected error when decrypting with wrong password")
	}
}

func TestDecryptWithNilData(t *testing.T) {
	_, err := Decrypt(nil, "password")
	if err == nil {
		t.Error("Expected error when decrypting nil data")
	}
}

func TestGenerateRandomBytes(t *testing.T) {
	length := 32
	bytes, err := GenerateRandomBytes(length)
	if err != nil {
		t.Fatalf("GenerateRandomBytes failed: %v", err)
	}
	
	if len(bytes) != length {
		t.Errorf("Expected %d bytes, got %d", length, len(bytes))
	}
	
	// Generate another set and ensure they're different
	bytes2, err := GenerateRandomBytes(length)
	if err != nil {
		t.Fatalf("GenerateRandomBytes failed: %v", err)
	}
	
	if len(bytes2) != length {
		t.Errorf("Expected %d bytes, got %d", length, len(bytes2))
	}
	
	// Very unlikely that two random byte arrays are identical
	identical := true
	for i := 0; i < length; i++ {
		if bytes[i] != bytes2[i] {
			identical = false
			break
		}
	}
	
	if identical {
		t.Error("Two random byte arrays should not be identical")
	}
}

func TestHashPassword(t *testing.T) {
	password := "testpassword"
	
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	
	if hash == "" {
		t.Error("Hash should not be empty")
	}
	
	// Hash the same password again - should be different due to different salt
	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	
	if hash == hash2 {
		t.Error("Two hashes of the same password should be different due to different salts")
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "testpassword"
	
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	
	// Verify with correct password
	valid, err := VerifyPassword(password, hash)
	if err != nil {
		t.Fatalf("VerifyPassword failed: %v", err)
	}
	if !valid {
		t.Error("Password verification should succeed with correct password")
	}
	
	// Verify with wrong password
	valid, err = VerifyPassword("wrongpassword", hash)
	if err != nil {
		t.Fatalf("VerifyPassword failed: %v", err)
	}
	if valid {
		t.Error("Password verification should fail with wrong password")
	}
}

func TestVerifyPasswordInvalidHash(t *testing.T) {
	// Test with invalid hash format
	_, err := VerifyPassword("password", "invalidhash")
	if err == nil {
		t.Error("Expected error with invalid hash format")
	}
}

func TestConstantTimeCompare(t *testing.T) {
	a := []byte("hello")
	b := []byte("hello")
	c := []byte("world")
	d := []byte("hell")
	
	if !constantTimeCompare(a, b) {
		t.Error("Identical byte arrays should compare equal")
	}
	
	if constantTimeCompare(a, c) {
		t.Error("Different byte arrays should not compare equal")
	}
	
	if constantTimeCompare(a, d) {
		t.Error("Different length byte arrays should not compare equal")
	}
}

func TestZeroBytes(t *testing.T) {
	bytes := []byte{1, 2, 3, 4, 5}
	zeroBytes(bytes)
	
	for _, b := range bytes {
		if b != 0 {
			t.Error("All bytes should be zero after zeroBytes")
		}
	}
}
