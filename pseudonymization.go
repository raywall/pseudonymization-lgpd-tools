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

// Result representa a saída de uma operação de pseudonimização.
// Eu projetei esta estrutura para conter todos os artefatos necessários para
// o uso seguro e a eventual reversão de um dado pseudonimizado.
type Result struct {
	// OriginalHash é o hash SHA-256 do valor original. Serve como uma impressão digital
	// para buscas e verificações sem expor o dado original.
	OriginalHash string `json:"original_hash_value"`

	// Pseudonym é um UUID v4 gerado para substituir o dado original. Este é o valor
	// que deve ser usado na maioria dos contextos de processamento de dados.
	Pseudonym string `json:"client_id"`

	// EncryptedValue é o valor original criptografado com AES-GCM. Este artefato
	// permite a reversão (re-identificação) do dado, mas apenas por quem possui a chave secreta.
	EncryptedValue string `json:"encrypted_original_value"`

	// Timestamp é o registro de quando a operação de pseudonimização ocorreu (em Unix timestamp).
	Timestamp int64 `json:"anonymization_at"`
}

// Service fornece os métodos para pseudonimização e reversão de dados.
// Ele encapsula a chave de criptografia para garantir que as operações sejam realizadas de forma segura.
type Service struct {
	encryptionKey []byte
}

// NewService cria uma nova instância do serviço de pseudonimização.
// A chave de criptografia é o componente mais crítico para a segurança.
//
// Parâmetros:
//   - encryptionKey: Uma chave de 32 bytes para usar o algoritmo AES-256.
//     Em produção, esta chave deve ser carregada de um local seguro, como um
//     cofre de segredos (AWS Secrets Manager, HashiCorp Vault, etc.).
func NewService(encryptionKey []byte) *Service {
	return &Service{
		encryptionKey: encryptionKey,
	}
}

// Pseudonymize processa um valor sensível e retorna os artefatos de pseudonimização.
// Este é o método central para proteger um dado pessoal.
//
// Parâmetros:
//   - value: O dado sensível a ser pseudonimizado (ex: CPF, email, nome completo).
//   - purpose: A finalidade da pseudonimização (para trilhas de auditoria).
//   - system: O sistema de origem que solicitou a operação (para trilhas de auditoria).
//
// Retorna:
//   - Um ponteiro para a struct Result contendo os artefatos da pseudonimização.
//   - Um erro se a operação falhar.
func (s *Service) Pseudonymize(value, purpose, system string) (*Result, error) {
	if len(value) == 0 {
		return nil, errors.New("o valor a ser pseudonimizado não pode estar vazio")
	}

	// 1. Gera o hash SHA-256 do valor original para servir como referência segura.
	hash := sha256.Sum256([]byte(value))
	hashStr := hex.EncodeToString(hash[:])

	// 2. Criptografa o valor original para permitir a reversão controlada.
	encrypted, err := s.encrypt(value)
	if err != nil {
		return nil, fmt.Errorf("a criptografia do valor original falhou: %w", err)
	}

	// 3. Gera um pseudônimo (UUID v4) para substituir o valor original.
	pseudonym := uuid.New().String()

	return &Result{
		OriginalHash:   hashStr,
		Pseudonym:      pseudonym,
		EncryptedValue: encrypted,
		Timestamp:      time.Now().Unix(),
	}, nil
}

// Revert descriptografa um valor de volta à sua forma original.
// Este método deve ser usado com extremo cuidado e apenas em contextos
// onde a re-identificação do titular dos dados é legalmente permitida e necessária.
//
// Parâmetros:
//   - encryptedValue: O valor criptografado (codificado em base64) obtido do Result.
//
// Retorna:
//   - O valor original em texto plano.
//   - Um erro se a descriptografia falhar (ex: chave incorreta, dado corrompido).
func (s *Service) Revert(encryptedValue string) (string, error) {
	plaintext, err := s.decrypt(encryptedValue)
	if err != nil {
		return "", fmt.Errorf("a descriptografia do valor falhou: %w", err)
	}
	return plaintext, nil
}

// Hash gera um hash SHA-256 de um valor e o codifica como uma string hexadecimal.
// É uma função de conveniência para gerar hashes consistentes com o processo de pseudonimização.
func (s *Service) Hash(value string) string {
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}

// encrypt é a função interna que realiza a criptografia AES-GCM.
// AES-GCM é um modo de criptografia de bloco autenticada, que garante tanto
// a confidencialidade quanto a integridade do dado.
func (s *Service) encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// O nonce (number used once) deve ser único para cada criptografia com a mesma chave.
	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return "", err
	}

	// Seal anexa o nonce ao início do ciphertext.
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decrypt é a função interna que realiza a descriptografia AES-GCM.
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
		return "", errors.New("o dado criptografado é muito curto")
	}

	// Extrai o nonce do início do dado.
	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintextBytes, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		// Este erro ocorre se a autenticação falhar (dado corrompido ou chave errada).
		return "", err
	}

	return string(plaintextBytes), nil
}
