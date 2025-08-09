package config

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DBPool *pgxpool.Pool

// InitDB initializes the database connection.
func InitDB() {
	cfg := GetConfig()
	dbURL := cfg.DBUrl // From Config struct

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Error parsing database URL: %v", err)
	}

	config.MaxConns = 10                  // Max allowed connections
	config.MinConns = 2                   // Minimum active connections
	config.HealthCheckPeriod = 1 * time.Minute // Check connection health every 1 minute

	// Create connection pool
	DBPool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Verify the connection with a ping
	if err = DBPool.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("✅ Database connected successfully!")
}

// CloseDB closes the database connection when the app stops.
func CloseDB() {
	if DBPool != nil {
		DBPool.Close()
		log.Println("🔌 Database connection closed.")
	}
}
