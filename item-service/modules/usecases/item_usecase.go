package usecases

import (
	"context"
	"shop-crud/item-service/modules/models"
	"shop-crud/item-service/modules/repositories"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type ItemUsecase interface {
	CreateItem(ctx context.Context, req models.CreateItemRequest) (*models.Item, error)
	GetAllItems(ctx context.Context) ([]models.Item, error)
	GetItemByID(ctx context.Context, id uuid.UUID) (*models.Item, error)
	UpdateItem(ctx context.Context, id uuid.UUID, req models.UpdateItemRequest) (*models.Item, error)
	DeleteItem(ctx context.Context, id uuid.UUID) error
}

type itemUsecase struct {
	itemRepo repositories.ItemRepository
}

func NewItemUsecase(itemRepo repositories.ItemRepository) ItemUsecase {
	return &itemUsecase{itemRepo: itemRepo}
}

var tracer = otel.Tracer("item-service-usecase")

func (u *itemUsecase) CreateItem(ctx context.Context, req models.CreateItemRequest) (*models.Item, error) {
	ctx, span := tracer.Start(ctx, "ItemUsecase.CreateItem")
	defer span.End()

	span.SetAttributes(
		attribute.String("item.name", req.Name),
		attribute.String("item.description", req.Description),
		attribute.Float64("item.price", req.Price),
		attribute.Int("item.stock", req.Stock),
	)

	newItem := &models.Item{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := u.itemRepo.Create(ctx, newItem)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return newItem, nil
}

func (u *itemUsecase) GetAllItems(ctx context.Context) ([]models.Item, error) {
	ctx, span := tracer.Start(ctx, "ItemUsecase.GetAllItems")
	defer span.End()

	items, err := u.itemRepo.FindAll(ctx)
	if err != nil {
		span.RecordError(err)
	}
	return items, err
}

func (u *itemUsecase) GetItemByID(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	ctx, span := tracer.Start(ctx, "ItemUsecase.GetItemByID")
	defer span.End()

	span.SetAttributes(attribute.String("item.id", id.String()))

	item, err := u.itemRepo.FindByID(ctx, id)
	if err != nil {
		span.RecordError(err)
	}
	return item, err
}

func (u *itemUsecase) UpdateItem(ctx context.Context, id uuid.UUID, req models.UpdateItemRequest) (*models.Item, error) {
	existingItem, err := u.itemRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err 
	}

	// Update field
	existingItem.Name = req.Name
	existingItem.Description = req.Description
	existingItem.Price = req.Price
	existingItem.Stock = req.Stock
	existingItem.UpdatedAt = time.Now()

	err = u.itemRepo.Update(ctx, existingItem)
	if err != nil {
		return nil, err
	}
	return existingItem, nil
}

func (u *itemUsecase) DeleteItem(ctx context.Context, id uuid.UUID) error {
	_, err := u.itemRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return u.itemRepo.Delete(ctx, id)
}