package merchantgraphqlmapper

import (
	"github.com/MamangRust/microservice-point-of-sale-apigateway/internal/model"
	merchantpb "github.com/MamangRust/microservice-point-of-sale-pb/merchant"
)

type MerchantGraphqlMapper interface {
	ToGraphqlResponseMerchant(res *merchantpb.ApiResponseMerchant) *model.APIResponseMerchant
	ToGraphqlResponsesMerchant(res *merchantpb.ApiResponsesMerchant) *model.APIResponsesMerchant
	ToGraphqlResponseMerchantDeleteAt(res *merchantpb.ApiResponseMerchantDeleteAt) *model.APIResponseMerchantDeleteAt
	ToGraphqlResponseMerchantDelete(res *merchantpb.ApiResponseMerchantDelete) *model.APIResponseMerchantDelete
	ToGraphqlResponseMerchantAll(res *merchantpb.ApiResponseMerchantAll) *model.APIResponseMerchantAll
	ToGraphqlResponsePaginationMerchant(res *merchantpb.ApiResponsePaginationMerchant) *model.APIResponsePaginationMerchant
	ToGraphqlResponsePaginationMerchantDeleteAt(res *merchantpb.ApiResponsePaginationMerchantDeleteAt) *model.APIResponsePaginationMerchantDeleteAt
	ToGraphqlResponseMerchantRestore(res *merchantpb.ApiResponseMerchant) *model.APIResponseMerchantDeleteAt
}
