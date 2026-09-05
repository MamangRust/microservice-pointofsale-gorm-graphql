package repository

import (
	"context"

	"github.com/MamangRust/microservice-point-of-sale-pkg/adapter"
	"github.com/MamangRust/microservice-point-of-sale-pkg/database/models"
	"github.com/MamangRust/microservice-point-of-sale-pkg/resilience"
	orderitem_errors "github.com/MamangRust/microservice-point-of-sale-shared/errors/order_item_errors"
	pborderitem "github.com/MamangRust/microservice-pointofsale-grpc/pb/order_item"
)

type orderItemQueryRepository struct {
	client pborderitem.OrderItemServiceClient
	guard  *resilience.DependencyGuard
}

func NewOrderItemQueryRepository(client pborderitem.OrderItemServiceClient, opts ...adapter.GuardOption) OrderItemQueryRepository {
	r := &orderItemQueryRepository{
		client: client,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *orderItemQueryRepository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func (r *orderItemQueryRepository) CalculateTotalPrice(ctx context.Context, orderID int) (*int32, error) {
	items, err := r.FindOrderItemByOrder(ctx, orderID)
	if err != nil {
		return nil, orderitem_errors.ErrCalculateTotalPrice
	}

	var total int32 = 0
	for _, item := range items {
		if item != nil {
			total += item.Quantity * int32(item.Price)
		}
	}

	return &total, nil
}

func (r *orderItemQueryRepository) FindOrderItemByOrder(ctx context.Context, orderID int) ([]*models.OrderItem, error) {
	var resp *pborderitem.ApiResponsesOrderItem
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.client.FindOrderItemByOrder(ctx, &pborderitem.FindByIdOrderItemRequest{
			Id: int32(orderID),
		})
		return callErr
	})
	if err != nil {
		return nil, orderitem_errors.ErrFindOrderItemByOrder
	}

	if resp == nil || resp.Data == nil {
		return nil, orderitem_errors.ErrFindOrderItemByOrder
	}

	var res []*models.OrderItem
	for _, item := range resp.Data {
		if item == nil {
			continue
		}
		res = append(res, &models.OrderItem{
			OrderItemID: item.Id,
			OrderID:     item.OrderId,
			ProductID:   item.ProductId,
			Quantity:    item.Quantity,
			Price:       int64(item.Price),
			CreatedAt:   parseTimePtr(item.CreatedAt),
			UpdatedAt:   parseTimePtr(item.UpdatedAt),
		})
	}

	return res, nil
}
