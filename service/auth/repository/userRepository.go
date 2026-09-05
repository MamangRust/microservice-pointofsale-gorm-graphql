package repository

import (
	"context"

	"github.com/MamangRust/microservice-point-of-sale-pkg/database/models"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/microservice-point-of-sale-shared/errors"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/user_errors"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindById(ctx context.Context, user_id int) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("user_id = ?", user_id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &user, nil
}

func (r *userRepository) FindByEmailAndVerify(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ? AND is_verified = ?", email, true).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &user, nil
}

func (r *userRepository) FindByVerificationCode(ctx context.Context, verification_code string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("verification_code = ?", verification_code).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, request *requests.RegisterRequest) (*models.User, error) {
	user := &models.User{
		Firstname:        request.FirstName,
		Lastname:         request.LastName,
		Email:            request.Email,
		Password:         request.Password,
		VerificationCode: request.VerifiedCode,
		IsVerified:       &request.IsVerified,
	}

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return user, nil
}

func (r *userRepository) UpdateUserIsVerified(ctx context.Context, user_id int, is_verified bool) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("user_id = ?", user_id).First(&user).Error
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	user.IsVerified = &is_verified
	if err := r.db.WithContext(ctx).Save(&user).Error; err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &user, nil
}

func (r *userRepository) UpdateUserPassword(ctx context.Context, user_id int, password string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("user_id = ?", user_id).First(&user).Error
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	user.Password = password
	if err := r.db.WithContext(ctx).Save(&user).Error; err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &user, nil
}
