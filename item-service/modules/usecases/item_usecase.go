package usecases

import (
	"context"
	"shop-crud/item-service/modules/models"
	"shop-crud/item-service/modules/repositories"
	"time"

	"github.com/google/uuid"
)

// ItemUsecase mendefinisikan logika bisnis untuk item.
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

// NewItemUsecase adalah constructor untuk usecase item.
func NewItemUsecase(itemRepo repositories.ItemRepository) ItemUsecase {
	return &itemUsecase{itemRepo: itemRepo}
}

func (u *itemUsecase) CreateItem(ctx context.Context, req models.CreateItemRequest) (*models.Item, error) {
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
		return nil, err
	}
	return newItem, nil
}

func (u *itemUsecase) GetAllItems(ctx context.Context) ([]models.Item, error) {
	return u.itemRepo.FindAll(ctx)
}

func (u *itemUsecase) GetItemByID(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	return u.itemRepo.FindByID(ctx, id)
}

func (u *itemUsecase) UpdateItem(ctx context.Context, id uuid.UUID, req models.UpdateItemRequest) (*models.Item, error) {
	// Pertama, dapatkan item yang ada untuk memastikan item tersebut ada
	existingItem, err := u.itemRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err // Akan mengembalikan error jika tidak ditemukan
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
	// Pastikan item ada sebelum menghapus
	_, err := u.itemRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return u.itemRepo.Delete(ctx, id)
}