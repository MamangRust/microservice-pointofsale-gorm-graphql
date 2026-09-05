package repository

import (
	"context"
	"fmt"

	"github.com/MamangRust/microservice-point-of-sale-pkg/database/models"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/order_errors"
	"gorm.io/gorm"
)

type productCommandRepository struct {
	db *gorm.DB
}

func NewProductCommandRepository(db *gorm.DB) ProductCommandRepository {
	return &productCommandRepository{db: db}
}

func (r *productCommandRepository) DecrementProductCountStock(ctx context.Context, productID int, quantity int) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).Model(&models.Product{}).Where("product_id = ?", productID).
		Update("count_in_stock", gorm.Expr("count_in_stock - ?", quantity)).Error
	if err != nil {
		return nil, err
	}

	// Check if stock went negative and revert
	if err := r.db.WithContext(ctx).First(&product, productID).Error; err != nil {
		return nil, err
	}
	if product.CountInStock < 0 {
		// Revert
		r.db.WithContext(ctx).Model(&models.Product{}).Where("product_id = ?", productID).Update("count_in_stock", gorm.Expr("count_in_stock + ?", quantity))
		return nil, order_errors.ErrInsufficientProductStock
	}

	return &product, nil
}

func (r *productCommandRepository) IncrementProductCountStock(ctx context.Context, productID int, quantity int) (*models.Product, error) {
	var product models.Product
	if err := r.db.WithContext(ctx).Model(&models.Product{}).Where("product_id = ?", productID).
		Update("count_in_stock", gorm.Expr("count_in_stock + ?", quantity)).Error; err != nil {
		return nil, fmt.Errorf("failed to increment stock: %w", err)
	}

	if err := r.db.WithContext(ctx).First(&product, productID).Error; err != nil {
		return nil, err
	}

	return &product, nil
}
