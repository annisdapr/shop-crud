package repositories

import (
	"context"
	"shop-crud/item-service/modules/models"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ItemRepository mendefinisikan interface untuk operasi data item.
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

// NewItemRepository adalah constructor untuk repository item.
func NewItemRepository(db *pgxpool.Pool) ItemRepository {
	return &itemRepository{db: db}
}

func (r *itemRepository) Create(ctx context.Context, item *models.Item) error {
	query := `INSERT INTO items (id, name, description, price, stock, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(ctx, query, item.ID, item.Name, item.Description, item.Price, item.Stock, item.CreatedAt, item.UpdatedAt)
	return err
}

func (r *itemRepository) FindAll(ctx context.Context) ([]models.Item, error) {
	var items []models.Item
	query := `SELECT id, name, description, price, stock, created_at, updated_at FROM items ORDER BY created_at DESC`

	// 1. Jalankan query. Ini mengembalikan 'rows' untuk diiterasi.
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	// 2. Pastikan untuk menutup rows setelah selesai. Ini sangat penting.
	defer rows.Close()

	// 3. Iterasi melalui setiap baris hasil query.
	for rows.Next() {
		var item models.Item
		// 4. Scan setiap kolom dari baris saat ini ke dalam struct 'item'.
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
			return nil, err
		}
		// 5. Tambahkan item yang sudah di-scan ke dalam slice.
		items = append(items, item)
	}

	// 6. Cek apakah ada error selama iterasi.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *itemRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	var item models.Item
	query := `SELECT id, name, description, price, stock, created_at, updated_at FROM items WHERE id = $1`
	
	// Pola ini sama dengan yang kita gunakan di user_repository.
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
		// pgx.ErrNoRows adalah error standar jika tidak ada baris yang ditemukan.
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