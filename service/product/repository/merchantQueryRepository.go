package repository

import (
	"context"

	"github.com/MamangRust/microservice-point-of-sale-pkg/database/models"
	"gorm.io/gorm"
)

type merchantQueryRepository struct {
	db *gorm.DB
}

func NewMerchantQueryRepository(db *gorm.DB) MerchantQueryRepository {
	return &merchantQueryRepository{db: db}
}

func (r *merchantQueryRepository) FindById(ctx context.Context, merchantID int) (*models.Merchant, error) {
	var item models.Merchant
	if err := r.db.WithContext(ctx).Where("merchant_id = ? AND deleted_at IS NULL", merchantID).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
