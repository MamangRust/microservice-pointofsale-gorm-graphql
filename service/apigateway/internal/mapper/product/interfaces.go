package productgraphqlmapper

import (
	"github.com/MamangRust/microservice-point-of-sale-apigateway/internal/model"
	productpb "github.com/MamangRust/microservice-point-of-sale-pb/product"
)

type ProductGraphqlMapper interface {
	ToGraphqlResponseProduct(res *productpb.ApiResponseProduct) *model.APIResponseProduct
	ToGraphqlResponsesProduct(res *productpb.ApiResponsesProduct) *model.APIResponsesProduct
	ToGraphqlResponseProductDeleteAt(res *productpb.ApiResponseProductDeleteAt) *model.APIResponseProductDeleteAt
	ToGraphqlResponseProductDelete(res *productpb.ApiResponseProductDelete) *model.APIResponseProductDelete
	ToGraphqlResponseProductAll(res *productpb.ApiResponseProductAll) *model.APIResponseProductAll
	ToGraphqlResponsePaginationProduct(res *productpb.ApiResponsePaginationProduct) *model.APIResponsePaginationProduct
	ToGraphqlResponsePaginationProductDeleteAt(res *productpb.ApiResponsePaginationProductDeleteAt) *model.APIResponsePaginationProductDeleteAt
}
