# --- STAGE 1: BUILDER ---
FROM golang:1.23-alpine AS builder

# Set main working directory
WORKDIR /app

# Copy ALL project files into the container.
# This ensures 'replace' in go.mod can resolve relative paths like '../item-service'.
COPY . .

# ✅ Changed: Move to the specific service directory
WORKDIR /app/user-service

# Run 'go mod tidy' to install dependencies for this service
RUN go mod tidy

# Build the app from the current directory ('.')
# Output binary named 'main' is placed in the container root
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /main .

# --- STAGE 2: FINAL IMAGE ---
FROM alpine:latest
WORKDIR /app

# Copy ONLY the built binary from the builder stage
COPY --from=builder /main .

# Expose service port
EXPOSE 5000

# Run the binary
CMD ["./main"]
