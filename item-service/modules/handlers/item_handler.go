package handlers

import (
	"database/sql"
	"net/http"
	"shop-crud/item-service/modules/models"
	"shop-crud/item-service/modules/usecases"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ItemHandler struct {
	itemUsecase usecases.ItemUsecase
}

func NewItemHandler(itemUsecase usecases.ItemUsecase) *ItemHandler {
	return &ItemHandler{itemUsecase: itemUsecase}
}

func (h *ItemHandler) RegisterRoutes(router *echo.Group, authMiddleware echo.MiddlewareFunc) {
	itemGroup := router.Group("/items")

	itemGroup.GET("", h.GetAllItems)
	itemGroup.GET("/:id", h.GetItemByID)

	itemGroup.POST("", h.CreateItem, authMiddleware)
	itemGroup.PUT("/:id", h.UpdateItem, authMiddleware)
	itemGroup.DELETE("/:id", h.DeleteItem, authMiddleware)
}

func (h *ItemHandler) CreateItem(c echo.Context) error {
	var req models.CreateItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	item, err := h.itemUsecase.CreateItem(c.Request().Context(), req)
	if err != nil {
		c.Logger().Errorf("Error creating item: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create item"})
	}

	return c.JSON(http.StatusCreated, item)
}

func (h *ItemHandler) GetAllItems(c echo.Context) error {
	items, err := h.itemUsecase.GetAllItems(c.Request().Context())
	if err != nil {
		c.Logger().Errorf("Error getting all items: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve items"})
	}
	return c.JSON(http.StatusOK, items)
}

func (h *ItemHandler) GetItemByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid item ID"})
	}

	item, err := h.itemUsecase.GetItemByID(c.Request().Context(), id)
	if err != nil {
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