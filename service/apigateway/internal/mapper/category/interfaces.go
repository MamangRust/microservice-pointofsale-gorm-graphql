package categorygraphqlmapper

import (
	"github.com/MamangRust/microservice-point-of-sale-apigateway/internal/model"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/category"
	pbstats "github.com/MamangRust/microservice-pointofsale-grpc/pb/stats"
)

type CategoryGraphqlMapper interface {
	ToGraphqlResponseCategory(res *pb.ApiResponseCategory) *model.APIResponseCategory
	ToGraphqlResponsesCategory(res *pb.ApiResponsesCategory) *model.APIResponsesCategory
	ToGraphqlResponseCategoryDeleteAt(res *pb.ApiResponseCategoryDeleteAt) *model.APIResponseCategoryDeleteAt
	ToGraphqlResponseCategoryDelete(res *pb.ApiResponseCategoryDelete) *model.APIResponseCategoryDelete
	ToGraphqlResponseCategoryAll(res *pb.ApiResponseCategoryAll) *model.APIResponseCategoryAll
	ToGraphqlResponsePaginationCategory(res *pb.ApiResponsePaginationCategory) *model.APIResponsePaginationCategory
	ToGraphqlResponseCategoryMonthlyTotalPrice(res *pbstats.ApiResponseCategoryMonthlyTotalPrice) *model.APIResponseCategoryMonthlyTotalPrice
	ToGraphqlResponseCategoryYearlyTotalPrice(res *pbstats.ApiResponseCategoryYearlyTotalPrice) *model.APIResponseCategoryYearlyTotalPrice
	ToGraphqlResponseCategoryMonthlyPrice(res *pbstats.ApiResponseCategoryMonthPrice) *model.APIResponseCategoryMonthPrice
	ToGraphqlResponseCategoryYearlyPrice(res *pbstats.ApiResponseCategoryYearPrice) *model.APIResponseCategoryYearPrice
	ToGraphqlResponsePaginationCategoryDeleteAt(res *pb.ApiResponsePaginationCategoryDeleteAt) *model.APIResponsePaginationCategoryDeleteAt
}
