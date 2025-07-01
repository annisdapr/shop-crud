# --- STAGE 1: BUILDER ---
FROM golang:1.23-alpine AS builder

# Set direktori kerja utama
WORKDIR /app

# Salin SEMUA file dari proyek Anda ke dalam container terlebih dahulu.
# Ini penting agar 'replace' di go.mod bisa menemukan path relatif seperti '../item-service'.
COPY . .

# ✅ DIUBAH: Pindah ke direktori layanan yang spesifik.
WORKDIR /app/user-service

# Jalankan 'go mod tidy' dari dalam direktori layanan.
# Ini akan mengunduh dependensi berdasarkan go.mod di folder ini.
RUN go mod tidy

# Build aplikasi dari direktori saat ini ('.').
# Go akan otomatis menemukan main.go di sini.
# Output binary diberi nama 'main' dan diletakkan di root direktori container.
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /main .

# --- STAGE 2: FINAL IMAGE ---
FROM alpine:latest
WORKDIR /app

# Salin HANYA binary yang sudah di-build dari stage builder.
COPY --from=builder /main .

EXPOSE 5000
CMD ["./main"]