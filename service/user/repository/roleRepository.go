package repository

import (
	"context"

	"github.com/MamangRust/microservice-point-of-sale-pkg/database/models"
	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *roleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) FindByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	err := r.db.WithContext(ctx).Where("role_name = ? AND deleted_at IS NULL", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}
