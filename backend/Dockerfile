FROM golang:1.23-alpine

# Instalasi tool yang dibutuhkan: git (untuk go mod), make (untuk Makefile), migrate (untuk database)
RUN apk add --no-cache git make
RUN go install github.com/air-verse/air@latest
RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /app

# Download dependencies dulu agar build lebih cepat (caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy semua file ke dalam container
COPY . .

# Jalankan aplikasi dengan Air untuk development
CMD ["air", "-c", ".air.toml"]
