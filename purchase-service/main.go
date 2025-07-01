package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"purchase-service/config" 

	"purchase-service/modules/handlers"
	"purchase-service/modules/repositories"
	"purchase-service/modules/usecases"

	authmiddle "purchase-service/middleware"
	itemRepositories "shop-crud/item-service/modules/repositories"
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
	// Inisialisasi koneksi DB dari config
	config.InitDB()
	defer config.CloseDB()

	appPort := os.Getenv("PORT")
	if appPort == "" {
		appPort = "5002"
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

	// Init repo & usecase dengan shared DB
	purchaseRepo := repositories.NewPurchaseRepository(config.DBPool)
	itemRepo := itemRepositories.NewItemRepository(config.DBPool)
	purchaseUsecase := usecases.NewPurchaseUsecase(purchaseRepo, itemRepo)

	// Handler
	purchaseHandler := handlers.NewPurchaseHandler(purchaseUsecase)

	// Routes dengan JWT Middleware
	// purchaseHandler.RegisterRoutes(v1, middleware.JWTWithConfig(middleware.JWTConfig{
	// 	SigningKey: []byte(jwtSecret),
	// }))
	purchaseHandler.RegisterRoutes(v1, authmiddle.JWTAuthMiddleware(jwtSecret))


	// Start server
	addr := fmt.Sprintf(":%s", appPort)
	log.Printf("✅ Purchase service berjalan di port %s", appPort)
	if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Gagal menjalankan server: %v", err)
	}
}
