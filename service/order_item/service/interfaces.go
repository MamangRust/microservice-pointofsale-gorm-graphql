package service

import (
	"context"

	"github.com/MamangRust/microservice-point-of-sale-pkg/database/models"
	"github.com/MamangRust/microservice-point-of-sale-order-item/repository"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
)

type OrderItemQueryService interface {
	FindAllOrderItems(ctx context.Context, req *requests.FindAllOrderItems) ([]*repository.OrderItemResult, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*repository.OrderItemResultDeleteAt, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*repository.OrderItemResultDeleteAt, *int, error)
	FindOrderItemByOrder(ctx context.Context, orderID int) ([]*models.OrderItem, error)
}
