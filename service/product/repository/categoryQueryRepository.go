package repository

import (
	"context"

	"github.com/MamangRust/microservice-point-of-sale-pkg/database/models"
	"gorm.io/gorm"
)

type categoryQueryRepository struct {
	db *gorm.DB
}

func NewCategoryQueryRepository(db *gorm.DB) CategoryQueryRepository {
	return &categoryQueryRepository{db: db}
}

func (r *categoryQueryRepository) FindById(ctx context.Context, categoryID int) (*models.Category, error) {
	var item models.Category
	if err := r.db.WithContext(ctx).Where("category_id = ? AND deleted_at IS NULL", categoryID).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
