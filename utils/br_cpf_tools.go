package utils

import (
	"crypto/rand"
	"fmt"
)

// IsValidCPF checks if a string is a valid CPF number according to Brazilian rules
// It removes formatting characters and validates the check digits
//
// Parameters:
// - cpf: The CPF string to validate (can include formatting like . and -)
//
// Returns:
// - bool: true if valid, false otherwise
func IsValidCPF(cpf string) bool {
	// Remove all non-digit characters
	cleaned := cleanCPF(cpf)

	// Check length (must be 11 digits)
	if len(cleaned) != 11 {
		return false
	}

	// Check for invalid patterns (all digits same)
	if allDigitsSame(cleaned) {
		return false
	}

	// Calculate first check digit
	firstDigit := calculateCPFCheckDigit(cleaned[:9], 10)

	// Calculate second check digit
	secondDigit := calculateCPFCheckDigit(cleaned[:10], 11)

	// Verify check digits
	return cleaned[9] == firstDigit && cleaned[10] == secondDigit
}

// GenerateSyntheticCPF creates a valid synthetic CPF for testing purposes
// The generated CPF follows the same validation rules as real CPFs but uses
// a known prefix to indicate it's synthetic (prefix 999)
//
// Returns:
// - string: A valid synthetic CPF (with formatting)
// - error: Only returns error if random number generation fails
func GenerateSyntheticCPF() (string, error) {
	// Use 999 as prefix to clearly identify synthetic CPFs
	prefix := "999"

	// Generate 6 random digits
	randomDigits := make([]byte, 6)
	_, err := rand.Read(randomDigits)
	if err != nil {
		return "", fmt.Errorf("failed to generate random digits: %w", err)
	}

	// Convert to digits 0-9
	for i := range randomDigits {
		randomDigits[i] = '0' + (randomDigits[i] % 10)
	}

	// Combine prefix and random digits (9 digits total)
	partialCPF := prefix + string(randomDigits)

	// Calculate first check digit
	firstDigit := calculateCPFCheckDigit(partialCPF, 10)
	partialCPF += string(firstDigit)

	// Calculate second check digit
	secondDigit := calculateCPFCheckDigit(partialCPF, 11)
	fullCPF := partialCPF + string(secondDigit)

	// Format with standard CPF punctuation
	return formatCPF(fullCPF), nil
}

// Helper function to calculate CPF check digit
func calculateCPFCheckDigit(partialCPF string, weight int) byte {
	var sum int
	for _, c := range partialCPF {
		sum += int(c-'0') * weight
		weight--
	}

	remainder := sum % 11
	if remainder < 2 {
		return '0'
	}
	return byte('0' + (11 - remainder))
}

// Helper function to remove all non-digit characters from CPF
func cleanCPF(cpf string) string {
	var cleaned []rune
	for _, c := range cpf {
		if c >= '0' && c <= '9' {
			cleaned = append(cleaned, c)
		}
	}
	return string(cleaned)
}

// Helper function to check if all digits are the same
func allDigitsSame(cpf string) bool {
	first := cpf[0]
	for i := 1; i < len(cpf); i++ {
		if cpf[i] != first {
			return false
		}
	}
	return true
}

// Helper function to format CPF with standard punctuation
func formatCPF(cpf string) string {
	if len(cpf) != 11 {
		return cpf
	}
	return fmt.Sprintf("%s.%s.%s-%s", cpf[:3], cpf[3:6], cpf[6:9], cpf[9:])
}
