package handlers

import (
	"errors"
	"net/http"
	"user-service/module/models"
	"user-service/module/usecases"

	"github.com/labstack/echo/v4"
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
	
	// 1. Binding request body ke struct.
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// 2. Validasi struct menggunakan validator yang kita daftarkan di main.go.
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user, err := h.userUsecase.Register(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, usecases.ErrEmailExists) {
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()}) // 409 Conflict
		}
		// Kirim log error ke server untuk debugging (best practice)
		c.Logger().Errorf("Internal server error on register: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to register user"})
	}

	return c.JSON(http.StatusCreated, user) // 201 Created
}

// Login adalah handler untuk endpoint login, versi Echo.
func (h *UserHandler) Login(c echo.Context) error {
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	
	res, err := h.userUsecase.Login(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, usecases.ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()}) // 401 Unauthorized
		}
		c.Logger().Errorf("Internal server error on login: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to login"})
	}

	return c.JSON(http.StatusOK, res) // 200 OK
}

