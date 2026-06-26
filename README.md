# CashTrackr

CashTrackr é uma API REST desenvolvida em Go para gerenciamento de finanças pessoais.

O objetivo do projeto é servir como um estudo aprofundado de desenvolvimento backend utilizando Go, MongoDB e arquitetura em camadas, implementando boas práticas de organização, autenticação, testes e agregações de dados.

---

## Tecnologias

* Go
* Fiber
* MongoDB
* JWT
* bcrypt
* Testify
* Docker (ambiente de desenvolvimento)

---

## Arquitetura

O projeto segue uma arquitetura em camadas:

```
HTTP
   │
Handler
   │
Service
   │
Repository
   │
MongoDB
```

Cada camada possui responsabilidades bem definidas:

### Handler

Responsável por:

* Receber requisições HTTP
* Validar parâmetros e corpo da requisição
* Traduzir erros para códigos HTTP
* Chamar os serviços

### Service

Responsável pelas regras de negócio:

* Validações
* Fluxo da aplicação
* Comunicação entre repositórios

### Repository

Responsável pelo acesso aos dados:

* Consultas
* Inserções
* Atualizações
* Exclusões
* Aggregation Pipelines

---

# Funcionalidades

## Usuários

* Cadastro
* Login
* Alteração do nome
* Alteração da senha
* Autenticação via JWT

---

## Categorias

* Cadastro
* Listagem
* Busca por ID
* Atualização
* Exclusão
* Validação de nomes duplicados
* Proteção contra exclusão de categorias vinculadas a transações

---

## Transações

* Cadastro
* Listagem
* Busca por ID
* Atualização
* Exclusão
* Associação com categorias
* Valores monetários armazenados em centavos para evitar problemas de precisão

---

## Relatórios

### Resumo mensal

Retorna:

* Total de receitas
* Total de despesas
* Saldo do mês

---

### Gastos por categoria

Retorna:

* Total gasto por categoria
* Ordenação do maior para o menor gasto

---

### Evolução mensal

Retorna:

* Receitas por mês
* Despesas por mês
* Saldo mensal

---

# Estrutura do projeto

```
internal/
├── config/
├── dto/
├── handler/
├── middleware/
├── model/
├── repository/
├── routes/
├── service/
├── utils/
```

---

# Autenticação

A API utiliza JWT.

Após realizar o login, as rotas protegidas devem receber o token no header:

```
Authorization: Bearer <token>
```

---

# Persistência

Banco de dados:

* MongoDB

Coleções:

* users
* categories
* transactions

---

# Modelagem

## User

```
id
name
email
password_hash
created_at
updated_at
```

---

## Category

```
id
user_id
name
normalized_name
created_at
updated_at
```

---

## Transaction

```
id
user_id
category_id
title
description
type
amount
transaction_date
created_at
updated_at
```

O campo `amount` é armazenado em centavos (`int64`) para evitar problemas de precisão com números de ponto flutuante.

---

# Testes

O projeto possui testes unitários para a camada de serviços.

Cenários cobertos incluem:

* Criação de transações
* Validação de categorias
* Tratamento de erros de repositório
* CRUD de categorias

Os testes utilizam mocks para isolamento das dependências.

---

# Executando o projeto

## Clonar

```bash
git clone https://github.com/<seu-usuario>/cashtrackr.git
```

## Instalar dependências

```bash
go mod download
```

## Configurar variáveis de ambiente

Exemplo:

```
MONGO_URI=
DATABASE_NAME=
JWT_SECRET=
```

## Executar

```bash
go run cmd/api/main.go
```

---

# Executando os testes

Todos os testes:

```bash
go test ./...
```

Cobertura:

```bash
go test ./... -cover
```

Relatório HTML:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

# Próximas funcionalidades

* Contas bancárias
* Orçamentos mensais
* Dashboard consolidado
* Relatórios avançados
* Documentação OpenAPI/Swagger
* Testes de integração

---

# Objetivo do projeto

Este projeto foi desenvolvido como estudo de backend com foco em:

* Arquitetura em camadas
* Boas práticas em Go
* MongoDB
* APIs REST
* Testes
* Autenticação
* Aggregation Pipeline
* Organização e escalabilidade de código
