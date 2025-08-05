// Copyright 2025 Raywall Malheiros
//
/*
Package pseudonymization fornece um conjunto de ferramentas que eu desenvolvi para facilitar a
pseudonimização de dados pessoais, um requisito fundamental para a conformidade com leis de
proteção de dados como a LGPD (Lei Geral de Proteção de Dados) do Brasil e a GDPR europeia.

# Visão Geral

Como desenvolvedor, meu objetivo foi criar uma solução prática que implementa as melhores
práticas de segurança para substituir dados sensíveis por pseudônimos, ao mesmo tempo que
permite a reversão (re-identificação) de forma controlada e segura.

O processo gera três artefatos principais:
  - Um hash SHA-256 do dado original, que pode ser usado para verificações sem expor o dado.
  - Um pseudônimo (UUID v4), que serve como o identificador substituto.
  - O valor original criptografado com AES-GCM, que permite a recuperação do dado original
    apenas por partes autorizadas que possuem a chave de criptografia.

# Principais Funcionalidades

- Hashing seguro com SHA-256 para referências não reversíveis.
- Criptografia simétrica com AES-GCM, um padrão moderno e seguro para criptografia autenticada.
- Geração de pseudônimos únicos usando UUID v4.
- Uma API simples e direta, agnóstica à camada de armazenamento ou transporte.

# Exemplo de Uso Básico

	// 1. Em produção, a chave DEVE vir de um sistema de gerenciamento seguro (ex: AWS Secrets Manager, Vault).
	key := make([]byte, 32) // Chave de 32 bytes para AES-256
	if _, err := rand.Read(key); err != nil {
		log.Fatal(err)
	}

	// 2. Crie uma instância do serviço de pseudonimização.
	svc := pseudonymization.NewService(key)

	// 3. Pseudonimize um dado sensível, como um CPF.
	resultado, err := svc.Pseudonymize("123.456.789-00", "processamento de folha de pagamento", "sistema_rh")
	if err != nil {
		log.Fatalf("Falha ao pseudonimizar: %v", err)
	}

	// Agora você pode armazenar o resultado de forma segura.
	// O pseudônimo pode ser usado como chave estrangeira, por exemplo.
	fmt.Printf("Pseudônimo gerado: %s\n", resultado.Pseudonym)
	fmt.Printf("Hash do original: %s\n", resultado.OriginalHash)

	// 4. Para reverter o processo (somente em contextos autorizados):
	valorOriginal, err := svc.Revert(resultado.EncryptedValue)
	if err != nil {
		log.Fatalf("Falha ao reverter: %v", err)
	}

	fmt.Printf("Valor original recuperado: %s\n", valorOriginal)

# Considerações de Segurança

A segurança de todo o processo depende da confidencialidade da chave de criptografia (`encryptionKey`).
Ela nunca deve ser armazenada em texto plano no código-fonte ou em sistemas de controle de versão.
*/
package pseudonymization
