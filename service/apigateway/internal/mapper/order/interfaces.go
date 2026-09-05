package ordergraphqlmapper

import (
	"github.com/MamangRust/microservice-point-of-sale-apigateway/internal/model"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/order"
	pbstats "github.com/MamangRust/microservice-pointofsale-grpc/pb/stats"
)

type OrderGraphqlMapper interface {
	ToGraphqlResponseOrder(res *pb.ApiResponseOrder) *model.APIResponseOrder
	ToGraphqlResponsesOrder(res *pb.ApiResponsesOrder) *model.APIResponsesOrder
	ToGraphqlResponseOrderDeleteAt(res *pb.ApiResponseOrderDeleteAt) *model.APIResponseOrderDeleteAt
	ToGraphqlResponseOrderDelete(res *pb.ApiResponseOrderDelete) *model.APIResponseOrderDelete
	ToGraphqlResponseOrderAll(res *pb.ApiResponseOrderAll) *model.APIResponseOrderAll
	ToGraphqlResponsePaginationOrder(res *pb.ApiResponsePaginationOrder) *model.APIResponsePaginationOrder
	ToGraphqlResponsePaginationOrderDeleteAt(res *pb.ApiResponsePaginationOrderDeleteAt) *model.APIResponsePaginationOrderDeleteAt
	ToGraphqlResponseMonthlyRevenue(res *pbstats.ApiResponseOrderMonthly) *model.APIResponseOrderMonthly
	ToGraphqlResponseYearlyRevenue(res *pbstats.ApiResponseOrderYearly) *model.APIResponseOrderYearly
	ToGraphqlResponseMonthlyTotalRevenue(res *pbstats.ApiResponseOrderMonthlyTotalRevenue) *model.APIResponseOrderMonthlyTotalRevenue
	ToGraphqlResponseYearlyTotalRevenue(res *pbstats.ApiResponseOrderYearlyTotalRevenue) *model.APIResponseOrderYearlyTotalRevenue
}
