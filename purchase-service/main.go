package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"purchase-service/config"

	"purchase-service/modules/handlers"
	"purchase-service/modules/repositories"
	"purchase-service/modules/usecases"
	"purchase-service/pkg/tracing"

	authmiddle "purchase-service/middleware"
	//itemRepositories "shop-crud/item-service/modules/repositories"
	"purchase-service/modules/clients"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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
   tp, err := tracing.InitTracerProvider("purchase-service", "tempo:4318")
   if err != nil {
       log.Fatal(err)
   }
   defer func() {
       if err := tp.Shutdown(context.Background()); err != nil {
           log.Fatal(err)
       }
   }()
   otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

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
	e := echo.New()
	e.Use(otelecho.Middleware("purchase-service"))
	// Setup Echo

	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	v1 := e.Group("/api/v1")

	// Init repo & usecase dengan shared DB
	purchaseRepo := repositories.NewPurchaseRepository(config.DBPool)
	//itemRepo := itemRepositories.NewItemRepository(config.DBPool)
	//purchaseUsecase := usecases.NewPurchaseUsecase(purchaseRepo, itemRepo)
	itemClient := clients.NewItemClient("http://item-service:5001/api/v1")
	purchaseUsecase := usecases.NewPurchaseUsecase(purchaseRepo, itemClient)
	// Handler
	purchaseHandler := handlers.NewPurchaseHandler(purchaseUsecase)

	purchaseHandler.RegisterRoutes(v1, authmiddle.JWTAuthMiddleware(jwtSecret))


	// Start server
	addr := fmt.Sprintf(":%s", appPort)
	log.Printf("✅ Purchase service berjalan di port %s", appPort)
	if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Gagal menjalankan server: %v", err)
	}
}
