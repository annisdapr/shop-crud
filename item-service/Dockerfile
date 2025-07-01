# --- STAGE 1: BUILDER ---
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .

# ✅ DIUBAH: Pindah ke direktori layanan yang spesifik.
WORKDIR /app/item-service

RUN go mod tidy
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /main .

# --- STAGE 2: FINAL IMAGE ---
FROM alpine:latest
WORKDIR /app
COPY --from=builder /main .
EXPOSE 5001
CMD ["./main"]