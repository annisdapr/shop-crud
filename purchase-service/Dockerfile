# --- STAGE 1: BUILDER ---
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Copy all project files
COPY . .

# ✅ Changed: Switch to the purchase-service directory
WORKDIR /app/purchase-service

# Download and tidy dependencies
RUN go mod tidy

# Build the binary with optimizations
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /main .

# --- STAGE 2: FINAL IMAGE ---
FROM alpine:latest
WORKDIR /app

# Copy only the compiled binary
COPY --from=builder /main .

# Expose service port
EXPOSE 5002

# Run the binary
CMD ["./main"]
