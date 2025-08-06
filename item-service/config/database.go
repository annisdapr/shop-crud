package config

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DBPool *pgxpool.Pool

// InitDB menginisialisasi koneksi database.
func InitDB() {
	cfg := GetConfig()
	dbURL := cfg.DBUrl // Menggunakan nilai dari struct Config

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Error parsing database URL: %v", err)
	}

	config.MaxConns = 10                  // Maksimum koneksi yang diizinkan
	config.MinConns = 2                   // Minimum koneksi yang aktif
	config.HealthCheckPeriod = 1 * time.Minute // Mengecek kesehatan koneksi setiap 1 menit

	// Buat pool koneksi
	DBPool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Tidak dapat terhubung ke database: %v", err)
	}

	if err = DBPool.Ping(context.Background()); err != nil {
		log.Fatalf("Gagal melakukan ping ke database: %v", err)
	}

	log.Println("✅ Database berhasil terhubung!")
}

// CloseDB untuk menutup koneksi saat aplikasi berhenti.
func CloseDB() {
	if DBPool != nil {
		DBPool.Close()
		log.Println("🔌 Koneksi database ditutup.")
	}
}