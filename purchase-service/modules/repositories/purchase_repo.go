package repositories

import (
	"context"
	purchaseModels "purchase-service/modules/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
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
	tr := otel.Tracer("purchase-repository")
	ctx, span := tr.Start(ctx, "PurchaseRepository.CreatePurchaseInTx")
	defer span.End()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		span.RecordError(err)
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Insert ke tabel purchases
	purchaseQuery := `INSERT INTO purchases (id, user_id, total_amount, created_at) VALUES ($1, $2, $3, $4)`
	_, err = tx.Exec(ctx, purchaseQuery, purchase.ID, purchase.UserID, purchase.TotalAmount, purchase.CreatedAt)
	if err != nil {
		_, insertSpan := tr.Start(ctx, "InsertPurchase")
		insertSpan.End()
		span.RecordError(err)
		return err
	}

	itemQuery := `INSERT INTO purchase_items (id, purchase_id, item_id, quantity, price_at_purchase) VALUES ($1, $2, $3, $4, $5)`
	updateStockQuery := `UPDATE items SET stock = stock - $1 WHERE id = $2 AND stock >= $1`

	for _, item := range items {
		// 2. Sub-span untuk insert purchase item
		_, itemSpan := tr.Start(ctx, "InsertPurchaseItem")
		_, err = tx.Exec(ctx, itemQuery, uuid.New(), purchase.ID, item.ItemID, item.Quantity, item.PriceAtPurchase)
		itemSpan.End()
		if err != nil {
			span.RecordError(err)
			return err
		}

		// 3. Sub-span untuk update stock
		_, stockSpan := tr.Start(ctx, "UpdateItemStock")
		result, err := tx.Exec(ctx, updateStockQuery, item.Quantity, item.ItemID)
		stockSpan.End()
		if err != nil {
			span.RecordError(err)
			return err
		}
		if result.RowsAffected() == 0 {
			return pgx.ErrNoRows
		}
	}

	// 4. Commit transaction
	_, commitSpan := tr.Start(ctx, "CommitTransaction")
	err = tx.Commit(ctx)
	commitSpan.End()
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
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
