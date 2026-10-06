package cashiergraphqlmapper

import (
	"github.com/MamangRust/microservice-point-of-sale-apigateway/internal/model"
	cashierpb "github.com/MamangRust/microservice-point-of-sale-pb/cashier"
)

type CashierGraphqlMapper interface {
	ToGraphqlResponseCashier(res *cashierpb.ApiResponseCashier) *model.APIResponseCashier
	ToGraphqlResponsesCashier(res *cashierpb.ApiResponsesCashier) *model.APIResponsesCashier
	ToGraphqlResponseCashierDeleteAt(res *cashierpb.ApiResponseCashierDeleteAt) *model.APIResponseCashierDeleteAt
	ToGraphqlResponseCashierDelete(res *cashierpb.ApiResponseCashierDelete) *model.APIResponseCashierDelete
	ToGraphqlResponseCashierAll(res *cashierpb.ApiResponseCashierAll) *model.APIResponseCashierAll
	ToGraphqlResponsePaginationCashier(res *cashierpb.ApiResponsePaginationCashier) *model.APIResponsePaginationCashier
	ToGraphqlResponsePaginationCashierDeleteAt(res *cashierpb.ApiResponsePaginationCashierDeleteAt) *model.APIResponsePaginationCashierDeleteAt
	ToGraphqlResponseMonthlyTotalSales(res *cashierpb.ApiResponseCashierMonthlyTotalSales) *model.APIResponseCashierMonthlyTotalSales
	ToGraphqlResponseMonthlySales(res *cashierpb.ApiResponseCashierMonthSales) *model.APIResponseCashierMonthSales
	ToGraphqlResponseYearlySales(res *cashierpb.ApiResponseCashierYearSales) *model.APIResponseCashierYearSales
	ToGraphqlResponseYearlyTotalSales(res *cashierpb.ApiResponseCashierYearlyTotalSales) *model.APIResponseCashierYearlyTotalSales
}
