package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"user-service/config"
	"user-service/pkg/tracing"

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
   tp, err := tracing.InitTracerProvider("user-service", "tempo:4318")
   if err != nil {
       log.Fatal(err)
   }
   defer func() {
       if err := tp.Shutdown(context.Background()); err != nil {
           log.Fatal(err)
       }
   }()

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

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	v1 := e.Group("/api/v1")

	userRepo := repositories.NewUserRepository(config.DBPool)
	userUsecase := usecases.NewUserUsecase(userRepo, jwtSecret)
	userHandler := handlers.NewUserHandler(userUsecase)
	userHandler.RegisterRoutes(v1)

	addr := fmt.Sprintf(":%s", appPort)
	log.Printf("✅ User service berjalan di port %s", appPort)
	if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Gagal menjalankan server: %v", err)
	}
}
