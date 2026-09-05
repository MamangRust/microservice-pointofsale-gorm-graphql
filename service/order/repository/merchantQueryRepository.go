package repository

import (
	"context"

	"github.com/MamangRust/microservice-point-of-sale-pkg/adapter"
	"github.com/MamangRust/microservice-point-of-sale-pkg/database/models"
	"github.com/MamangRust/microservice-point-of-sale-pkg/resilience"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/merchant_errors"
	pbmerchant "github.com/MamangRust/microservice-pointofsale-grpc/pb/merchant"
)

type merchantQueryRepository struct {
	client pbmerchant.MerchantServiceClient
	guard  *resilience.DependencyGuard
}

func NewMerchantQueryRepository(client pbmerchant.MerchantServiceClient, opts ...adapter.GuardOption) MerchantQueryRepository {
	r := &merchantQueryRepository{
		client: client,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *merchantQueryRepository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func parseNullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func (r *merchantQueryRepository) FindById(ctx context.Context, merchantID int) (*models.Merchant, error) {
	var resp *pbmerchant.ApiResponseMerchant
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.client.FindById(ctx, &pbmerchant.FindByIdMerchantRequest{
			Id: int32(merchantID),
		})
		return callErr
	})
	if err != nil {
		return nil, merchant_errors.ErrFindById
	}

	if resp == nil || resp.Data == nil {
		return nil, merchant_errors.ErrFindById
	}

	m := resp.Data
	return &models.Merchant{
		MerchantID:   m.Id,
		UserID:       m.UserId,
		Name:         m.Name,
		Description:  parseNullableString(m.Description),
		Address:      parseNullableString(m.Address),
		ContactEmail: parseNullableString(m.ContactEmail),
		ContactPhone: parseNullableString(m.ContactPhone),
		Status:       m.Status,
		CreatedAt:    parseTimePtr(m.CreatedAt),
		UpdatedAt:    parseTimePtr(m.UpdatedAt),
	}, nil
}
