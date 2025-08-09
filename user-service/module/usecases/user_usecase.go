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

var (
	ErrEmailExists      = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type UserUsecase interface {
	Register(ctx context.Context, req models.RegisterRequest) (*models.User, error)
	Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error)
}

type userUsecase struct {
	userRepo  repositories.UserRepository
	jwtSecret string
}

func NewUserUsecase(userRepo repositories.UserRepository, jwtSecret string) UserUsecase {
	return &userUsecase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (u *userUsecase) Register(ctx context.Context, req models.RegisterRequest) (*models.User, error) {
   tracer := otel.Tracer("user-service-usecase")
   ctx, span := tracer.Start(ctx, "UserUsecase.Register")
   defer span.End()

   span.SetAttributes(
       attribute.String("user.email", req.Email),
       attribute.String("user.name", req.Name),
   )
   	logger.Info(ctx, "🔍 Checking for existing email: "+req.Email)
   // 1. Email duplication check 
	_, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err == nil { 
		logger.Warn(ctx, "⚠️ Email already registered: "+req.Email)
		return nil, ErrEmailExists
	}
	if !errors.Is(err, sql.ErrNoRows) { 
		logger.	Error(ctx, "❌ Error checking email in database: "+err.Error())
		return nil, err
	}

	// 2. Hash password 
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(ctx, "❌ Error checking email in database: "+err.Error())
		return nil, err
	}

	// 3. Prepare new user data
	newUser := &models.User{
		ID:           uuid.New(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 4. Save to repo
	if err := u.userRepo.Create(ctx, newUser); err != nil {
		logger.Error(ctx, "❌ Failed to save user to DB: "+err.Error())
		return nil, err
	}
	logger.Info(ctx, "✅ User registration successful: "+newUser.Email)
	return newUser, nil
}


func (u *userUsecase) Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {
   tracer := otel.Tracer("user-service-usecase")
   ctx, span := tracer.Start(ctx, "UserUsecase.Login")
   defer span.End()
   // annotate span with request attributes
   span.SetAttributes(
       attribute.String("user.email", req.Email),
   )
   	logger.Info(ctx, "🔐 Login attempt for: "+req.Email)
	// 1. Find user by email
	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn(ctx, "⚠️ Email not found: "+req.Email)
			return nil, ErrInvalidCredentials
		}
		logger.Error(ctx, "❌ Failed to query email: "+err.Error())
		return nil, err
	}
	logger.Info(ctx, "🔍 Verifying password for: "+req.Email)
	// 2. Password coomparison
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		logger.Error(ctx, "❌ Invalid password for: "+req.Email)
		return nil, ErrInvalidCredentials 
	}

	logger.Info(ctx, "🔑 Generating JWT token for: "+req.Email)
	// 3. If matched, claim JWT
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"name":  user.Name,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(), 
		"iat":   time.Now().Unix(),                     
	}

	// 4. Make and signed token with secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		logger.Error(ctx, "❌ Failed to sign JWT token: "+err.Error())
		return nil, err
	}

	logger.Info(ctx, "✅ Login successful for: "+req.Email)
	return &models.LoginResponse{AccessToken: tokenString}, nil
}
