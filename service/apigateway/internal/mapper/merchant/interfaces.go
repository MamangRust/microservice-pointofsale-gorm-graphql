package merchantgraphqlmapper

import (
	"github.com/MamangRust/microservice-point-of-sale-apigateway/internal/model"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/merchant"
)

type MerchantGraphqlMapper interface {
	ToGraphqlResponseMerchant(res *pb.ApiResponseMerchant) *model.APIResponseMerchant
	ToGraphqlResponsesMerchant(res *pb.ApiResponsesMerchant) *model.APIResponsesMerchant
	ToGraphqlResponseMerchantDeleteAt(res *pb.ApiResponseMerchantDeleteAt) *model.APIResponseMerchantDeleteAt
	ToGraphqlResponseMerchantDelete(res *pb.ApiResponseMerchantDelete) *model.APIResponseMerchantDelete
	ToGraphqlResponseMerchantAll(res *pb.ApiResponseMerchantAll) *model.APIResponseMerchantAll
	ToGraphqlResponsePaginationMerchant(res *pb.ApiResponsePaginationMerchant) *model.APIResponsePaginationMerchant
	ToGraphqlResponsePaginationMerchantDeleteAt(res *pb.ApiResponsePaginationMerchantDeleteAt) *model.APIResponsePaginationMerchantDeleteAt
	ToGraphqlResponseMerchantRestore(res *pb.ApiResponseMerchant) *model.APIResponseMerchantDeleteAt
}
