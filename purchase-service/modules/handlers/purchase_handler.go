package handlers

import (
	"errors"
	"net/http"
	"purchase-service/middleware"
	purchaseModels "purchase-service/modules/models"
	purchaseUsecases "purchase-service/modules/usecases"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PurchaseHandler struct {
	purchaseUsecase purchaseUsecases.PurchaseUsecase
}

func NewPurchaseHandler(purchaseUsecase purchaseUsecases.PurchaseUsecase) *PurchaseHandler {
	return &PurchaseHandler{purchaseUsecase: purchaseUsecase}
}

func (h *PurchaseHandler) RegisterRoutes(router *echo.Group, authMiddleware echo.MiddlewareFunc) {
	purchaseGroup := router.Group("/purchases", authMiddleware) // Semua endpoint di sini terproteksi
	{
		purchaseGroup.POST("", h.CreatePurchase)
		// purchaseGroup.GET("", h.GetHistory) // Bisa ditambahkan nanti
	}
}

func (h *PurchaseHandler) CreatePurchase(c echo.Context) error {
	// 1. Ambil data user dari context yang sudah di-set oleh middleware
	claims, ok := middleware.GetUserFromContext(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token claims"})
	}
	userID, err := uuid.Parse(claims["sub"].(string))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID in token"})
	}

	// 2. Bind dan validasi request body
	var req purchaseModels.CreatePurchaseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// 3. Panggil usecase
	purchase, err := h.purchaseUsecase.CreatePurchase(c.Request().Context(), userID, req)
	if err != nil {
		if errors.Is(err, purchaseUsecases.ErrItemNotFound) || errors.Is(err, purchaseUsecases.ErrStockNotSufficient) {
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}
		c.Logger().Errorf("Error creating purchase: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create purchase"})
	}

	return c.JSON(http.StatusCreated, purchase)
}
