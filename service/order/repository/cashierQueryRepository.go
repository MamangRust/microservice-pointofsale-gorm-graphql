package repository

import (
	"context"
	"time"

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
	r := &cashierQueryRepository{
		client: client,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *cashierQueryRepository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func (r *cashierQueryRepository) FindById(ctx context.Context, cashierID int) (*models.Cashier, error) {
	var resp *pbcashier.ApiResponseCashier
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.client.FindById(ctx, &pbcashier.FindByIdCashierRequest{
			Id: int32(cashierID),
		})
		return callErr
	})
	if err != nil {
		return nil, cashier_errors.ErrFindCashierById
	}

	if resp == nil || resp.Data == nil {
		return nil, cashier_errors.ErrFindCashierById
	}

	c := resp.Data
	return &models.Cashier{
		CashierID:  c.Id,
		MerchantID: c.MerchantId,
		Name:       c.Name,
		CreatedAt:  parseTimePtr(c.CreatedAt),
		UpdatedAt:  parseTimePtr(c.UpdatedAt),
	}, nil
}

func parseTimePtr(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		t, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return nil
		}
	}
	return &t
}
