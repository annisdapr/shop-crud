# --- STAGE 1: BUILDER ---
FROM golang:1.23-alpine AS builder

# Install git dan bash karena beberapa package bisa butuh itu saat build
RUN apk add --no-cache git bash

# Set working directory
WORKDIR /app

# Salin go.mod dan go.sum untuk cache dependency layer
COPY go.mod go.sum ./
RUN go mod download

# Salin seluruh source code (termasuk folder app/)
COPY . .

# Build binary dari main.go di folder app/
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./app

# --- STAGE 2: FINAL IMAGE ---
FROM alpine:latest

# Menyalin binary hasil build ke image final
COPY --from=builder /app/main /app/main

# Set working directory
WORKDIR /app

# Dokumentasi port (tidak wajib, tapi membantu pembaca)
EXPOSE 5000

# Jalankan binary
CMD ["./main"]
