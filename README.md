# CashTrackr

API REST para gerenciamento de finanças pessoais desenvolvida em Go.

O objetivo do projeto é servir como estudo e demonstração de boas práticas de desenvolvimento backend, incluindo autenticação, arquitetura em camadas, persistência de dados, mensageria, cache e observabilidade.

## Tecnologias

- Go
- Fiber
- MongoDB
- JWT
- Docker (em desenvolvimento)
- Redis (planejado)
- NATS (planejado)

## Funcionalidades

### Usuários

- Criar usuário
- Login
- Consultar usuário autenticado
- Alterar nome
- Alterar senha

### Transações

- Criar transação
- Consultar transação por id
- Consultar transações paginadas por filtro
- Atualizar transação
- Apagar transação

### Categorias

- Criar categoria
- Consultar categoria por id
- Consultar categorias paginadas por filtro
- Atualizar categoria
- Apagar categoria

## Decisões técnicas

### Valores monetários

Os valores financeiros são armazenados como `int64`, representando centavos, para evitar problemas de precisão associados ao uso de `float`.

## Arquitetura

O projeto segue uma arquitetura em camadas:

```text
Handler
  ↓
Service
  ↓
Repository
  ↓
MongoDB
```

### Responsabilidades

- Handler: HTTP, requests e responses
- Service: regras de negócio
- Repository: persistência de dados

## Estrutura de pastas

```text
cmd/
internal/
├── app_error/
├── auth/
├── database/
├── dto/
├── handler/
├── middleware/
├── model/
├── repository/
└── service/
```

## Como executar

### Pré-requisitos

- Go 1.24+
- MongoDB

### Instalação

```bash
git clone https://github.com/seu-usuario/cashtrackr.git

cd cashtrackr

go mod download
```

### Variáveis de ambiente

```env
MONGODB_URI=mongodb://localhost:27017
DATABASE_NAME=cashtrackr_dev

JWT_SECRET=your-secret
JWT_EXPIRATION=10 // tempo em minutos
```

### Executar

```bash
go run cmd/api/main.go
```

## Próximos passos

- [x] Cadastro de usuários
- [x] Login JWT
- [x] Alteração de senha
- [x] Criação de transações
- [x] Listagem de transações
- [x] Atualização de transações
- [x] Exclusão de transações
- [x] Paginação
- [x] Categorias
- [ ] Relatórios
- [ ] Redis
- [ ] NATS
- [ ] Docker
- [ ] Testes automatizados