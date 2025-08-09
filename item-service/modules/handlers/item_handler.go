package handlers

import (
	"database/sql"
	"net/http"
	"shop-crud/item-service/modules/models"
	"shop-crud/item-service/modules/usecases"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("item-service-handler")
type ItemHandler struct {
	itemUsecase usecases.ItemUsecase
}

func NewItemHandler(itemUsecase usecases.ItemUsecase) *ItemHandler {
	return &ItemHandler{itemUsecase: itemUsecase}
}

func (h *ItemHandler) RegisterRoutes(router *echo.Group, authMiddleware echo.MiddlewareFunc) {
	itemGroup := router.Group("/items")
	
	// Endpoint no need authentication
	itemGroup.GET("", h.GetAllItems)
	itemGroup.GET("/:id", h.GetItemByID)

	// Endpoint REQUIRED authentication
	itemGroup.POST("", h.CreateItem, authMiddleware)
	itemGroup.PUT("/:id", h.UpdateItem, authMiddleware)
	itemGroup.DELETE("/:id", h.DeleteItem, authMiddleware)
}

func (h *ItemHandler) CreateItem(c echo.Context) error {
	ctx, span := tracer.Start(c.Request().Context(), "ItemHandler.CreateItem")
	defer span.End()

	var req models.CreateItemRequest
	if err := c.Bind(&req); err != nil {
		span.RecordError(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	if err := c.Validate(&req); err != nil {
		span.RecordError(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Add attribute tracing
	span.SetAttributes(
		attribute.String("item.name", req.Name),
		attribute.Float64("item.price", req.Price),
		attribute.Int("item.stock", req.Stock),
	)

	item, err := h.itemUsecase.CreateItem(ctx, req)
	if err != nil {
		span.RecordError(err)
		c.Logger().Errorf("Error creating item: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create item"})
	}

	return c.JSON(http.StatusCreated, item)
}

func (h *ItemHandler) GetAllItems(c echo.Context) error {
	ctx, span := tracer.Start(c.Request().Context(), "ItemHandler.GetAllItems")
	defer span.End()

	items, err := h.itemUsecase.GetAllItems(ctx)
	if err != nil {
		span.RecordError(err)
		c.Logger().Errorf("Error getting all items: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve items"})
	}
	return c.JSON(http.StatusOK, items)
}

func (h *ItemHandler) GetItemByID(c echo.Context) error {
	ctx, span := tracer.Start(c.Request().Context(), "ItemHandler.GetItemByID")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		span.RecordError(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid item ID"})
	}

	span.SetAttributes(attribute.String("item.id", id.String()))

	item, err := h.itemUsecase.GetItemByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Item not found"})
		}
		c.Logger().Errorf("Error getting item by id: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve item"})
	}
	return c.JSON(http.StatusOK, item)
}

func (h *ItemHandler) UpdateItem(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid item ID"})
	}
	
	var req models.UpdateItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	item, err := h.itemUsecase.UpdateItem(c.Request().Context(), id, req)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Item not found"})
		}
		c.Logger().Errorf("Error updating item: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update item"})
	}
	return c.JSON(http.StatusOK, item)
}

func (h *ItemHandler) DeleteItem(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid item ID"})
	}

	err = h.itemUsecase.DeleteItem(c.Request().Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Item not found"})
		}
		c.Logger().Errorf("Error deleting item: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete item"})
	}
	return c.NoContent(http.StatusNoContent)
}