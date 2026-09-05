package repository

import (
	"context"

	"github.com/MamangRust/microservice-point-of-sale-pkg/database/models"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
)

type UserResult struct {
	UserID    int32
	Firstname string
	Lastname  string
	Email     string
	CreatedAt *string
	UpdatedAt *string
	DeletedAt *string
	TotalCount int64
}

type UserQueryRepository interface {
	FindAllUsers(ctx context.Context, req *requests.FindAllUsers) ([]*UserResult, error)
	FindById(ctx context.Context, user_id int) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*UserResult, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*UserResult, error)
}

type UserCommandRepository interface {
	CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*models.User, error)
	UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*models.User, error)
	TrashedUser(ctx context.Context, user_id int) (*models.User, error)
	RestoreUser(ctx context.Context, user_id int) (*models.User, error)
	DeleteUserPermanent(ctx context.Context, user_id int) (bool, error)
	RestoreAllUser(ctx context.Context) (bool, error)
	DeleteAllUserPermanent(ctx context.Context) (bool, error)
}

type RoleQueryRepository interface {
	FindByName(ctx context.Context, name string) (*models.Role, error)
}
