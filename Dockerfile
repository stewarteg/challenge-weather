# Etapa 1: Build
FROM golang:1.23 AS builder

# Definir o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copiar os arquivos do projeto para o contêiner
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Compilar o binário de forma estática
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Etapa 2: Runtime
FROM debian:bullseye-slim

# Instalar os certificados CA
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Definir o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copiar o binário da etapa de build para a etapa de runtime
COPY --from=builder /app/main .

# Expor a porta que a aplicação utiliza
EXPOSE 8080

# Comando para executar a aplicação
CMD ["./main"]