# --- Estágio 1: Build ---
# Usamos uma imagem oficial do Go com Alpine Linux, que é menor.
# 'as builder' dá um nome a este estágio para que possamos nos referir a ele mais tarde.
FROM golang:1.22-alpine as builder

# Define o diretório de trabalho dentro do container.
WORKDIR /app

# Instala dependências necessárias para compilação
RUN apk add --no-cache git

# Copia os arquivos de gerenciamento de dependências primeiro.
# Isso aproveita o cache do Docker. Se esses arquivos não mudarem, o Docker não baixará as dependências novamente.
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o resto do código fonte.
COPY . .

# Compila a aplicação.
# CGO_ENABLED=0 desabilita o CGO para criar um binário estático.
# GOOS=linux garante que o binário seja compilado para o ambiente Linux do container final.
# -o /app/server define o nome e o local do arquivo de saída.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/server ./cmd/api/main.go

# --- Estágio 2: Final ---
# Começamos com uma imagem Alpine Linux limpa, que é extremamente leve (cerca de 5MB).
FROM alpine:latest

# Instala ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates

# Define o diretório de trabalho.
WORKDIR /app

# A mágica acontece aqui: copiamos APENAS o binário compilado do estágio 'builder'.
# A imagem final não terá o código fonte, nem as ferramentas do Go, apenas o executável.
COPY --from=builder /app/server .

# Também precisamos copiar nosso arquivo de configuração para que a aplicação possa lê-lo.
COPY ./config.yaml ./

# Cria um usuário não-root para segurança
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Muda a propriedade dos arquivos para o usuário não-root
RUN chown -R appuser:appgroup /app

# Muda para o usuário não-root
USER appuser

# Expõe a porta que nosso servidor Gin vai usar dentro do container.
EXPOSE 8080

# O comando para executar quando o container iniciar.
CMD ["./server"] 