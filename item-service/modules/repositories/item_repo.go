package repositories

import (
	"context"
	"shop-crud/item-service/modules/models"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ItemRepository interface {
	Create(ctx context.Context, item *models.Item) error
	FindAll(ctx context.Context) ([]models.Item, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Item, error)
	Update(ctx context.Context, item *models.Item) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type itemRepository struct {
	db *pgxpool.Pool
}

func NewItemRepository(db *pgxpool.Pool) ItemRepository {
	return &itemRepository{db: db}
}

var tracer = otel.Tracer("item-service-repository")

func (r *itemRepository) Create(ctx context.Context, item *models.Item) error {
	ctx, span := tracer.Start(ctx, "ItemRepository.Create")
	defer span.End()

	span.SetAttributes(
		attribute.String("item.id", item.ID.String()),
		attribute.String("item.name", item.Name),
		attribute.Float64("item.price", item.Price),
	)

	query := `INSERT INTO items (id, name, description, price, stock, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(ctx, query, item.ID, item.Name, item.Description, item.Price, item.Stock, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		span.RecordError(err)
	}
	return err
}

func (r *itemRepository) FindAll(ctx context.Context) ([]models.Item, error) {
	ctx, span := tracer.Start(ctx, "ItemRepository.FindAll")
	defer span.End()

	var items []models.Item
	query := `SELECT id, name, description, price, stock, created_at, updated_at FROM items ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Price,
			&item.Stock,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			span.RecordError(err)
			return nil, err
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(attribute.Int("item.count", len(items)))
	return items, nil
}

func (r *itemRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	ctx, span := tracer.Start(ctx, "ItemRepository.FindByID")
	defer span.End()

	span.SetAttributes(attribute.String("item.id", id.String()))

	var item models.Item
	query := `SELECT id, name, description, price, stock, created_at, updated_at FROM items WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.Price,
		&item.Stock,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &item, nil
}

func (r *itemRepository) Update(ctx context.Context, item *models.Item) error {
	query := `UPDATE items SET name = $1, description = $2, price = $3, stock = $4, updated_at = $5 WHERE id = $6`
	_, err := r.db.Exec(ctx, query, item.Name, item.Description, item.Price, item.Stock, item.UpdatedAt, item.ID)
	return err
}

func (r *itemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM items WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}