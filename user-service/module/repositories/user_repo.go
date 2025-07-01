package repositories

import (
	"context"
	"user-service/module/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository mendefinisikan interface untuk operasi data user.
// Penggunaan interface memudahkan untuk testing (mocking).
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}

// Struct ini adalah implementasi konkret dari interface di atas.
type userRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository adalah constructor untuk membuat instance baru dari UserRepository.
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

// Create menyimpan user baru ke dalam database.
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (id, name, email, password_hash, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6)`
	// ExecContext digunakan untuk query yang tidak mengembalikan baris data (INSERT, UPDATE, DELETE).
	_, err := r.db.Exec(ctx, query, user.ID, user.Name, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt)
	return err
}

// FindByEmail mencari user berdasarkan alamat email.
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, name, email, password_hash, created_at, updated_at FROM users WHERE email = $1`
	// GetContext digunakan untuk query yang diharapkan mengembalikan satu baris data.
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	// Jika ada error (termasuk jika user tidak ditemukan), kembalikan error tersebut.
	if err != nil {
		return nil, err
	}
	return &user, nil
}
