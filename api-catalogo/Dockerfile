# Etapa 1: Compilación (Builder)
FROM golang:1.25.7-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
# Si el go.sum no está todavía, ejecutaremos go mod tidy después
# Para compilar, primero copiamos e intentamos
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/catalogo-service ./cmd/api/main.go

# Etapa 2: Imagen Final Minimalista (Alpine)
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/catalogo-service /app/catalogo-service
COPY --from=builder /app/db/migrations /app/db/migrations

EXPOSE 9094

ENV PORT=9094

CMD ["/app/catalogo-service"]
