package handlers

import (
//    "context"
   "errors"
   "net/http"
   "user-service/module/models"
   "user-service/module/usecases"
   "github.com/labstack/echo/v4"

   "go.opentelemetry.io/otel"
   "go.opentelemetry.io/otel/attribute"
//    "go.opentelemetry.io/otel/trace"
"user-service/pkg/logger" 
)

// UserHandler memegang dependency ke usecase. Strukturnya tetap sama.
type UserHandler struct {
	userUsecase usecases.UserUsecase
}

// NewUserHandler adalah constructor, tidak berubah.
func NewUserHandler(userUsecase usecases.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

// RegisterRoutes mendaftarkan semua endpoint yang berhubungan dengan user ke router Echo.
func (h *UserHandler) RegisterRoutes(router *echo.Group) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("/register", h.Register)
		userGroup.POST("/login", h.Login)
	}
}

// Register adalah handler untuk endpoint registrasi, versi Echo.
func (h *UserHandler) Register(c echo.Context) error {
   var req models.RegisterRequest
   // start tracing span for RegisterHandler
   tracer := otel.Tracer("user-service-handler")
   ctx, span := tracer.Start(c.Request().Context(), "RegisterHandler")
   defer span.End()
   // set route and email attributes after binding
	
	// 1. Binding request body ke struct.
   if err := c.Bind(&req); err != nil {
		logger.Error(ctx, "Failed to bind register request: "+err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
   // annotate span with attributes
   span.SetAttributes(
       attribute.String("http.route", c.Path()),
       attribute.String("user.email", req.Email),
   )

	// 2. Validasi struct menggunakan validator yang kita daftarkan di main.go.
	if err := c.Validate(&req); err != nil {
		logger.Error(ctx, "Validation failed on register: "+err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	logger.Info(ctx, "Attempting to register user: "+req.Email)
   user, err := h.userUsecase.Register(ctx, req)
	if err != nil {
		if errors.Is(err, usecases.ErrEmailExists) {
			logger.Warn(ctx, "Registration failed, email already exists: "+req.Email)
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()}) // 409 Conflict
		}
		logger.Error(ctx, "Internal server error on register: "+err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to register user"})
	}
	logger.Info(ctx, "✅ Register user success: "+user.Email)
	return c.JSON(http.StatusCreated, user) // 201 Created
}

// Login adalah handler untuk endpoint login, versi Echo.
func (h *UserHandler) Login(c echo.Context) error {
	var req models.LoginRequest
	tracer := otel.Tracer("user-service-handler")
	ctx, span := tracer.Start(c.Request().Context(), "LoginHandler")
	defer span.End()

	if err := c.Bind(&req); err != nil {
		logger.Error(ctx, "❌ Failed to bind login request: "+err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		logger.Error(ctx, "❌ Validation failed on login: "+err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	span.SetAttributes(
       attribute.String("http.route", c.Path()),
       attribute.String("user.email", req.Email),
	)
	res, err := h.userUsecase.Login(ctx, req)
	if err != nil {
		if errors.Is(err, usecases.ErrInvalidCredentials) {
			logger.Warn(ctx, "⚠️ Invalid login credentials for email: "+req.Email)
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()}) // 401 Unauthorized
		}
		logger.Error(ctx, "❌ Internal server error on login: "+err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to login"})
	}

	return c.JSON(http.StatusOK, res) // 200 OK
}

