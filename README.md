# Pseudonymization Library

A language-agnostic pseudonymization library for data protection compliance (LGPD/GDPR).

## Features

- Secure one-way hashing (SHA-256)
- Reversible encryption (AES-256-GCM)
- UUID v4 pseudonym generation
- Storage/transport layer agnostic
- Compliance with data protection regulations

## Installation

```bash
go get github.com/raywall/pseudonymization-lgpd-tools
```

## Usage

### Basic Pseudonymization

```go
package main

import (
	"crypto/rand"
	"fmt"
	"log"

	"github.com/yourorg/pseudonymization"
)

func main() {
	// Generate encryption key (use proper key management in production)
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		log.Fatal(err)
	}

	svc := pseudonymization.NewService(key)

	// Pseudonymize sensitive data
	result, err := svc.Pseudonymize("sensitive-value", "purpose", "system")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Pseudonym: %s\n", result.Pseudonym)
	fmt.Printf("Original hash: %s\n", result.OriginalHash)

	// Revert when needed
	original, err := svc.Revert(result.EncryptedValue)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Original: %s\n", original)
}
```

## Security Considerations

- Always use proper key management (HSM/KMS) in production
- Store encryption keys separately from pseudonymized data
- Implement proper access controls for reverting pseudonymization
- Audit all pseudonymization/reversion operations

## Compliance

This library helps implement the following data protection principles:

- Data minimization: Only process necessary data
- Storage limitation: Original data is encrypted
- Integrity and confidentiality: Strong encryption standards
- Accountability: Audit trails through operation metadata

## License

MIT
