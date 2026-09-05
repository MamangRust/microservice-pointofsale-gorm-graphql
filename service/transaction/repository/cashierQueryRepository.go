package repository

import (
	"context"
	
	"github.com/MamangRust/microservice-point-of-sale-pkg/adapter"
	"github.com/MamangRust/microservice-point-of-sale-pkg/database/models"
	"github.com/MamangRust/microservice-point-of-sale-pkg/resilience"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/cashier_errors"
	pbcashier "github.com/MamangRust/microservice-pointofsale-grpc/pb/cashier"
)

type cashierQueryRepository struct {
	client pbcashier.CashierServiceClient
	guard  *resilience.DependencyGuard
}

func NewCashierQueryRepository(client pbcashier.CashierServiceClient, opts ...adapter.GuardOption) CashierQueryRepository {
	r := &cashierQueryRepository{client: client}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *cashierQueryRepository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func (r *cashierQueryRepository) FindById(ctx context.Context, cashier_id int) (*models.Cashier, error) {
	var resp *pbcashier.ApiResponseCashier
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.client.FindById(ctx, &pbcashier.FindByIdCashierRequest{Id: int32(cashier_id)})
		return callErr
	})
	if err != nil {
		return nil, cashier_errors.ErrFindCashierById
	}
	if resp == nil || resp.Data == nil {
		return nil, cashier_errors.ErrFindCashierById
	}
	c := resp.Data
	createdAt := parseTxTime(c.CreatedAt)
	updatedAt := parseTxTime(c.UpdatedAt)
	return &models.Cashier{
		CashierID:  c.Id,
		MerchantID: c.MerchantId,
		Name:       c.Name,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}, nil
}
