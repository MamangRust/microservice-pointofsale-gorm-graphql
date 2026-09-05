package service

import (
	"context"

	"github.com/MamangRust/microservice-point-of-sale-pkg/database/models"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/microservice-point-of-sale-user/repository"
)

type UserQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllUsers) ([]*repository.UserResult, *int, error)
	FindByID(ctx context.Context, id int) (*models.User, error)
	FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*repository.UserResult, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*repository.UserResult, *int, error)
}

type UserCommandService interface {
	CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*models.User, error)
	UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*models.User, error)
	TrashedUser(ctx context.Context, user_id int) (*models.User, error)
	RestoreUser(ctx context.Context, user_id int) (*models.User, error)
	DeleteUserPermanent(ctx context.Context, user_id int) (bool, error)
	RestoreAllUser(ctx context.Context) (bool, error)
	DeleteAllUserPermanent(ctx context.Context) (bool, error)
}
