package repositories

import (
	"context"
	purchaseModels "purchase-service/modules/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PurchaseRepository interface {
	// CreatePurchaseInTx menjalankan semua query dalam satu transaksi database.
	CreatePurchaseInTx(ctx context.Context, purchase *purchaseModels.Purchase, items []purchaseModels.PurchaseItem) error
	FindPurchasesByUserID(ctx context.Context, userID uuid.UUID) ([]purchaseModels.Purchase, error)
	FindPurchaseItemsByPurchaseID(ctx context.Context, purchaseID uuid.UUID) ([]purchaseModels.PurchaseItem, error)
}

type purchaseRepository struct {
	db *pgxpool.Pool
}

func NewPurchaseRepository(db *pgxpool.Pool) PurchaseRepository {
	return &purchaseRepository{db: db}
}

func (r *purchaseRepository) CreatePurchaseInTx(ctx context.Context, purchase *purchaseModels.Purchase, items []purchaseModels.PurchaseItem) error {
	// 1. Memulai transaksi menggunakan db.Begin() dari pgxpool
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	// Defer Rollback akan membatalkan transaksi jika terjadi panic atau error sebelum Commit.
	defer tx.Rollback(ctx)

	// 2. Buat record di tabel 'purchases' menggunakan tx.Exec()
	purchaseQuery := `INSERT INTO purchases (id, user_id, total_amount, created_at) VALUES ($1, $2, $3, $4)`
	_, err = tx.Exec(ctx, purchaseQuery, purchase.ID, purchase.UserID, purchase.TotalAmount, purchase.CreatedAt)
	if err != nil {
		return err
	}

	// 3. Loop melalui setiap item yang dibeli
	itemQuery := `INSERT INTO purchase_items (id, purchase_id, item_id, quantity, price_at_purchase) VALUES ($1, $2, $3, $4, $5)`
	updateStockQuery := `UPDATE items SET stock = stock - $1 WHERE id = $2 AND stock >= $1`

	for _, item := range items {
		// 3a. Buat record di 'purchase_items'
		_, err = tx.Exec(ctx, itemQuery, uuid.New(), purchase.ID, item.ItemID, item.Quantity, item.PriceAtPurchase)
		if err != nil {
			return err
		}

		// 3b. Update stok di tabel 'items'.
		result, err := tx.Exec(ctx, updateStockQuery, item.Quantity, item.ItemID)
		if err != nil {
			return err
		}
		// Cek apakah ada baris yang terpengaruh untuk memastikan stok cukup.
		if result.RowsAffected() == 0 {
			// Jika tidak ada baris yang diupdate, berarti stok tidak cukup.
			// Kembalikan pgx.ErrNoRows agar bisa ditangani di usecase.
			return pgx.ErrNoRows
		}
	}

	// 4. Jika semua query berhasil, commit transaksi.
	return tx.Commit(ctx)
}

func (r *purchaseRepository) FindPurchasesByUserID(ctx context.Context, userID uuid.UUID) ([]purchaseModels.Purchase, error) {
	var purchases []purchaseModels.Purchase
	query := `SELECT id, user_id, total_amount, created_at FROM purchases WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p purchaseModels.Purchase
		err := rows.Scan(&p.ID, &p.UserID, &p.TotalAmount, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		purchases = append(purchases, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return purchases, nil
}

func (r *purchaseRepository) FindPurchaseItemsByPurchaseID(ctx context.Context, purchaseID uuid.UUID) ([]purchaseModels.PurchaseItem, error) {
	var items []purchaseModels.PurchaseItem
	query := `SELECT id, purchase_id, item_id, quantity, price_at_purchase FROM purchase_items WHERE purchase_id = $1`

	rows, err := r.db.Query(ctx, query, purchaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var i purchaseModels.PurchaseItem
		err := rows.Scan(&i.ID, &i.PurchaseID, &i.ItemID, &i.Quantity, &i.PriceAtPurchase)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
