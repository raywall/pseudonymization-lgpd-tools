package pseudonymization

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Result represents the output of a pseudonymization operation
type Result struct {
	OriginalHash   string `json:"original_hash_value"`      // SHA-256 hash of original value (hex encoded)
	Pseudonym      string `json:"client_id"`                // Generated UUID v4 pseudonym
	EncryptedValue string `json:"encrypted_original_value"` // AES-GCM encrypted original value (base64 encoded)
	Timestamp      int64  `json:"anonymization_at"`         // Unix timestamp of operation
}

// Service provides pseudonymization methods
type Service struct {
	encryptionKey []byte
}

// NewService creates a new pseudonymization service instance
//
// Parameters:
//   - encryptionKey: 32-byte key for AES-256 encryption
//     In production, should come from secure key management
func NewService(encryptionKey []byte) *Service {
	return &Service{
		encryptionKey: encryptionKey,
	}
}

// Pseudonymize processes a sensitive value and returns pseudonymization artifacts
//
// Parameters:
// - value: The sensitive value to pseudonymize
// - purpose: Reason for pseudonymization (for audit trails)
// - system: Originating system (for audit trails)
//
// Returns:
// - Result containing pseudonymization artifacts
// - error if operation fails
func (s *Service) Pseudonymize(value, purpose, system string) (*Result, error) {
	if len(value) == 0 {
		return nil, errors.New("value cannot be empty")
	}

	// Generate SHA-256 hash of original value
	hash := sha256.Sum256([]byte(value))
	hashStr := hex.EncodeToString(hash[:])

	// Encrypt the original value
	encrypted, err := s.encrypt(value)
	if err != nil {
		return nil, fmt.Errorf("encryption failed: %w", err)
	}

	// Generate UUID v4 pseudonym
	pseudonym := uuid.New().String()

	return &Result{
		OriginalHash:   hashStr,
		Pseudonym:      pseudonym,
		EncryptedValue: encrypted,
		Timestamp:      time.Now().Unix(),
	}, nil
}

// Revert decrypts an encrypted value back to its original form
//
// Parameters:
// - encryptedValue: Base64-encoded encrypted value
//
// Returns:
// - Original plaintext value
// - error if decryption fails
func (s *Service) Revert(encryptedValue string) (string, error) {
	plaintext, err := s.decrypt(encryptedValue)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}
	return plaintext, nil
}

// Hash generates a SHA-256 hash of a value (hex encoded)
func (s *Service) Hash(value string) string {
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}

// encrypt performs AES-GCM encryption of plaintext
func (s *Service) encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decrypt performs AES-GCM decryption of ciphertext
func (s *Service) decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintextBytes, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintextBytes), nil
}
