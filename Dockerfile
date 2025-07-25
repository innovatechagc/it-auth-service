# Multi-stage build para optimizar el tamaño de la imagen
FROM golang:1.22-alpine AS builder

# Instalar dependencias necesarias
RUN apk add --no-cache git ca-certificates tzdata

# Crear directorio de trabajo
WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# Imagen final
FROM alpine:latest

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Crear usuario no-root
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Crear directorio para la aplicación
WORKDIR /app

# Copiar el binario desde el builder y establecer permisos
COPY --from=builder /app/main .
RUN chmod +x main && chown appuser:appgroup main

# Cambiar al usuario no-root
USER appuser

# Exponer puerto
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./main"]