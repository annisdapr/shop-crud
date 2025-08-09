package usecases

import (
	"context"
	"database/sql"
	"errors"

	//itemRepos "shop-crud/item-service/modules/repositories"
	"purchase-service/modules/clients"
	purchaseModels "purchase-service/modules/models"
	purchaseRepos "purchase-service/modules/repositories"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrStockNotSufficient = errors.New("stock for an item is not sufficient")
	ErrItemNotFound       = errors.New("one or more items not found")
)

type PurchaseUsecase interface {
	CreatePurchase(ctx context.Context, userID uuid.UUID, req purchaseModels.CreatePurchaseRequest) (*purchaseModels.Purchase, error)
	GetPurchaseHistory(ctx context.Context, userID uuid.UUID) ([]purchaseModels.Purchase, error)
}

type purchaseUsecase struct {
	purchaseRepo purchaseRepos.PurchaseRepository
	//itemRepo     itemRepos.ItemRepository
	itemClient   clients.ItemClient 
}

func NewPurchaseUsecase(purchaseRepo purchaseRepos.PurchaseRepository, itemClient clients.ItemClient) PurchaseUsecase {
	return &purchaseUsecase{
		purchaseRepo: purchaseRepo,
		itemClient:   itemClient,
	}
}

func (u *purchaseUsecase) CreatePurchase(ctx context.Context, userID uuid.UUID, req purchaseModels.CreatePurchaseRequest) (*purchaseModels.Purchase, error) {
	var totalAmount float64
	var purchaseItems []purchaseModels.PurchaseItem
	var purchaseItemResponses []purchaseModels.PurchaseItemResponse

	tr := otel.Tracer("purchase-usecase")
	ctx, span := tr.Start(ctx, "PurchaseUsecase.CreatePurchase")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.id", userID.String()),
		attribute.Int("item.count", len(req.Items)),
	)

	// Validasi dan kalkulasi total harga
	for _, reqItem := range req.Items {
		// Sub-span untuk GetItemByID
		itemSpanCtx, itemSpan := tr.Start(ctx, "ItemClient.GetItemByID",
			trace.WithAttributes(attribute.String("item.id", reqItem.ItemID.String())),
		)
		item, err := u.itemClient.GetItemByID(itemSpanCtx, reqItem.ItemID)
		itemSpan.End()

		if err != nil {
			if err == sql.ErrNoRows {
				return nil, ErrItemNotFound
			}
			return nil, err
		}
		if item.Stock < reqItem.Quantity {
			return nil, ErrStockNotSufficient
		}

		totalAmount += float64(reqItem.Quantity) * item.Price
		purchaseItems = append(purchaseItems, purchaseModels.PurchaseItem{
			ItemID:          reqItem.ItemID,
			Quantity:        reqItem.Quantity,
			PriceAtPurchase: item.Price,
		})
		purchaseItemResponses = append(purchaseItemResponses, purchaseModels.PurchaseItemResponse{
			ItemID:   reqItem.ItemID,
			Quantity: reqItem.Quantity,
			Name:     item.Name,
			Price:    item.Price,
		})
	}

	newPurchase := &purchaseModels.Purchase{
		ID:          uuid.New(),
		UserID:      userID,
		TotalAmount: totalAmount,
		CreatedAt:   time.Now(),
	}

	// Sub-span untuk penyimpanan DB
	dbSpanCtx, dbSpan := tr.Start(ctx, "PurchaseRepo.CreatePurchaseInTx")
	err := u.purchaseRepo.CreatePurchaseInTx(dbSpanCtx, newPurchase, purchaseItems)
	dbSpan.End()

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrStockNotSufficient
		}
		return nil, err
	}

	newPurchase.Items = purchaseItemResponses
	return newPurchase, nil
}


func (u *purchaseUsecase) GetPurchaseHistory(ctx context.Context, userID uuid.UUID) ([]purchaseModels.Purchase, error) {
	purchases, err := u.purchaseRepo.FindPurchasesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	for i, p := range purchases {
		items, err := u.purchaseRepo.FindPurchaseItemsByPurchaseID(ctx, p.ID)
		if err != nil {
			return nil, err
		}

		var itemResponses []purchaseModels.PurchaseItemResponse
		for _, item := range items {
			// Ambil detail item dari itemRepo
			itemDetail, err := u.itemClient.GetItemByID(ctx, item.ItemID)
			if err != nil {
				return nil, err
			}

			itemResponses = append(itemResponses, purchaseModels.PurchaseItemResponse{
				ItemID:   item.ItemID,
				Quantity: item.Quantity,
				Name:     itemDetail.Name,
				Price:    item.PriceAtPurchase,
			})
		}

		purchases[i].Items = itemResponses
	}

	return purchases, nil
}