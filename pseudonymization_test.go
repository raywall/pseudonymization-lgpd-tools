package pseudonymization

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPseudonymization(t *testing.T) {
	// Generate random encryption key
	key := make([]byte, 32)
	_, err := rand.Read(key)
	assert.NoError(t, err)

	svc := NewService(key)

	// Test pseudonymization
	value := "sensitive-data-123"
	result, err := svc.Pseudonymize(value, "test", "test")
	assert.NoError(t, err)
	assert.NotEmpty(t, result.OriginalHash)
	assert.NotEmpty(t, result.Pseudonym)
	assert.NotEmpty(t, result.EncryptedValue)
	assert.NotZero(t, result.Timestamp)

	// Test hash consistency
	hash1 := svc.Hash(value)
	hash2 := svc.Hash(value)
	assert.Equal(t, hash1, hash2)
	assert.Equal(t, result.OriginalHash, hash1)

	// Test revert
	original, err := svc.Revert(result.EncryptedValue)
	assert.NoError(t, err)
	assert.Equal(t, value, original)

	// Test empty value
	_, err = svc.Pseudonymize("", "test", "test")
	assert.Error(t, err)
}

func TestEncryption(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	assert.NoError(t, err)

	svc := NewService(key)

	plaintext := "test-value-456"
	encrypted, err := svc.encrypt(plaintext)
	assert.NoError(t, err)
	assert.NotEmpty(t, encrypted)
	assert.NotEqual(t, plaintext, encrypted)

	decrypted, err := svc.decrypt(encrypted)
	assert.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)

	// Test invalid ciphertext
	_, err = svc.decrypt("invalid-base64")
	assert.Error(t, err)
}
