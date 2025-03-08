# Etapa 1: Build
FROM golang:1.20 AS builder

# Definir o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copiar os arquivos do projeto para o contêiner
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compilar o binário da aplicação
RUN go build -o main .

# Etapa 2: Runtime
FROM debian:bullseye-slim

# Definir o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copiar o binário da etapa de build para a etapa de runtime
COPY --from=builder /app/main .

# Expor a porta que a aplicação utiliza
EXPOSE 8080

# Comando para executar a aplicação
CMD ["./main"]