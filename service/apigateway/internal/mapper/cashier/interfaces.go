package cashiergraphqlmapper

import (
	"github.com/MamangRust/microservice-point-of-sale-apigateway/internal/model"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/cashier"
	pbstats "github.com/MamangRust/microservice-pointofsale-grpc/pb/stats"
)

type CashierGraphqlMapper interface {
	ToGraphqlResponseCashier(res *pb.ApiResponseCashier) *model.APIResponseCashier
	ToGraphqlResponsesCashier(res *pb.ApiResponsesCashier) *model.APIResponsesCashier
	ToGraphqlResponseCashierDeleteAt(res *pb.ApiResponseCashierDeleteAt) *model.APIResponseCashierDeleteAt
	ToGraphqlResponseCashierDelete(res *pb.ApiResponseCashierDelete) *model.APIResponseCashierDelete
	ToGraphqlResponseCashierAll(res *pb.ApiResponseCashierAll) *model.APIResponseCashierAll
	ToGraphqlResponsePaginationCashier(res *pb.ApiResponsePaginationCashier) *model.APIResponsePaginationCashier
	ToGraphqlResponsePaginationCashierDeleteAt(res *pb.ApiResponsePaginationCashierDeleteAt) *model.APIResponsePaginationCashierDeleteAt
	ToGraphqlResponseMonthlyTotalSales(res *pbstats.ApiResponseCashierMonthlyTotalSales) *model.APIResponseCashierMonthlyTotalSales
	ToGraphqlResponseMonthlySales(res *pbstats.ApiResponseCashierMonthSales) *model.APIResponseCashierMonthSales
	ToGraphqlResponseYearlySales(res *pbstats.ApiResponseCashierYearSales) *model.APIResponseCashierYearSales
	ToGraphqlResponseYearlyTotalSales(res *pbstats.ApiResponseCashierYearlyTotalSales) *model.APIResponseCashierYearlyTotalSales
}
