# Desafio Client-Server-API

Neste desafio vamos aplicar o que aprendemos sobre webserver http, contextos,
banco de dados e manipulação de arquivos com Go.

## Estrutura do Projeto

O projeto consiste em dois sistemas em Go:

- **client.go**: Realiza uma requisição HTTP no server.go solicitando a cotação do dólar e salva o valor em um arquivo `cotacao.txt`.
- **server.go**: Consome a API de câmbio dólar-real e retorna o resultado em JSON para o cliente. Registra cada cotação no banco de dados SQLite usando contextos para gerenciar timeouts.

## Requisitos

1. **client.go**:
   - Realiza uma requisição HTTP para `/cotacao` no servidor.
   - Salva o valor da cotação em um arquivo `cotacao.txt`.
   - Usa contexto para garantir que a operação de recebimento não exceda 300ms.

2. **server.go**:
   - Expõe um endpoint `/cotacao` na porta 8080.
   - Consome a API `https://economia.awesomeapi.com.br/json/last/USD-BRL`.
   - Salva a cotação no banco de dados SQLite utilizando contexto para garantir:
     - Timeout máximo de 200ms na chamada da API externa.
     - Timeout máximo de 10ms na persistência dos dados no banco.

## Como Executar

Para executar os sistemas:

1. Certifique-se de ter Go instalado na sua máquina.
2. Clone o repositório.
3. Execute `go run client.go` para iniciar o cliente.
4. Execute `go run server.go` para iniciar o servidor.

## Dependências

- `github.com/mattn/go-sqlite3`: Driver SQLite para Go.