package repository

import (
	"context"

	"github.com/MamangRust/microservice-point-of-sale-pkg/adapter"
	"github.com/MamangRust/microservice-point-of-sale-pkg/database/models"
	"github.com/MamangRust/microservice-point-of-sale-pkg/resilience"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/order_errors"
	pborder "github.com/MamangRust/microservice-pointofsale-grpc/pb/order"
)

type orderQueryRepository struct {
	client pborder.OrderServiceClient
	guard  *resilience.DependencyGuard
}

func NewOrderQueryRepository(client pborder.OrderServiceClient, opts ...adapter.GuardOption) OrderQueryRepository {
	r := &orderQueryRepository{client: client}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *orderQueryRepository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func (r *orderQueryRepository) FindById(ctx context.Context, order_id int) (*models.Order, error) {
	var resp *pborder.ApiResponseOrder
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.client.FindById(ctx, &pborder.FindByIdOrderRequest{Id: int32(order_id)})
		return callErr
	})
	if err != nil {
		return nil, order_errors.ErrFindById
	}
	if resp == nil || resp.Data == nil {
		return nil, order_errors.ErrFindById
	}
	o := resp.Data
	createdAt := parseTxTime(o.CreatedAt)
	updatedAt := parseTxTime(o.UpdatedAt)
	return &models.Order{
		OrderID:    o.Id,
		MerchantID: o.MerchantId,
		CashierID:  o.CashierId,
		TotalPrice: int64(o.TotalPrice),
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}, nil
}
