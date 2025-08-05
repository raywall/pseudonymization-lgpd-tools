# Ferramentas de Pseudonimização para LGPD em Go

Olá! Eu sou Raywall Malheiros. Como desenvolvedor que atua no Brasil, sei da importância e dos desafios de estar em conformidade com a **Lei Geral de Proteção de Dados (LGPD)**. Eu criei esta biblioteca para oferecer um conjunto de ferramentas práticas e seguras que facilitam a implementação da **pseudonimização**, uma das técnicas recomendadas pela lei para a proteção de dados pessoais.

Meu objetivo foi criar uma solução de código aberto, simples e robusta para que outros desenvolvedores possam proteger os dados de seus usuários de forma eficaz.

## O que é Pseudonimização?

De acordo com a LGPD (Art. 5º, III), pseudonimização é "o tratamento por meio do qual um dado perde a possibilidade de associação, direta ou indireta, a um indivíduo, senão pelo uso de informação adicional mantida separadamente pelo controlador em ambiente controlado e seguro."

Em termos práticos, esta biblioteca faz o seguinte:

1.  **Substitui** um dado pessoal (como um CPF ou e-mail) por um **pseudônimo** (um identificador aleatório).
2.  **Mantém** uma forma segura de **reverter** o processo (re-identificar o indivíduo), mas apenas em um ambiente controlado e com a chave de segurança correta.

Isso permite que os sistemas continuem operando com os dados (usando os pseudônimos), minimizando o risco de exposição dos dados pessoais originais.

## O Processo de Pseudonimização da Biblioteca

Quando você usa o método `Pseudonymize`, ele gera três informações essenciais:

1.  `Pseudonym` **(Pseudônimo)**

    - **O que é?** Um UUID v4, que é um identificador universalmente único e aleatório.
    - **Para que serve?** Para substituir o dado original nos seus sistemas. Por exemplo, em vez de armazenar o CPF do usuário na tabela de pedidos, você armazena o pseudônimo dele.

2.  `OriginalHash` **(Hash do Original)**

    - **O que é?** Um hash criptográfico SHA-256 do dado original.
    - **Para que serve?** Para permitir buscas e verificações de unicidade sem expor o dado original. Como o hash é de mão única, você pode, por exemplo, verificar se um CPF já foi cadastrado comparando o hash do novo CPF com os hashes já armazenados.

3.  `EncryptedValue` **(Valor Criptografado)**
    - **O que é?** O dado original criptografado com o algoritmo seguro AES-GCM.
    - **Para que serve?** Esta é a "informação adicional mantida separadamente" que a LGPD menciona. Ela permite que, em contextos autorizados e com a posse da chave secreta, você possa reverter o processo e obter o dado original de volta.

## Como Usar

### Instalação

```bash
go get [github.com/raywall/pseudonymization-lgpd-tools](https://github.com/raywall/pseudonymization-lgpd-tools)
```

### Exemplo de Uso Principal

O fluxo básico envolve criar um `Service`, usar o método `Pseudonymize` para proteger os dados e, opcionalmente, o método `Revert` para recuperá-los.

```go
package main

import (
    "crypto/rand"
    "fmt"
    "log"

    "[github.com/raywall/pseudonymization-lgpd-tools](https://github.com/raywall/pseudonymization-lgpd-tools)"
)

func main() {
    // 1. GERE OU CARREGUE SUA CHAVE DE CRIPTOGRAFIA
    // A chave deve ter 32 bytes para AES-256.
    // EM PRODUÇÃO, NUNCA GERE A CHAVE DESTA FORMA. CARREGUE-A DE UM AMBIENTE SEGURO.
    key := make([]byte, 32)
    if _, err := rand.Read(key); err != nil {
        log.Fatal("Falha ao gerar a chave:", err)
    }

    // 2. CRIE UMA INSTÂNCIA DO SERVIÇO
    svc := pseudonymization.NewService(key)

    // 3. PSEUDONIMIZE UM DADO SENSÍVEL
    dadoSensivel := "123.456.789-00" // Exemplo: CPF
    resultado, err := svc.Pseudonymize(dadoSensivel, "processamento de analytics", "sistema_bi")
    if err != nil {
        log.Fatal("Falha na pseudonimização:", err)
    }

    // Você pode agora armazenar os artefatos gerados no seu banco de dados.
    // O `Pseudonym` pode ser usado em outras tabelas. O `OriginalHash` para buscas.
    fmt.Printf("Dado Original: %s\n", dadoSensivel)
    fmt.Printf("-> Pseudônimo: %s\n", resultado.Pseudonym)
    fmt.Printf("-> Hash: %s\n", resultado.OriginalHash)
    fmt.Printf("-> Valor Criptografado: %s\n", resultado.EncryptedValue)

    // 4. REVERTA O PROCESSO (somente quando necessário e autorizado)
    valorOriginal, err := svc.Revert(resultado.EncryptedValue)
    if err != nil {
        log.Fatal("Falha ao reverter o dado:", err)
    }

    fmt.Printf("\nValor original recuperado: %s\n", valorOriginal)
}
```

## Utilitários (`utils`)

Eu incluí um pacote `utils` com ferramentas úteis para o contexto brasileiro.

### Validação e Geração de CPF

Você pode validar CPFs ou gerar CPFs sintéticos (válidos, mas não reais) para usar em seus testes.

```go
package main

import (
    "fmt"
    "[github.com/raywall/pseudonymization-lgpd-tools/utils](https://github.com/raywall/pseudonymization-lgpd-tools/utils)"
)

func main() {
    // Validar um CPF
    cpfValido := "529.982.247-25"
    cpfInvalido := "111.111.111-11"
    fmt.Printf("CPF '%s' é válido? %v\n", cpfValido, utils.IsValidCPF(cpfValido))
    fmt.Printf("CPF '%s' é válido? %v\n", cpfInvalido, utils.IsValidCPF(cpfInvalido))

    // Gerar um CPF sintético para testes
    cpfSintetico, _ := utils.GenerateSyntheticCPF()
    fmt.Printf("\nCPF sintético gerado: %s\n", cpfSintetico)
    fmt.Printf("Ele é válido? %v\n", utils.IsValidCPF(cpfSintetico))
}
```

## ⚠️ Considerações de Segurança Críticas

A segurança de todo o processo de pseudonimização **depende inteiramente do sigilo da sua chave de criptografia (`encryptionKey`)**.

- **NUNCA** armazene a chave em texto plano no seu código-fonte ou em sistemas de controle de versão (como o Git).
- **NUNCA** use chaves fracas ou previsíveis. Use a função `crypto/rand` para gerar chaves fortes.
- **SEMPRE** carregue a chave de um ambiente seguro em produção. As melhores práticas incluem o uso de serviços como **AWS Secrets Manager**, **Google Secret Manager**, **HashiCorp Vault** ou, no mínimo, variáveis de ambiente seguras.

Se a sua chave for comprometida, a pseudonimização de todos os seus dados pode ser revertida por um ator malicioso.

---

Espero que esta biblioteca seja útil para seus projetos e ajude a construir aplicações mais seguras e em conformidade com a LGPD.

Atenciosamente,

**Raywall Malheiros**
