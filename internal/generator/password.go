package generator

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

// Character sets for password generation
const (
	Lowercase = "abcdefghijklmnopqrstuvwxyz"
	Uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numbers   = "0123456789"
	Symbols   = "!@#$%^&*()_+-=[]{}|;:,.<>?"
)

// PasswordConfig holds configuration for password generation
type PasswordConfig struct {
	Length     int
	Uppercase  bool
	Lowercase  bool
	Numbers    bool
	Symbols    bool
	Exclude    string // Characters to exclude
	NoRepeating bool  // Avoid consecutive repeating characters
}

// DefaultConfig returns a default password configuration
func DefaultConfig() *PasswordConfig {
	return &PasswordConfig{
		Length:     16,
		Uppercase:  true,
		Lowercase:  true,
		Numbers:    true,
		Symbols:    true,
		Exclude:    "",
		NoRepeating: true,
	}
}

// GeneratePassword creates a strong password based on the configuration
func GeneratePassword(config *PasswordConfig) (string, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// Validate configuration
	if err := validateConfig(config); err != nil {
		return "", fmt.Errorf("invalid configuration: %w", err)
	}

	// Build character set based on configuration
	charSet := buildCharSet(config)
	if len(charSet) == 0 {
		return "", fmt.Errorf("no character sets selected")
	}

	// Ensure minimum length
	if config.Length < len(charSet) {
		return "", fmt.Errorf("password length %d is too short for selected character sets", config.Length)
	}

	// Generate password
	password := make([]byte, config.Length)
	
	// First, ensure at least one character from each selected set
	password = ensureCharacterSets(password, config)
	
	// Fill remaining positions randomly
	for i := 0; i < config.Length; i++ {
		if password[i] == 0 {
			char, err := randomChar(charSet)
			if err != nil {
				return "", fmt.Errorf("failed to generate random character: %w", err)
			}
			password[i] = char
		}
	}

	// Apply no-repeating rule if enabled
	if config.NoRepeating {
		password = applyNoRepeatingRule(password)
	}

	// Shuffle the password to avoid predictable patterns
	shufflePassword(password)

	return string(password), nil
}

// validateConfig validates the password configuration
func validateConfig(config *PasswordConfig) error {
	if config.Length < 8 {
		return fmt.Errorf("password length must be at least 8 characters")
	}
	if config.Length > 128 {
		return fmt.Errorf("password length cannot exceed 128 characters")
	}
	
	// At least one character set must be selected
	if !config.Uppercase && !config.Lowercase && !config.Numbers && !config.Symbols {
		return fmt.Errorf("at least one character set must be selected")
	}
	
	return nil
}

// buildCharSet builds the character set based on configuration
func buildCharSet(config *PasswordConfig) string {
	var charSet strings.Builder
	
	if config.Uppercase {
		charSet.WriteString(Uppercase)
	}
	if config.Lowercase {
		charSet.WriteString(Lowercase)
	}
	if config.Numbers {
		charSet.WriteString(Numbers)
	}
	if config.Symbols {
		charSet.WriteString(Symbols)
	}
	
	// Remove excluded characters
	result := charSet.String()
	if config.Exclude != "" {
		for _, char := range config.Exclude {
			result = strings.ReplaceAll(result, string(char), "")
		}
	}
	
	return result
}

// ensureCharacterSets ensures at least one character from each selected set
func ensureCharacterSets(password []byte, config *PasswordConfig) []byte {
	positions := make([]int, 0, 4)
	
	// Collect available positions
	for i := range password {
		positions = append(positions, i)
	}
	
	// Shuffle positions to randomize placement
	shuffleInts(positions)
	posIndex := 0
	
	// Ensure uppercase if selected
	if config.Uppercase {
		if posIndex < len(positions) {
			char, _ := randomChar(Uppercase)
			password[positions[posIndex]] = char
			posIndex++
		}
	}
	
	// Ensure lowercase if selected
	if config.Lowercase {
		if posIndex < len(positions) {
			char, _ := randomChar(Lowercase)
			password[positions[posIndex]] = char
			posIndex++
		}
	}
	
	// Ensure numbers if selected
	if config.Numbers {
		if posIndex < len(positions) {
			char, _ := randomChar(Numbers)
			password[positions[posIndex]] = char
			posIndex++
		}
	}
	
	// Ensure symbols if selected
	if config.Symbols {
		if posIndex < len(positions) {
			char, _ := randomChar(Symbols)
			password[positions[posIndex]] = char
			posIndex++
		}
	}
	
	return password
}

// applyNoRepeatingRule ensures no consecutive repeating characters
func applyNoRepeatingRule(password []byte) []byte {
	charSet := buildCharSet(&PasswordConfig{
		Uppercase: true,
		Lowercase: true,
		Numbers:   true,
		Symbols:   true,
	})
	
	for i := 1; i < len(password); i++ {
		if password[i] == password[i-1] {
			// Find a different character
			for attempts := 0; attempts < 10; attempts++ {
				newChar, err := randomChar(charSet)
				if err == nil && newChar != password[i-1] {
					password[i] = newChar
					break
				}
			}
		}
	}
	
	return password
}

// randomChar selects a random character from the given character set
func randomChar(charSet string) (byte, error) {
	if len(charSet) == 0 {
		return 0, fmt.Errorf("empty character set")
	}
	
	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charSet))))
	if err != nil {
		return 0, fmt.Errorf("failed to generate random index: %w", err)
	}
	
	return charSet[index.Int64()], nil
}

// shufflePassword shuffles the password using Fisher-Yates algorithm
func shufflePassword(password []byte) {
	for i := len(password) - 1; i > 0; i-- {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			continue
		}
		password[i], password[j.Int64()] = password[j.Int64()], password[i]
	}
}

// shuffleInts shuffles a slice of integers
func shuffleInts(ints []int) {
	for i := len(ints) - 1; i > 0; i-- {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			continue
		}
		ints[i], ints[j.Int64()] = ints[j.Int64()], ints[i]
	}
}

// AnalyzePasswordStrength analyzes the strength of a password
func AnalyzePasswordStrength(password string) map[string]interface{} {
	analysis := map[string]interface{}{
		"length":        len(password),
		"has_uppercase": false,
		"has_lowercase": false,
		"has_numbers":   false,
		"has_symbols":   false,
		"unique_chars":  0,
		"strength_score": 0,
		"strength_level": "",
	}
	
	if len(password) == 0 {
		analysis["strength_level"] = "Empty"
		return analysis
	}
	
	// Check character types
	uniqueChars := make(map[rune]bool)
	for _, char := range password {
		uniqueChars[char] = true
		
		if char >= 'A' && char <= 'Z' {
			analysis["has_uppercase"] = true
		} else if char >= 'a' && char <= 'z' {
			analysis["has_lowercase"] = true
		} else if char >= '0' && char <= '9' {
			analysis["has_numbers"] = true
		} else {
			analysis["has_symbols"] = true
		}
	}
	
	analysis["unique_chars"] = len(uniqueChars)
	
	// Calculate strength score
	score := 0
	
	// Length contribution
	if len(password) >= 8 {
		score += 1
	}
	if len(password) >= 12 {
		score += 1
	}
	if len(password) >= 16 {
		score += 1
	}
	
	// Character variety contribution
	if analysis["has_uppercase"].(bool) {
		score += 1
	}
	if analysis["has_lowercase"].(bool) {
		score += 1
	}
	if analysis["has_numbers"].(bool) {
		score += 1
	}
	if analysis["has_symbols"].(bool) {
		score += 1
	}
	
	// Uniqueness contribution
	uniqueRatio := float64(analysis["unique_chars"].(int)) / float64(len(password))
	if uniqueRatio >= 0.8 {
		score += 1
	}
	
	analysis["strength_score"] = score
	
	// Determine strength level
	switch score {
	case 0, 1:
		analysis["strength_level"] = "Very Weak"
	case 2:
		analysis["strength_level"] = "Weak"
	case 3:
		analysis["strength_level"] = "Fair"
	case 4:
		analysis["strength_level"] = "Good"
	case 5:
		analysis["strength_level"] = "Strong"
	case 6, 7:
		analysis["strength_level"] = "Very Strong"
	default:
		analysis["strength_level"] = "Excellent"
	}
	
	return analysis
}
