package repository

import (
	"context"

	"github.com/MamangRust/microservice-point-of-sale-pkg/database/models"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/merchant_errors"
	pbmerchant "github.com/MamangRust/microservice-pointofsale-grpc/pb/merchant"
)

type merchantQueryRepository struct {
	client pbmerchant.MerchantServiceClient
}

func NewMerchantQueryRepository(client pbmerchant.MerchantServiceClient) MerchantQueryRepository {
	return &merchantQueryRepository{client: client}
}

func (r *merchantQueryRepository) FindById(ctx context.Context, id int) (*models.Merchant, error) {
	res, err := r.client.FindById(ctx, &pbmerchant.FindByIdMerchantRequest{Id: int32(id)})
	if err != nil || res == nil || res.Data == nil {
		return nil, merchant_errors.ErrFindById
	}

	merchant := &models.Merchant{
		MerchantID:   res.Data.Id,
		UserID:       res.Data.UserId,
		Name:         res.Data.Name,
		Status:       res.Data.Status,
	}
	if res.Data.Description != "" {
		merchant.Description = &res.Data.Description
	}
	if res.Data.Address != "" {
		merchant.Address = &res.Data.Address
	}
	if res.Data.ContactEmail != "" {
		merchant.ContactEmail = &res.Data.ContactEmail
	}
	if res.Data.ContactPhone != "" {
		merchant.ContactPhone = &res.Data.ContactPhone
	}
	return merchant, nil
}
