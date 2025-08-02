package main

import (
	"fmt"
	"log"

	"github.com/raywall/pseudonymization-lgpd-tools/utils"
)

func ExampleCPFFunctions() {
	// Validate a CPF
	valid := utils.IsValidCPF("529.982.247-25")
	fmt.Printf("Is valid CPF: %v\n", valid)

	// Generate a synthetic CPF
	synthetic, err := utils.GenerateSyntheticCPF()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Synthetic CPF: %s\n", synthetic)
	fmt.Printf("Is valid: %v\n", utils.IsValidCPF(synthetic))

	// Output:
	// Is valid CPF: true
	// Synthetic CPF: 999.XXX.XXX-XX
	// Is valid: true
}

func main() {
	ExampleCPFFunctions()
}
