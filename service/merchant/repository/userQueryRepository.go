package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-point-of-sale-pkg/database/models"
	sharedErrors "github.com/MamangRust/microservice-point-of-sale-shared/errors"
	pbuser "github.com/MamangRust/microservice-pointofsale-grpc/pb/user"
)

type userQueryRepository struct {
	client pbuser.UserServiceClient
}

func NewUserQueryRepository(client pbuser.UserServiceClient) UserQueryRepository {
	return &userQueryRepository{
		client: client,
	}
}

func (r *userQueryRepository) FindById(ctx context.Context, userID int) (*models.User, error) {
	res, err := r.client.FindById(ctx, &pbuser.FindByIdUserRequest{
		Id: int32(userID),
	})
	if err != nil || res == nil || res.Data == nil {
		return nil, sharedErrors.ErrInternal
	}

	user := &models.User{
		UserID:    res.Data.Id,
		Firstname: res.Data.Firstname,
		Lastname:  res.Data.Lastname,
		Email:     res.Data.Email,
	}

	if res.Data.CreatedAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", res.Data.CreatedAt); err == nil {
			user.CreatedAt = &t
		} else if t, err = time.Parse(time.RFC3339, res.Data.CreatedAt); err == nil {
			user.CreatedAt = &t
		}
	}
	if res.Data.UpdatedAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", res.Data.UpdatedAt); err == nil {
			user.UpdatedAt = &t
		} else if t, err = time.Parse(time.RFC3339, res.Data.UpdatedAt); err == nil {
			user.UpdatedAt = &t
		}
	}

	return user, nil
}
