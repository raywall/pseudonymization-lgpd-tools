package main

import (
	"crypto/rand"
	"fmt"
	"log"

	"github.com/raywall/pseudonymization-lgpd-tools"
)

func ExampleService_Pseudonymize() {
	// Generate random encryption key (in production, use proper key management)
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Fatal("Failed to generate key:", err)
	}

	// Create pseudonymization service
	svc := pseudonymization.NewService(key)

	// Pseudonymize a sensitive value (e.g., CPF, email)
	value := "12345678901" // Example CPF
	result, err := svc.Pseudonymize(value, "data-processing", "analytics-system")
	if err != nil {
		log.Fatal("Pseudonymization failed:", err)
	}

	fmt.Printf("Original hash: %s\n", result.OriginalHash)
	fmt.Printf("Generated pseudonym: %s\n", result.Pseudonym)

	// Later, when you need the original value (authorized contexts only)
	original, err := svc.Revert(result.EncryptedValue)
	if err != nil {
		log.Fatal("Revert failed:", err)
	}

	fmt.Printf("Original value: %s\n", original)

	// Output:
	// Original hash: 15e24a16abfc4fef8b8b32bce3d1d06bfd274f3660793b754d730448d3b1e53f
	// Generated pseudonym: <uuid>
	// Original value: 12345678901
}

func main() {
	ExampleService_Pseudonymize()
}
