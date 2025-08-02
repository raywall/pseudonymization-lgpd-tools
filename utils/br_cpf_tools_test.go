package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCPFValidation(t *testing.T) {
	testCases := []struct {
		cpf     string
		isValid bool
	}{
		{"529.982.247-25", true},  // Valid formatted CPF
		{"52998224725", true},     // Valid unformatted CPF
		{"111.111.111-11", false}, // Invalid (all same digits)
		{"123.456.789-00", false}, // Invalid (wrong check digits)
		{"529.982.247-26", false}, // Invalid (one wrong digit)
		{"", false},               // Empty
		{"123", false},            // Too short
		{"529982247252", false},   // Too long
	}

	for _, tc := range testCases {
		t.Run(tc.cpf, func(t *testing.T) {
			assert.Equal(t, tc.isValid, IsValidCPF(tc.cpf))
		})
	}
}

func TestSyntheticCPFGeneration(t *testing.T) {
	// Test multiple generations
	for i := 0; i < 100; i++ {
		cpf, err := GenerateSyntheticCPF()
		assert.NoError(t, err)
		assert.True(t, strings.HasPrefix(cleanCPF(cpf), "999"))
		assert.True(t, IsValidCPF(cpf))
	}
}
