package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"shop-crud/item-service/config"
	authmiddle "shop-crud/item-service/middleware"
	"shop-crud/item-service/modules/handlers"
	"shop-crud/item-service/modules/repositories"
	"shop-crud/item-service/modules/usecases"
	"shop-crud/item-service/pkg/tracing"

	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	//	echojwt "github.com/labstack/echo-jwt/v4"
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
	   // Initialize tracing provider
   tp, err := tracing.InitTracerProvider("item-service", "tempo:4318")
   if err != nil {
       log.Fatal(err)
   }
   defer func() {
       if err := tp.Shutdown(context.Background()); err != nil {
           log.Fatal(err)
       }
   }()
   otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	// Inisialisasi koneksi database dari config package
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

	e.Use(otelecho.Middleware("item-service"))
	v1 := e.Group("/api/v1")

	// Inisialisasi repository, usecase, handler
	itemRepo := repositories.NewItemRepository(config.DBPool)
	itemUsecase := usecases.NewItemUsecase(itemRepo)
	itemHandler := handlers.NewItemHandler(itemUsecase)

	// Registrasi route dengan middleware JWT (opsional)
	// itemHandler.RegisterRoutes(v1, middleware.JWTWithConfig(middleware.JWTConfig{
	// 	SigningKey: []byte(jwtSecret),
	// }))
	// itemHandler.RegisterRoutes(v1, echojwt.WithConfig(echojwt.Config{
	// SigningKey: []byte(jwtSecret),
	// }))
	itemHandler.RegisterRoutes(v1, authmiddle.JWTAuthMiddleware(jwtSecret))


	// Jalankan server
	addr := fmt.Sprintf(":%s", appPort)
	log.Printf("✅ Item service berjalan di port %s", appPort)
	if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Gagal menjalankan server: %v", err)
	}
}
