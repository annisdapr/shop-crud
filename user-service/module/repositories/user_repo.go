package repositories

import (
	"context"
	"user-service/module/models"
	"user-service/pkg/logger"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/jackc/pgx"
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
   // start tracing for Create repository method
   tracer := otel.Tracer("user-service-repo")
   ctx, span := tracer.Start(ctx, "UserRepository.Create")
   defer span.End()
   span.SetAttributes(
       attribute.String("user.id", user.ID.String()),
       attribute.String("user.email", user.Email),
   )
   logger.Info(ctx, "💾 Creating user in DB: "+user.Email)
   query := `INSERT INTO users (id, name, email, password_hash, created_at, updated_at) 
          VALUES ($1, $2, $3, $4, $5, $6)`
  _, err := r.db.Exec(ctx, query, user.ID, user.Name, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		span.RecordError(err)
		logger.Error(ctx, "❌ Error creating user in DB: "+err.Error())
		return err
	}
	logger.Info(ctx, "✅ Successfully created user in DB: "+user.Email)
   return err
}

// FindByEmail mencari user berdasarkan alamat email.
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	tracer := otel.Tracer("user-service-repo")
	ctx, span := tracer.Start(ctx, "UserRepository.FindByEmail")
	defer span.End()
	var user models.User

	span.SetAttributes(
		attribute.String("user.email", email),
		attribute.String("db.statement", "SELECT ... FROM users WHERE email = $1"),
	)
	logger.Info(ctx, "🔍 Finding user by email in DB: "+email)
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
		span.RecordError(err)
		if err == pgx.ErrNoRows {
			// Menggunakan level WARN untuk "tidak ditemukan" agar mudah difilter
			logger.Warn(ctx, "❓ User not found in DB: "+email)
		} else {
			// Menggunakan level ERROR untuk masalah database lainnya
			logger.Error(ctx, "❌ Error finding user in DB: "+err.Error())
		}
		return nil, err
	}

	logger.Info(ctx, "✅ Successfully found user in DB: "+user.Email)
	return &user, nil
}