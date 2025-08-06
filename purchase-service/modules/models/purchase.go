package models

import (
	"time"
	"github.com/google/uuid"
)


type Purchase struct {
	ID          uuid.UUID `db:"id" json:"id"`
	UserID      uuid.UUID `db:"user_id" json:"user_id"`
	TotalAmount float64   `db:"total_amount" json:"total_amount"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	Items       []PurchaseItemResponse `json:"items"` // Akan diisi oleh usecase
}


type PurchaseItem struct {
	ID                uuid.UUID `db:"id"`
	PurchaseID        uuid.UUID `db:"purchase_id"`
	ItemID            uuid.UUID `db:"item_id"`
	Quantity          int       `db:"quantity"`
	PriceAtPurchase   float64   `db:"price_at_purchase"`
}


type CreatePurchaseRequest struct {
	Items []PurchaseItemRequest `json:"items" validate:"required,min=1,dive"`
}


type PurchaseItemRequest struct {
	ItemID   uuid.UUID `json:"item_id" validate:"required"`
	Quantity int       `json:"quantity" validate:"required,gt=0"`
}


type PurchaseItemResponse struct {
	ItemID   uuid.UUID `json:"item_id"`
	Quantity int       `json:"quantity"`
	Name     string    `json:"name"`
	Price    float64   `json:"price"`
}

type PurchaseHistoryResponse struct {
	PurchaseID   uuid.UUID             `json:"purchase_id"`
	TotalAmount  float64               `json:"total_amount"`
	PurchasedAt  time.Time             `json:"purchased_at"`
	Items        []PurchaseItemHistory `json:"items"`
}

type PurchaseItemHistory struct {
	ItemID          uuid.UUID `json:"item_id"`
	Name            string    `json:"name"`
	Quantity        int       `json:"quantity"`
	PriceAtPurchase float64   `json:"price_at_purchase"`
	TotalPrice      float64   `json:"total_price"`
}