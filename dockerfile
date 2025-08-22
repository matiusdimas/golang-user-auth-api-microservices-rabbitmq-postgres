# ---------- Build Stage ----------
FROM golang:1.23 AS builder

# Set working directory
WORKDIR /app

# Copy go.mod & go.sum dulu biar cache build lebih efisien
COPY go.mod go.sum ./
RUN go mod download

# Copy semua source code
COPY . .

# Build binary
RUN go build -o user-api cmd/main.go

# ---------- Runtime Stage ----------
FROM debian:bookworm-slim

WORKDIR /app

# Install CA certificates (biar bisa konek HTTPS/RabbitMQ)
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy binary dari builder
COPY --from=builder /app/user-api .

# Copy .env (opsional, bisa juga pakai docker-compose/env var langsung)
COPY .env .

# Expose port
EXPOSE 8080

# Jalankan binary
CMD ["./user-api"]
