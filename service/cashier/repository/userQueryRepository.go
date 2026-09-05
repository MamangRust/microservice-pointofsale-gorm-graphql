package repository

import (
	"context"

	"github.com/MamangRust/microservice-point-of-sale-pkg/database/models"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/user_errors"
	pbuser "github.com/MamangRust/microservice-pointofsale-grpc/pb/user"
)

type userQueryRepository struct {
	client pbuser.UserServiceClient
}

func NewUserQueryRepository(client pbuser.UserServiceClient) UserQueryRepository {
	return &userQueryRepository{client: client}
}

func (r *userQueryRepository) FindById(ctx context.Context, id int) (*models.User, error) {
	res, err := r.client.FindById(ctx, &pbuser.FindByIdUserRequest{Id: int32(id)})
	if err != nil || res == nil || res.Data == nil {
		return nil, user_errors.ErrUserNotFound
	}

	user := &models.User{
		UserID:    res.Data.Id,
		Firstname: res.Data.Firstname,
		Lastname:  res.Data.Lastname,
		Email:     res.Data.Email,
	}
	return user, nil
}
