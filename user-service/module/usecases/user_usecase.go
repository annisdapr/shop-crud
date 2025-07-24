package usecases

import (
   "context"
   "database/sql"
   "errors"
   "user-service/module/models"
   "user-service/module/repositories"
   "user-service/pkg/logger" 
   "time"

   "go.opentelemetry.io/otel"
   "go.opentelemetry.io/otel/attribute"

   "github.com/golang-jwt/jwt/v5"
   "github.com/google/uuid"
   "golang.org/x/crypto/bcrypt"
)

// Custom errors untuk penanganan yang lebih baik di layer handler.
var (
	ErrEmailExists      = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

// UserUsecase mendefinisikan logika bisnis untuk user.
type UserUsecase interface {
	Register(ctx context.Context, req models.RegisterRequest) (*models.User, error)
	Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error)
}

type userUsecase struct {
	userRepo  repositories.UserRepository
	jwtSecret string
}

// NewUserUsecase adalah constructor untuk usecase.
func NewUserUsecase(userRepo repositories.UserRepository, jwtSecret string) UserUsecase {
	return &userUsecase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Register menangani logika pendaftaran user baru.
func (u *userUsecase) Register(ctx context.Context, req models.RegisterRequest) (*models.User, error) {
   // start tracing for Register usecase
   tracer := otel.Tracer("user-service-usecase")
   ctx, span := tracer.Start(ctx, "UserUsecase.Register")
   defer span.End()
   // annotate span with request attributes
   span.SetAttributes(
       attribute.String("user.email", req.Email),
       attribute.String("user.name", req.Name),
   )
   	logger.Info(ctx, "🔍 Checking for existing email: "+req.Email)
   // 1. Cek duplikasi email.
	_, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err == nil { // Jika tidak ada error, berarti user ditemukan.
		logger.Warn(ctx, "⚠️ Email already registered: "+req.Email)
		return nil, ErrEmailExists
	}
	if !errors.Is(err, sql.ErrNoRows) { // Handle error database selain "tidak ditemukan".
		logger.	Error(ctx, "❌ Error checking email in database: "+err.Error())
		return nil, err
	}

	// 2. Hash password menggunakan bcrypt untuk keamanan.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(ctx, "❌ Error checking email in database: "+err.Error())
		return nil, err
	}

	// 3. Siapkan data user baru.
	newUser := &models.User{
		ID:           uuid.New(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 4. Simpan ke database via repository.
	if err := u.userRepo.Create(ctx, newUser); err != nil {
		logger.Error(ctx, "❌ Failed to save user to DB: "+err.Error())
		return nil, err
	}
	logger.Info(ctx, "✅ User registration successful: "+newUser.Email)
	return newUser, nil
}

// Login menangani logika otentikasi user dan pembuatan token.
func (u *userUsecase) Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {
   tracer := otel.Tracer("user-service-usecase")
   ctx, span := tracer.Start(ctx, "UserUsecase.Login")
   defer span.End()
   // annotate span with request attributes
   span.SetAttributes(
       attribute.String("user.email", req.Email),
   )
   	logger.Info(ctx, "🔐 Login attempt for: "+req.Email)
	// 1. Cari user berdasarkan email.
	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		// Samarkan error "tidak ditemukan" menjadi "kredensial salah" untuk keamanan.
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn(ctx, "⚠️ Email not found: "+req.Email)
			return nil, ErrInvalidCredentials
		}
		logger.Error(ctx, "❌ Failed to query email: "+err.Error())
		return nil, err
	}
	logger.Info(ctx, "🔍 Verifying password for: "+req.Email)
	// 2. Bandingkan password dari request dengan hash di database.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		logger.Error(ctx, "❌ Invalid password for: "+req.Email)
		return nil, ErrInvalidCredentials // Jika tidak cocok, kredensial salah.
	}

	logger.Info(ctx, "🔑 Generating JWT token for: "+req.Email)
	// 3. Jika cocok, buat klaim untuk JWT.
	claims := jwt.MapClaims{
		"sub":   user.ID, // Subject (standard claim), diisi user ID.
		"name":  user.Name,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(), // Token berlaku 72 jam.
		"iat":   time.Now().Unix(),                      // Issued At (standard claim).
	}

	// 4. Buat dan tandatangani token dengan secret key.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		logger.Error(ctx, "❌ Failed to sign JWT token: "+err.Error())
		return nil, err
	}

	logger.Info(ctx, "✅ Login successful for: "+req.Email)
	return &models.LoginResponse{AccessToken: tokenString}, nil
}
