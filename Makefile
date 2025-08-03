.PHONY: help build run test clean db-up db-down db-reset migrate generate

# Variáveis
BINARY_NAME=boilerplate-api
MAIN_PATH=cmd/api/main.go

# Comandos principais
help: ## Mostra esta ajuda
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Compila o projeto
	@echo "Compilando o projeto..."
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

run: ## Executa o servidor
	@echo "Executando o servidor..."
	go run $(MAIN_PATH)

dev: ## Executa o servidor em modo desenvolvimento
	@echo "Executando o servidor em modo desenvolvimento..."
	DB_HOST=localhost DB_PORT=5433 go run $(MAIN_PATH)

test: ## Executa os testes
	@echo "Executando testes..."
	go test -v ./...

test-integration: ## Executa os testes de integração
	@echo "Executando testes de integração..."
	go test -v ./tests/integration/...

test-coverage: ## Executa os testes com cobertura
	@echo "Executando testes com cobertura..."
	go test -v -cover ./...

test-all: ## Executa todos os testes (unitários + integração)
	@echo "Executando todos os testes..."
	go test -v ./...
	go test -v ./tests/integration/...

clean: ## Limpa arquivos de build
	@echo "Limpando arquivos de build..."
	rm -rf bin/
	go clean -cache

# Comandos do banco de dados
db-up: ## Inicia o banco de dados
	@echo "Iniciando banco de dados..."
	docker-compose up -d

db-down: ## Para o banco de dados
	@echo "Parando banco de dados..."
	docker-compose down

db-reset: ## Reseta o banco de dados (para e inicia)
	@echo "Resetando banco de dados..."
	docker-compose down -v
	docker-compose up -d

db-logs: ## Mostra logs do banco de dados
	@echo "Logs do banco de dados:"
	docker-compose logs -f postgres

# Comandos de migração e geração de código
migrate: ## Executa as migrações com goose
	@echo "Executando migrações com goose..."
	$(HOME)/go/bin/goose -dir sql/migrations postgres "host=localhost port=5433 user=postgres password=secret dbname=boilerplate sslmode=disable" up

migrate-status: ## Verifica status das migrações
	@echo "Verificando status das migrações..."
	$(HOME)/go/bin/goose -dir sql/migrations postgres "host=localhost port=5433 user=postgres password=secret dbname=boilerplate sslmode=disable" status

migrate-rollback: ## Faz rollback da última migração
	@echo "Fazendo rollback da última migração..."
	$(HOME)/go/bin/goose -dir sql/migrations postgres "host=localhost port=5433 user=postgres password=secret dbname=boilerplate sslmode=disable" down

generate: ## Gera código com sqlc
	@echo "Gerando código com sqlc..."
	sqlc generate

generate-docs: ## Gera documentação Swagger
	@echo "Gerando documentação Swagger..."
	$(HOME)/go/bin/swag init -g cmd/api/main.go -o docs

# Comandos de desenvolvimento completo
setup: ## Configura o ambiente de desenvolvimento
	@echo "Configurando ambiente de desenvolvimento..."
	go mod tidy
	docker-compose up -d
	@echo "Aguardando banco de dados..."
	sleep 5
	psql -h localhost -p 5433 -U postgres -d boilerplate -f sql/migrations/001_create_users_table.sql
	@echo "Ambiente configurado!"

start: ## Inicia o ambiente completo (banco + servidor)
	@echo "Iniciando ambiente completo..."
	docker-compose up -d
	@echo "✅ Ambiente iniciado! API disponível em http://localhost:8080"
	@echo "📊 Para ver logs: docker-compose logs -f api"

stop: ## Para o ambiente completo
	@echo "Parando ambiente completo..."
	docker-compose down

# Comandos de verificação
check: ## Verifica se tudo está funcionando
	@echo "Verificando ambiente..."
	@echo "1. Verificando se o banco está rodando..."
	@docker-compose ps | grep -q "Up" || (echo "Banco não está rodando!" && exit 1)
	@echo "2. Verificando se o código compila..."
	@go build $(MAIN_PATH) || (echo "Erro na compilação!" && exit 1)
	@echo "3. Verificando se os testes passam..."
	@go test ./... || (echo "Testes falharam!" && exit 1)
	@echo "✅ Tudo funcionando!"

# Comandos de documentação
docs: ## Gera documentação da API
	@echo "Gerando documentação da API..."
	@echo "Acesse: http://localhost:8080/api/v1/users"
	@echo "Health check: http://localhost:8080/health"

# Comandos de exemplo
example-create-user: ## Exemplo de criação de usuário
	@echo "Exemplo de criação de usuário:"
	@echo 'curl -X POST http://localhost:8080/api/v1/users \\'
	@echo '  -H "Content-Type: application/json" \\'
	@echo '  -d '"'"'{"email":"admin@example.com","password":"123456","name":"Admin","role":"admin"}'"'"''

example-login: ## Exemplo de login
	@echo "Exemplo de login:"
	@echo 'curl -X POST http://localhost:8080/api/v1/auth/login \\'
	@echo '  -H "Content-Type: application/json" \\'
	@echo '  -d '"'"'{"email":"admin@example.com","password":"123456"}'"'"'' 