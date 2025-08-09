package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"shop-crud/item-service/config"
	"shop-crud/item-service/modules/handlers"
	"shop-crud/item-service/modules/repositories"
	"shop-crud/item-service/modules/usecases"
	 authmiddle"shop-crud/item-service/middleware"
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
	config.InitDB()
	defer config.CloseDB()

	appPort := os.Getenv("PORT")
	if appPort == "" {
		appPort = "5001"
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

	itemRepo := repositories.NewItemRepository(config.DBPool)
	itemUsecase := usecases.NewItemUsecase(itemRepo)
	itemHandler := handlers.NewItemHandler(itemUsecase)
	itemHandler.RegisterRoutes(v1, authmiddle.JWTAuthMiddleware(jwtSecret))

	addr := fmt.Sprintf(":%s", appPort)
	log.Printf("✅ Item service berjalan di port %s", appPort)
	if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Gagal menjalankan server: %v", err)
	}
}
