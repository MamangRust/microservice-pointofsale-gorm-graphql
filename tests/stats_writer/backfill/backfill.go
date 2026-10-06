// Package backfill implements the stats-writer `backfill` command: it reads
// historical OLTP rows from the POS PostgreSQL database (orders, order_items
// joined with products for category_id, transactions joined with orders for
// cashier_id) and materializes them into ClickHouse through the same batch
// repository used for live events.
//
// This is the bootstrap path for the stats pipeline — it lets the ClickHouse
// tables reflect pre-existing data without replaying every domain event.
package backfill

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/microservice-point-of-sale-pkg/database"
	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/events"
	"github.com/MamangRust/microservice-point-of-sale-stats-writer/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"go.uber.org/zap"
)

// backfillEventID derives a deterministic UUID per entity so re-running the
// backfill replaces the same ReplacingMergeTree key (with a newer version)
// instead of appending duplicates.
func backfillEventID(kind string, id int32) string {
	return uuid.NewSHA1(uuid.NameSpaceDNS, []byte(fmt.Sprintf("backfill:%s:%d", kind, id))).String()
}

// Backfiller reads OLTP rows and pushes them into ClickHouse. POS runs all
// services against one PostgreSQL database, so a single pool covers every
// stats source.
type Backfiller struct {
	log  logger.LoggerInterface
	repo repository.Repository
	db   *gorm.DB
}

// New opens the POS database connection and returns a ready Backfiller. Call
// Close to release it.
func New(log logger.LoggerInterface, repo repository.Repository) (*Backfiller, error) {
	db, err := database.NewGormClient(log)
	if err != nil {
		return nil, fmt.Errorf("connect POS database: %w", err)
	}
	return &Backfiller{log: log, repo: repo, db: db}, nil
}

func (b *Backfiller) Close() {
	if b.db != nil {
		if sqlDB, err := b.db.DB(); err == nil {
			sqlDB.Close()
		}
	}
}

// Run streams all stats sources into ClickHouse. The event version is the
// backfill run timestamp so re-running supersedes previous rows.
func (b *Backfiller) Run(ctx context.Context) error {
	version := uint64(time.Now().Unix())
	counts := map[string]int{}

	if err := b.backfillOrders(ctx, version, counts); err != nil {
		return err
	}
	if err := b.backfillOrderItems(ctx, version, counts); err != nil {
		return err
	}
	if err := b.backfillTransactions(ctx, version, counts); err != nil {
		return err
	}

	if err := b.repo.Flush(ctx); err != nil {
		return fmt.Errorf("flush backfill batches: %w", err)
	}

	b.log.Info("backfill complete",
		zap.Int("orders", counts["order"]),
		zap.Int("order_items", counts["order_item"]),
		zap.Int("transactions", counts["transaction"]),
	)
	return nil
}

type orderRow struct {
	OrderID    int32
	MerchantID int32
	CashierID  int32
	TotalPrice int64
	CreatedAt  time.Time
}

type itemRow struct {
	OrderItemID int32
	OrderID     int32
	ProductID   int32
	CategoryID  int32
	Quantity    int32
	UnitPrice   int32
	CreatedAt   time.Time
}

type txRow struct {
	TransactionID int32
	OrderID       int32
	MerchantID    int32
	CashierID     int32
	PaymentMethod string
	Amount        int32
	Status        string
	CreatedAt     time.Time
}

func (b *Backfiller) backfillOrders(ctx context.Context, version uint64, counts map[string]int) error {
	rows, err := b.db.WithContext(ctx).Raw(`SELECT order_id, merchant_id, cashier_id, total_price, created_at
		FROM orders WHERE deleted_at IS NULL`).Rows()
	if err != nil {
		return fmt.Errorf("query orders: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r orderRow
		if err := rows.Scan(&r.OrderID, &r.MerchantID, &r.CashierID, &r.TotalPrice, &r.CreatedAt); err != nil {
			return fmt.Errorf("scan order: %w", err)
		}
		event := events.OrderEvent{
			OrderID:    r.OrderID,
			CashierID:  r.CashierID,
			MerchantID: r.MerchantID,
			TotalPrice: r.TotalPrice,
			Status:     "created",
			EventTime:  r.CreatedAt.UTC().Format(time.RFC3339),
		}
		if err := b.repo.InsertOrderEvent(ctx, backfillEventID("order", r.OrderID), version, event); err != nil {
			return fmt.Errorf("insert order %d: %w", r.OrderID, err)
		}
		counts["order"]++
	}
	return rows.Err()
}

func (b *Backfiller) backfillOrderItems(ctx context.Context, version uint64, counts map[string]int) error {
	// order_items carry no category_id; join products to denormalize it so the
	// product/category sold statistics can be served purely from ClickHouse.
	rows, err := b.db.WithContext(ctx).Raw(`SELECT oi.order_item_id, oi.order_id, oi.product_id, p.category_id,
			oi.quantity, oi.price, oi.created_at
		FROM order_items oi
		JOIN products p ON p.product_id = oi.product_id
		WHERE oi.deleted_at IS NULL AND p.deleted_at IS NULL`).Rows()
	if err != nil {
		return fmt.Errorf("query order items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r itemRow
		if err := rows.Scan(&r.OrderItemID, &r.OrderID, &r.ProductID, &r.CategoryID, &r.Quantity, &r.UnitPrice, &r.CreatedAt); err != nil {
			return fmt.Errorf("scan order item: %w", err)
		}
		event := events.OrderItemEvent{
			OrderItemID: r.OrderItemID,
			OrderID:     r.OrderID,
			ProductID:   r.ProductID,
			CategoryID:  r.CategoryID,
			Quantity:    r.Quantity,
			UnitPrice:   r.UnitPrice,
			Subtotal:    r.Quantity * r.UnitPrice,
			EventTime:   r.CreatedAt.UTC().Format(time.RFC3339),
		}
		if err := b.repo.InsertOrderItemEvent(ctx, backfillEventID("order_item", r.OrderItemID), version, event); err != nil {
			return fmt.Errorf("insert order item %d: %w", r.OrderItemID, err)
		}
		counts["order_item"]++
	}
	return rows.Err()
}

func (b *Backfiller) backfillTransactions(ctx context.Context, version uint64, counts map[string]int) error {
	// transactions carry no cashier_id; join orders to denormalize it for the
	// cashier activity statistics.
	rows, err := b.db.WithContext(ctx).Raw(`SELECT t.transaction_id, t.order_id, t.merchant_id, o.cashier_id,
			t.payment_method, t.amount, t.payment_status, t.created_at
		FROM transactions t
		JOIN orders o ON o.order_id = t.order_id
		WHERE t.deleted_at IS NULL AND o.deleted_at IS NULL`).Rows()
	if err != nil {
		return fmt.Errorf("query transactions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r txRow
		if err := rows.Scan(&r.TransactionID, &r.OrderID, &r.MerchantID, &r.CashierID, &r.PaymentMethod, &r.Amount, &r.Status, &r.CreatedAt); err != nil {
			return fmt.Errorf("scan transaction: %w", err)
		}
		event := events.TransactionEvent{
			TransactionID: r.TransactionID,
			OrderID:       r.OrderID,
			CashierID:     r.CashierID,
			MerchantID:    r.MerchantID,
			PaymentMethod: r.PaymentMethod,
			Amount:        r.Amount,
			Status:        r.Status,
			EventTime:     r.CreatedAt.UTC().Format(time.RFC3339),
		}
		if err := b.repo.InsertTransactionEvent(ctx, backfillEventID("transaction", r.TransactionID), version, event); err != nil {
			return fmt.Errorf("insert transaction %d: %w", r.TransactionID, err)
		}
		counts["transaction"]++
	}
	return rows.Err()
}
