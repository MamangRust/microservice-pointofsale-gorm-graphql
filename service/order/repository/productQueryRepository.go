package repository

import (
	"context"

	"github.com/MamangRust/microservice-point-of-sale-pkg/adapter"
	"github.com/MamangRust/microservice-point-of-sale-pkg/database/models"
	"github.com/MamangRust/microservice-point-of-sale-pkg/resilience"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/product_errors"
	pbproduct "github.com/MamangRust/microservice-pointofsale-grpc/pb/product"
)

type productQueryRepository struct {
	client pbproduct.ProductServiceClient
	guard  *resilience.DependencyGuard
}

func NewProductQueryRepository(client pbproduct.ProductServiceClient, opts ...adapter.GuardOption) ProductQueryRepository {
	r := &productQueryRepository{
		client: client,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *productQueryRepository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func parseNullableInt32(i int32) *int32 {
	return &i
}

func (r *productQueryRepository) FindById(ctx context.Context, product_id int) (*models.Product, error) {
	var resp *pbproduct.ApiResponseProduct
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.client.FindById(ctx, &pbproduct.FindByIdProductRequest{
			Id: int32(product_id),
		})
		return callErr
	})
	if err != nil {
		return nil, product_errors.ErrFindById
	}

	if resp == nil || resp.Data == nil {
		return nil, product_errors.ErrFindById
	}

	p := resp.Data
	return &models.Product{
		ProductID:    p.Id,
		MerchantID:   p.MerchantId,
		CategoryID:   p.CategoryId,
		Name:         p.Name,
		Description:  parseNullableString(p.Description),
		Price:        p.Price,
		CountInStock: p.CountInStock,
		Brand:        parseNullableString(p.Brand),
		Weight:       parseNullableInt32(p.Weight),
		SlugProduct:  parseNullableString(p.SlugProduct),
		ImageProduct: parseNullableString(p.ImageProduct),
		Barcode:      parseNullableString(p.Barcode),
		CreatedAt:    parseTimePtr(p.CreatedAt),
		UpdatedAt:    parseTimePtr(p.UpdatedAt),
	}, nil
}
