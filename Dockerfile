# Estágio de compilação
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copia os arquivos de gerenciamento de dependências
COPY go.mod go.sum ./

# Baixa as dependências
RUN go mod download

# Copia o código-fonte
COPY . .

# Compila a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -o /locadora-api

# Estágio de produção
FROM alpine:latest

WORKDIR /root/

# Copia o binário compilado do estágio de compilação
COPY --from=builder /locadora-api .

# Expõe a porta em que a aplicação será executada
EXPOSE 8080

# Comando para executar a aplicação
CMD ["./locadora-api"]