# Etapa 1: Compilación (Builder)
FROM golang:1.25.7-alpine AS builder

# Instalar dependencias necesarias para compilar y certificados seguros
# hadolint ignore=DL3018
RUN apk add --no-cache tzdata

WORKDIR /app

# Descargar las dependencias del módulo Go
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente del microservicio
COPY . .

# Compilar un binario estático altamente optimizado
# -s -w elimina la tabla de símbolos y debuggers reduciendo hasta un 40% el peso final
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/auth-service ./cmd/api/main.go

# Etapa 2: Imagen Final Minimalista (Alpine)
FROM alpine:3.19

WORKDIR /app

# Copiar certificados de seguridad necesarios para solicitudes HTTPS salientes
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copiar el binario estático del builder
COPY --from=builder /app/auth-service /app/auth-service

# Copiar las migraciones de la base de datos para ejecución en runtime
COPY --from=builder /app/db/migrations /app/db/migrations

# Exponer el puerto por defecto expuesto
EXPOSE 9090

# Definir variables de entorno de producción por defecto
ENV PORT=9090

# Comando de ejecución
CMD ["/app/auth-service"]
