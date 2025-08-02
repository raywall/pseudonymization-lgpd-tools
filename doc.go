// Package pseudonymization provides tools for secure data pseudonymization
// in compliance with data protection regulations like LGPD and GDPR.
//
// Overview:
//
// This package implements best practices for pseudonymizing sensitive personal data,
// specifically focusing on:
// - Secure one-way hashing for identification
// - Reversible encryption of original values
// - UUID generation as pseudonyms
//
// Key Features:
// - SHA-256 hashing for secure reference
// - AES-GCM encryption for reversible pseudonymization
// - UUID v4 generation as pseudonyms
// - Complete agnosticism of storage/transport layer
//
// Basic Usage Example:
//
//	// Create a new service with encryption key
//	key := make([]byte, 32) // In production, use proper key management
//	svc := pseudonymization.NewService(key)
//
//	// Pseudonymize a value (e.g., CPF)
//	result, err := svc.Pseudonymize("12345678901", "purpose", "system")
//	if err != nil {
//	    log.Fatalf("Failed to pseudonymize: %v", err)
//	}
//
//	fmt.Printf("Generated pseudonym: %s\n", result.Pseudonym)
//
//	// Revert pseudonymization (when authorized)
//	original, err := svc.Revert(result.EncryptedValue)
//	if err != nil {
//	    log.Fatalf("Failed to revert: %v", err)
//	}
//
//	fmt.Printf("Original value: %s\n", original)
package pseudonymization
