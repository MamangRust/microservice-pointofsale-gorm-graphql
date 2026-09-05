package repository

import (
	"context"

	"github.com/MamangRust/microservice-point-of-sale-pkg/database/models"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	"gorm.io/gorm"
)

type userQueryRepository struct {
	db *gorm.DB
}

func NewUserQueryRepository(db *gorm.DB) *userQueryRepository {
	return &userQueryRepository{db: db}
}

func (r *userQueryRepository) FindAllUsers(ctx context.Context, req *requests.FindAllUsers) ([]*UserResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*UserResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT u.user_id, u.firstname, u.lastname, u.email,
			COALESCE(TO_CHAR(u.created_at, 'YYYY-MM-DD HH24:MI:SS'), '') as created_at,
			COALESCE(TO_CHAR(u.updated_at, 'YYYY-MM-DD HH24:MI:SS'), '') as updated_at,
			COALESCE(TO_CHAR(u.deleted_at, 'YYYY-MM-DD HH24:MI:SS'), '') as deleted_at,
			COUNT(*) OVER() as total_count
		FROM users u
		WHERE u.deleted_at IS NULL
			AND (? = '' OR u.firstname ILIKE ? OR u.lastname ILIKE ? OR u.email ILIKE ?)
		ORDER BY u.user_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%", req.PageSize, offset).Scan(&results).Error
	return results, err
}

func (r *userQueryRepository) FindById(ctx context.Context, userID int) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("user_id = ? AND deleted_at IS NULL", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userQueryRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ? AND deleted_at IS NULL", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*UserResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*UserResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT u.user_id, u.firstname, u.lastname, u.email,
			COALESCE(TO_CHAR(u.created_at, 'YYYY-MM-DD HH24:MI:SS'), '') as created_at,
			COALESCE(TO_CHAR(u.updated_at, 'YYYY-MM-DD HH24:MI:SS'), '') as updated_at,
			COALESCE(TO_CHAR(u.deleted_at, 'YYYY-MM-DD HH24:MI:SS'), '') as deleted_at,
			COUNT(*) OVER() as total_count
		FROM users u
		WHERE u.deleted_at IS NULL
			AND (? = '' OR u.firstname ILIKE ? OR u.lastname ILIKE ? OR u.email ILIKE ?)
		ORDER BY u.user_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%", req.PageSize, offset).Scan(&results).Error
	return results, err
}

func (r *userQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*UserResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*UserResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT u.user_id, u.firstname, u.lastname, u.email,
			COALESCE(TO_CHAR(u.created_at, 'YYYY-MM-DD HH24:MI:SS'), '') as created_at,
			COALESCE(TO_CHAR(u.updated_at, 'YYYY-MM-DD HH24:MI:SS'), '') as updated_at,
			COALESCE(TO_CHAR(u.deleted_at, 'YYYY-MM-DD HH24:MI:SS'), '') as deleted_at,
			COUNT(*) OVER() as total_count
		FROM users u
		WHERE u.deleted_at IS NOT NULL
			AND (? = '' OR u.firstname ILIKE ? OR u.lastname ILIKE ? OR u.email ILIKE ?)
		ORDER BY u.user_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%", req.PageSize, offset).Scan(&results).Error
	return results, err
}
