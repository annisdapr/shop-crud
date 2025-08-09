# --- STAGE 1: BUILDER ---
FROM golang:1.23-alpine AS builder

# Install git and bash (some packages may require them during build)
RUN apk add --no-cache git bash

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum to cache dependency layers
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code (including app/ folder)
COPY . .

# Build the binary from main.go inside the app/ folder
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./app

# --- STAGE 2: FINAL IMAGE ---
FROM alpine:latest

# Copy the built binary from the builder stage to the final image
COPY --from=builder /app/main /app/main

# Set working directory
WORKDIR /app

# Document the exposed port (optional but useful)
EXPOSE 5000

# Run the binary
CMD ["./main"]
