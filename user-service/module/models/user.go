package models

import (
	"time"
	"github.com/google/uuid"
)

// User merepresentasikan data pengguna di database.
// Tag `db` digunakan oleh sqlx untuk mapping, `json` oleh gin untuk response.
type User struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"` // Tanda `-` berarti jangan pernah kirim field ini dalam response JSON.
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// RegisterRequest adalah DTO (Data Transfer Object) untuk request registrasi.
// Tag `binding` digunakan oleh gin untuk validasi otomatis.
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"` // Password minimal 8 karakter.
}

// LoginRequest adalah DTO untuk request login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse adalah DTO untuk response login yang sukses.
type LoginResponse struct {
	AccessToken string `json:"access_token"`
}
