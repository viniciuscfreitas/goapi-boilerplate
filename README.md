# GoAPI-Boilerplate

Um boilerplate robusto e profissional para desenvolvimento de APIs RESTful em Go, seguindo rigorosamente os princípios de Clean Architecture. Este projeto serve como base sólida para aplicações de produção, com configuração flexível, observabilidade estruturada e deploy containerizado.

![Go Version](https://img.shields.io/badge/Go-1.22%2B-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)
![Docker](https://img.shields.io/badge/Docker-Ready-blue?logo=docker)
![Architecture](https://img.shields.io/badge/Architecture-Clean%20Architecture-orange.svg)
![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL-blue?logo=postgresql)

## 🚀 Características

- **Clean Architecture**: Separação clara entre domínio, casos de uso e infraestrutura
- **Configuração Flexível**: Suporte a arquivos YAML e variáveis de ambiente
- **Observabilidade**: Logging estruturado JSON com request tracking
- **Containerização**: Docker multi-stage para imagens otimizadas
- **Banco de Dados**: PostgreSQL com sqlc para queries type-safe
- **Autenticação**: Sistema de usuários com bcrypt e roles
- **Testes**: Estrutura preparada para testes unitários e de integração

### Pré-requisitos

- Go 1.21+
- PostgreSQL 14+
- Docker e Docker Compose
- sqlc CLI

### Opção 1: Docker (Recomendado)

1. **Clone o repositório**
```bash
git clone <repository-url>
cd goapi-boilerplate
```

2. **Execute com Docker**
```bash
docker-compose up --build
```

3. **Teste a API**
```bash
curl http://localhost:8080/health
```

### Opção 2: Desenvolvimento Local

1. **Clone e configure**
```bash
git clone <repository-url>
cd goapi-boilerplate
go mod tidy
```

2. **Configure o ambiente**
```bash
# Copie o arquivo de configuração e ajuste as variáveis
cp config.yaml config.local.yaml
# Edite config.local.yaml com suas configurações locais
```

3. **Configure o banco de dados**
```bash
# Execute as migrações (usando goose - recomendado)
go install github.com/pressly/goose/v3/cmd/goose@latest
goose -dir sql/migrations postgres "host=localhost port=5432 user=postgres password=password dbname=boilerplate sslmode=disable" up

# Ou execute manualmente
psql -d your_database -f sql/migrations/001_create_users_table.sql
```

4. **Gere o código do sqlc**
```bash
sqlc generate
```

5. **Execute a aplicação**
```bash
go run cmd/api/main.go
```

## 📊 Endpoints da API

### Autenticação
- `POST /api/v1/auth/login` - Login de usuário

### Usuários
- `POST /api/v1/users` - Criar usuário
- `GET /api/v1/users` - Listar usuários (com paginação)
- `GET /api/v1/users/{id}` - Buscar usuário por ID
- `GET /api/v1/users/email?email=...` - Buscar usuário por email
- `PUT /api/v1/users/{id}` - Atualizar usuário
- `DELETE /api/v1/users/{id}` - Deletar usuário

### Sistema
- `GET /health` - Health check da API

## 📁 Estrutura do Projeto

```
├── cmd/
│   └── api/                    # Ponto de entrada da aplicação
├── internal/
│   ├── domain/                 # Camada de domínio (regras de negócio)
│   │   ├── user/              # Entidade User
│   │   └── repository/        # Interfaces dos repositórios
│   ├── usecase/               # Casos de uso (lógica de aplicação)
│   └── infrastructure/        # Camada de infraestrutura
│       ├── database/          # Código gerado pelo sqlc
│       ├── repository/        # Implementações dos repositórios
│       └── http/              # Handlers HTTP
├── sql/
│   ├── migrations/            # Migrações do banco de dados
│   └── queries/               # Queries SQL para sqlc
├── pkg/                       # Pacotes compartilhados
│   ├── logger/                # Sistema de logging
│   └── config/                # Configuração da aplicação
└── docs/                      # Documentação
```

## 🚀 Início Rápido

### Pré-requisitos

- Go 1.21+
- PostgreSQL 14+
- Docker e Docker Compose
- sqlc CLI

### Opção 1: Docker (Recomendado)

1. **Clone o repositório**
```bash
git clone <repository-url>
cd goapi-boilerplate
```

2. **Execute com Docker**
```bash
docker-compose up --build
```

3. **Teste a API**
```bash
curl http://localhost:8080/health
```

### Opção 2: Desenvolvimento Local

1. **Clone e configure**
```bash
git clone <repository-url>
cd goapi-boilerplate
go mod tidy
```

2. **Configure o banco de dados**
```bash
# Execute as migrações
psql -d your_database -f sql/migrations/001_create_users_table.sql
```

3. **Gere o código do sqlc**
```bash
sqlc generate
```

4. **Execute a aplicação**
```bash
go run cmd/api/main.go
```

## 🏗️ Arquitetura

### Princípios Fundamentais

1. **API-First**: Backend agnóstico que expõe API RESTful sem conhecimento sobre clientes
2. **Clean Architecture**: Lógica de negócio completamente desacoplada da infraestrutura
3. **Tipagem Forte**: Sistema de tipos do Go para máxima segurança
4. **Testabilidade**: Arquitetura que facilita testes unitários e de integração
5. **Observabilidade**: Logs estruturados para monitoramento em produção

### Stack Tecnológica

#### Backend (API RESTful)
- **Linguagem**: Go (Golang)
- **Framework Web**: Gin
- **Banco de Dados**: PostgreSQL
- **Comunicação com DB**: sqlc
- **Containerização**: Docker
- **Configuração**: Viper
- **Logging**: slog

## 📚 Documentação das Camadas

### Domain Layer (`internal/domain/`)

#### Entidade User (`internal/domain/user/user.go`)
```go
type User struct {
    ID        string    `json:"id"`
    Email     string    `json:"email"`
    Password  string    `json:"-"` // Não exposto na serialização
    Name      string    `json:"name"`
    Role      Role      `json:"role"`
    IsActive  bool      `json:"is_active"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

**Características**:
- Validação de dados
- Hash de senha com bcrypt
- Métodos de negócio (UpdateName, UpdateEmail, etc.)
- Tipagem forte com Role enum

#### Interface Repository (`internal/domain/repository/user_repository.go`)
```go
type UserRepository interface {
    Create(ctx context.Context, user *user.User) error
    GetByID(ctx context.Context, id string) (*user.User, error)
    GetByEmail(ctx context.Context, email string) (*user.User, error)
    Update(ctx context.Context, user *user.User) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, offset, limit int) ([]*user.User, error)
    Count(ctx context.Context) (int64, error)
    ExistsByEmail(ctx context.Context, email string) (bool, error)
    ExistsByID(ctx context.Context, id string) (bool, error)
}
```

### Use Case Layer (`internal/usecase/`)

#### UserUseCase (`internal/usecase/user_usecase.go`)
Implementa a lógica de aplicação com:
- Criação de usuários com validação
- Autenticação
- Atualização com verificações de unicidade
- Listagem com paginação
- Operações CRUD completas

### Infrastructure Layer (`internal/infrastructure/`)

#### PostgreSQL Repository (`internal/infrastructure/repository/postgres_user_repository.go`)
- Implementa a interface `UserRepository`
- Usa sqlc para queries type-safe
- Mapeamento entre domínio e banco de dados
- Tratamento de erros robusto

## 🧪 Testes

### Executar Testes
```bash
# Todos os testes
go test ./...

# Testes específicos
go test ./internal/domain/user/
go test ./internal/usecase/
go test ./internal/infrastructure/repository/
```

### Cobertura de Testes
```bash
go test -cover ./...
```

## ⚙️ Configuração

### Variáveis de Ambiente

O projeto suporta configuração via arquivo `config.yaml` ou variáveis de ambiente:

```bash
# Servidor
APP_SERVER_PORT=8080
APP_SERVER_HOST=0.0.0.0

# Banco de Dados
APP_DB_HOST=localhost
APP_DB_PORT=5432
APP_DB_USER=postgres
APP_DB_PASSWORD=secret
APP_DB_NAME=boilerplate

# Logging
APP_LOG_LEVEL=info
APP_ENV=development
```

### Docker

Para desenvolvimento com Docker:

```bash
# Build e execução
docker-compose up --build

# Apenas execução
docker-compose up

# Execução em background
docker-compose up -d

# Parar serviços
docker-compose down
```

## 📈 Logs e Observabilidade

O projeto utiliza logging estruturado JSON com slog:

```json
{
  "time": "2025-08-02T20:49:42.11966443Z",
  "level": "INFO",
  "msg": "Request handled",
  "request_id": "09cdfb13-21af-4801-9a1c-a91e436cedcd",
  "status_code": 201,
  "method": "POST",
  "path": "/api/v1/users",
  "latency_ms": 69,
  "user_agent": "curl/8.7.1"
}
```

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## 🆘 Suporte

Para dúvidas ou problemas, abra uma issue no repositório. 