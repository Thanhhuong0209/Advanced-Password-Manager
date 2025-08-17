package generator

import (
	"strings"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	if config.Length != 16 {
		t.Errorf("Expected default length 16, got %d", config.Length)
	}
	if !config.Uppercase {
		t.Error("Expected default uppercase to be true")
	}
	if !config.Lowercase {
		t.Error("Expected default lowercase to be true")
	}
	if !config.Numbers {
		t.Error("Expected default numbers to be true")
	}
	if !config.Symbols {
		t.Error("Expected default symbols to be true")
	}
	if config.Exclude != "" {
		t.Error("Expected default exclude to be empty")
	}
	if !config.NoRepeating {
		t.Error("Expected default no-repeating to be true")
	}
}

func TestGeneratePasswordDefault(t *testing.T) {
	password, err := GeneratePassword(nil)
	if err != nil {
		t.Fatalf("GeneratePassword failed: %v", err)
	}
	
	if len(password) != 16 {
		t.Errorf("Expected password length 16, got %d", len(password))
	}
	
	// Check that password contains at least one character from each set
	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSymbol := false
	
	for _, char := range password {
		if strings.ContainsRune(Uppercase, char) {
			hasUpper = true
		} else if strings.ContainsRune(Lowercase, char) {
			hasLower = true
		} else if strings.ContainsRune(Numbers, char) {
			hasNumber = true
		} else if strings.ContainsRune(Symbols, char) {
			hasSymbol = true
		}
	}
	
	if !hasUpper {
		t.Error("Password should contain uppercase letters")
	}
	if !hasLower {
		t.Error("Password should contain lowercase letters")
	}
	if !hasNumber {
		t.Error("Password should contain numbers")
	}
	if !hasSymbol {
		t.Error("Password should contain symbols")
	}
}

func TestGeneratePasswordCustomLength(t *testing.T) {
	config := &PasswordConfig{
		Length:     24,
		Uppercase:  true,
		Lowercase:  true,
		Numbers:    true,
		Symbols:    true,
	}
	
	password, err := GeneratePassword(config)
	if err != nil {
		t.Fatalf("GeneratePassword failed: %v", err)
	}
	
	if len(password) != 24 {
		t.Errorf("Expected password length 24, got %d", len(password))
	}
}

func TestGeneratePasswordOnlyUppercase(t *testing.T) {
	config := &PasswordConfig{
		Length:    12,
		Uppercase: true,
	}
	
	password, err := GeneratePassword(config)
	if err != nil {
		t.Fatalf("GeneratePassword failed: %v", err)
	}
	
	if len(password) != 12 {
		t.Errorf("Expected password length 12, got %d", len(password))
	}
	
	// Check that password contains only uppercase letters
	for _, char := range password {
		if !strings.ContainsRune(Uppercase, char) {
			t.Errorf("Password should only contain uppercase letters, found: %c", char)
		}
	}
}

func TestGeneratePasswordOnlyLowercase(t *testing.T) {
	config := &PasswordConfig{
		Length:    12,
		Lowercase: true,
	}
	
	password, err := GeneratePassword(config)
	if err != nil {
		t.Fatalf("GeneratePassword failed: %v", err)
	}
	
	if len(password) != 12 {
		t.Errorf("Expected password length 12, got %d", len(password))
	}
	
	// Check that password contains only lowercase letters
	for _, char := range password {
		if !strings.ContainsRune(Lowercase, char) {
			t.Errorf("Password should only contain lowercase letters, found: %c", char)
		}
	}
}

func TestGeneratePasswordOnlyNumbers(t *testing.T) {
	config := &PasswordConfig{
		Length:  12,
		Numbers: true,
	}
	
	password, err := GeneratePassword(config)
	if err != nil {
		t.Fatalf("GeneratePassword failed: %v", err)
	}
	
	if len(password) != 12 {
		t.Errorf("Expected password length 12, got %d", len(password))
	}
	
	// Check that password contains only numbers
	for _, char := range password {
		if !strings.ContainsRune(Numbers, char) {
			t.Errorf("Password should only contain numbers, found: %c", char)
		}
	}
}

func TestGeneratePasswordOnlySymbols(t *testing.T) {
	config := &PasswordConfig{
		Length:  12,
		Symbols: true,
	}
	
	password, err := GeneratePassword(config)
	if err != nil {
		t.Fatalf("GeneratePassword failed: %v", err)
	}
	
	if len(password) != 12 {
		t.Errorf("Expected password length 12, got %d", len(password))
	}
	
	// Check that password contains only symbols
	for _, char := range password {
		if !strings.ContainsRune(Symbols, char) {
			t.Errorf("Password should only contain symbols, found: %c", char)
		}
	}
}

func TestGeneratePasswordWithExclude(t *testing.T) {
	config := &PasswordConfig{
		Length:    16,
		Uppercase: true,
		Lowercase: true,
		Numbers:   true,
		Symbols:   true,
		Exclude:   "0O1lI",
	}
	
	password, err := GeneratePassword(config)
	if err != nil {
		t.Fatalf("GeneratePassword failed: %v", err)
	}
	
	// Check that excluded characters are not in the password
	for _, excludedChar := range config.Exclude {
		if strings.ContainsRune(password, excludedChar) {
			t.Errorf("Password should not contain excluded character: %c", excludedChar)
		}
	}
}

func TestGeneratePasswordNoRepeating(t *testing.T) {
	config := &PasswordConfig{
		Length:      16,
		Uppercase:   true,
		Lowercase:   true,
		Numbers:     true,
		Symbols:     true,
		NoRepeating: true,
	}
	
	password, err := GeneratePassword(config)
	if err != nil {
		t.Fatalf("GeneratePassword failed: %v", err)
	}
	
	// Check that no consecutive characters are the same
	for i := 1; i < len(password); i++ {
		if password[i] == password[i-1] {
			t.Errorf("Password should not have consecutive repeating characters at position %d: %c%c", i, password[i-1], password[i])
		}
	}
}

func TestValidateConfig(t *testing.T) {
	// Test valid config
	validConfig := &PasswordConfig{
		Length:    12,
		Uppercase: true,
		Lowercase: true,
	}
	
	if err := validateConfig(validConfig); err != nil {
		t.Errorf("Valid config should not return error: %v", err)
	}
	
	// Test too short length
	shortConfig := &PasswordConfig{
		Length:    6,
		Uppercase: true,
	}
	
	if err := validateConfig(shortConfig); err == nil {
		t.Error("Config with length < 8 should return error")
	}
	
	// Test too long length
	longConfig := &PasswordConfig{
		Length:    200,
		Uppercase: true,
	}
	
	if err := validateConfig(longConfig); err == nil {
		t.Error("Config with length > 128 should return error")
	}
	
	// Test no character sets selected
	noCharSetsConfig := &PasswordConfig{
		Length: 16,
	}
	
	if err := validateConfig(noCharSetsConfig); err == nil {
		t.Error("Config with no character sets should return error")
	}
}

func TestBuildCharSet(t *testing.T) {
	config := &PasswordConfig{
		Uppercase: true,
		Lowercase: true,
		Numbers:   true,
		Symbols:   true,
	}
	
	charSet := buildCharSet(config)
	expectedLength := len(Uppercase) + len(Lowercase) + len(Numbers) + len(Symbols)
	
	if len(charSet) != expectedLength {
		t.Errorf("Expected character set length %d, got %d", expectedLength, len(charSet))
	}
	
	// Test with exclude
	config.Exclude = "0O1lI"
	charSetWithExclude := buildCharSet(config)
	
	if len(charSetWithExclude) >= len(charSet) {
		t.Error("Character set with exclusions should be smaller")
	}
}

func TestAnalyzePasswordStrength(t *testing.T) {
	// Test strong password
	strongPassword := "MyP@ssw0rd123!"
	analysis := AnalyzePasswordStrength(strongPassword)
	
	if analysis["length"] != 16 {
		t.Errorf("Expected length 16, got %v", analysis["length"])
	}
	if !analysis["has_uppercase"].(bool) {
		t.Error("Expected has_uppercase to be true")
	}
	if !analysis["has_lowercase"].(bool) {
		t.Error("Expected has_lowercase to be true")
	}
	if !analysis["has_numbers"].(bool) {
		t.Error("Expected has_numbers to be true")
	}
	if !analysis["has_symbols"].(bool) {
		t.Error("Expected has_symbols to be true")
	}
	
	// Test weak password
	weakPassword := "123"
	weakAnalysis := AnalyzePasswordStrength(weakPassword)
	
	if weakAnalysis["strength_level"] != "Very Weak" {
		t.Errorf("Expected strength level 'Very Weak', got %v", weakAnalysis["strength_level"])
	}
	
	// Test empty password
	emptyAnalysis := AnalyzePasswordStrength("")
	if emptyAnalysis["strength_level"] != "Empty" {
		t.Errorf("Expected strength level 'Empty', got %v", emptyAnalysis["strength_level"])
	}
}

func TestRandomChar(t *testing.T) {
	charSet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	
	// Test multiple random characters
	seen := make(map[byte]bool)
	for i := 0; i < 100; i++ {
		char, err := randomChar(charSet)
		if err != nil {
			t.Fatalf("randomChar failed: %v", err)
		}
		
		if !strings.ContainsRune(charSet, rune(char)) {
			t.Errorf("Generated character %c is not in character set", char)
		}
		
		seen[char] = true
	}
	
	// With 100 attempts, we should see most characters
	if len(seen) < len(charSet)/2 {
		t.Error("Random character generation seems too deterministic")
	}
}

func TestRandomCharEmptySet(t *testing.T) {
	_, err := randomChar("")
	if err == nil {
		t.Error("Expected error when character set is empty")
	}
}

func TestShufflePassword(t *testing.T) {
	password := []byte("ABCDEFGHIJKLMNOP")
	original := make([]byte, len(password))
	copy(original, password)
	
	shufflePassword(password)
	
	// Check that password was shuffled (very unlikely to be the same)
	if string(password) == string(original) {
		t.Error("Password should be shuffled")
	}
	
	// Check that all characters are still present
	for _, char := range original {
		if !strings.ContainsRune(string(password), rune(char)) {
			t.Errorf("Character %c missing after shuffle", char)
		}
	}
}

func TestShuffleInts(t *testing.T) {
	ints := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	original := make([]int, len(ints))
	copy(original, ints)
	
	shuffleInts(ints)
	
	// Check that slice was shuffled (very unlikely to be the same)
	identical := true
	for i := range ints {
		if ints[i] != original[i] {
			identical = false
			break
		}
	}
	
	if identical {
		t.Error("Integer slice should be shuffled")
	}
	
	// Check that all values are still present
	seen := make(map[int]bool)
	for _, val := range ints {
		seen[val] = true
	}
	
	for _, val := range original {
		if !seen[val] {
			t.Errorf("Value %d missing after shuffle", val)
		}
	}
}
