package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"user-service/config" 

	"user-service/module/handlers"
	"user-service/module/repositories"
	"user-service/module/usecases"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func main() {
	// Inisialisasi koneksi database dari config package
	config.InitDB()
	defer config.CloseDB()

	appPort := os.Getenv("PORT")
	if appPort == "" {
		appPort = "5000"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret"
	}

	// Setup Echo
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	v1 := e.Group("/api/v1")

	// Inisialisasi repository, usecase, dan handler
	userRepo := repositories.NewUserRepository(config.DBPool)
	userUsecase := usecases.NewUserUsecase(userRepo, jwtSecret)
	userHandler := handlers.NewUserHandler(userUsecase)
	userHandler.RegisterRoutes(v1)

	// Jalankan server
	addr := fmt.Sprintf(":%s", appPort)
	log.Printf("✅ User service berjalan di port %s", appPort)
	if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Gagal menjalankan server: %v", err)
	}
}
