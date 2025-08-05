// Package utils contém funções utilitárias que complementam o serviço de pseudonimização.
// Eu iniciei este pacote com ferramentas específicas para o contexto brasileiro,
// como a validação e geração de CPFs, que são dados pessoais frequentemente processados.
package utils

import (
	"crypto/rand"
	"fmt"
)

// IsValidCPF verifica se uma string corresponde a um número de CPF válido de acordo com
// o algoritmo oficial do Brasil. Esta função remove caracteres de formatação
// e valida os dígitos verificadores.
//
// Parâmetros:
//   - cpf: A string do CPF a ser validada (pode conter pontos e traço).
//
// Retorna:
//   - bool: `true` se o CPF for válido, `false` caso contrário.
func IsValidCPF(cpf string) bool {
	// Remove todos os caracteres que não são dígitos.
	cleaned := cleanCPF(cpf)

	// Verifica o comprimento (deve ter 11 dígitos).
	if len(cleaned) != 11 {
		return false
	}

	// Verifica padrões inválidos conhecidos (todos os dígitos iguais).
	if allDigitsSame(cleaned) {
		return false
	}

	// Calcula e compara os dígitos verificadores.
	firstDigit := calculateCPFCheckDigit(cleaned[:9], 10)
	secondDigit := calculateCPFCheckDigit(cleaned[:10], 11)

	return cleaned[9] == firstDigit && cleaned[10] == secondDigit
}

// GenerateSyntheticCPF cria um número de CPF sintético, porém válido, para uso em testes.
// O CPF gerado segue todas as regras de validação, mas eu o projetei para usar um prefixo
// conhecido (999) para indicar que não é um CPF real.
//
// Retorna:
//   - string: Um CPF sintético válido e formatado.
//   - error: Retorna um erro apenas se a geração de números aleatórios do sistema falhar.
func GenerateSyntheticCPF() (string, error) {
	// Uso o prefixo 999 para identificar claramente CPFs sintéticos.
	prefix := "999"

	// Gera 6 dígitos aleatórios para completar o corpo do CPF.
	randomDigits := make([]byte, 6)
	_, err := rand.Read(randomDigits)
	if err != nil {
		return "", fmt.Errorf("falha ao gerar dígitos aleatórios: %w", err)
	}

	for i := range randomDigits {
		randomDigits[i] = '0' + (randomDigits[i] % 10)
	}

	// Combina o prefixo e os dígitos aleatórios.
	partialCPF := prefix + string(randomDigits)

	// Calcula os dois dígitos verificadores.
	firstDigit := calculateCPFCheckDigit(partialCPF, 10)
	partialCPF += string(firstDigit)
	secondDigit := calculateCPFCheckDigit(partialCPF, 11)
	fullCPF := partialCPF + string(secondDigit)

	// Formata o resultado final com a pontuação padrão.
	return formatCPF(fullCPF), nil
}

// calculateCPFCheckDigit é uma função auxiliar para calcular um dígito verificador do CPF.
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

// cleanCPF é uma função auxiliar para remover todos os caracteres não numéricos de uma string.
func cleanCPF(cpf string) string {
	var cleaned []rune
	for _, c := range cpf {
		if c >= '0' && c <= '9' {
			cleaned = append(cleaned, c)
		}
	}
	return string(cleaned)
}

// allDigitsSame é uma função auxiliar para verificar se todos os dígitos da string são iguais.
func allDigitsSame(cpf string) bool {
	first := cpf[0]
	for i := 1; i < len(cpf); i++ {
		if cpf[i] != first {
			return false
		}
	}
	return true
}

// formatCPF é uma função auxiliar para formatar um CPF de 11 dígitos com a pontuação padrão.
func formatCPF(cpf string) string {
	if len(cpf) != 11 {
		return cpf
	}
	return fmt.Sprintf("%s.%s.%s-%s", cpf[:3], cpf[3:6], cpf[6:9], cpf[9:])
}
