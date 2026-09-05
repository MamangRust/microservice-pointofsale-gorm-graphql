package response_api

import (
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/response"
	pbstats "github.com/MamangRust/microservice-pointofsale-grpc/pb/stats"
)

// StatsResponseMapper maps the stats-reader (ClickHouse) gRPC responses to the
// REST shapes exposed by the API gateway (F4 §7.4).
type StatsResponseMapper interface {
	ToApiResponseOrderMonthlyRevenue(res *pbstats.ApiResponseOrderMonthlyRevenue) *response.ApiResponseOrderMonthlyRevenue
	ToApiResponseOrderYearlyRevenue(res *pbstats.ApiResponseOrderYearlyRevenue) *response.ApiResponseOrderYearlyRevenue
	ToApiResponseCashierMonthlyRevenue(res *pbstats.ApiResponseCashierMonthlyRevenue) *response.ApiResponseCashierMonthlyRevenue
	ToApiResponseProductMonthlySold(res *pbstats.ApiResponseProductMonthlySold) *response.ApiResponseProductMonthlySold
	ToApiResponseCategoryMonthlySold(res *pbstats.ApiResponseCategoryMonthlySold) *response.ApiResponseCategoryMonthlySold
	ToApiResponseTransactionMonthlySuccess(res *pbstats.ApiResponseTransactionMonthlySuccess) *response.ApiResponseTransactionMonthlySuccess
	ToApiResponseCashierMonthlyOrders(res *pbstats.ApiResponseCashierMonthlyOrders) *response.ApiResponseCashierMonthlyOrders
}

type statsResponseMapper struct{}

func NewStatsResponseMapper() *statsResponseMapper {
	return &statsResponseMapper{}
}

func (s *statsResponseMapper) ToApiResponseOrderMonthlyRevenue(res *pbstats.ApiResponseOrderMonthlyRevenue) *response.ApiResponseOrderMonthlyRevenue {
	data := make([]*response.OrderMonthlyRevenueResponse, 0, len(res.GetData()))
	for _, row := range res.GetData() {
		data = append(data, &response.OrderMonthlyRevenueResponse{
			Year:         row.GetYear(),
			Month:        row.GetMonth(),
			TotalRevenue: row.GetTotalRevenue(),
			OrderCount:   int(row.GetOrderCount()),
		})
	}
	return &response.ApiResponseOrderMonthlyRevenue{Status: res.GetStatus(), Message: res.GetMessage(), Data: data}
}

func (s *statsResponseMapper) ToApiResponseOrderYearlyRevenue(res *pbstats.ApiResponseOrderYearlyRevenue) *response.ApiResponseOrderYearlyRevenue {
	data := make([]*response.OrderYearlyRevenueResponse, 0, len(res.GetData()))
	for _, row := range res.GetData() {
		data = append(data, &response.OrderYearlyRevenueResponse{
			Year:         row.GetYear(),
			TotalRevenue: row.GetTotalRevenue(),
			OrderCount:   int(row.GetOrderCount()),
		})
	}
	return &response.ApiResponseOrderYearlyRevenue{Status: res.GetStatus(), Message: res.GetMessage(), Data: data}
}

func (s *statsResponseMapper) ToApiResponseCashierMonthlyRevenue(res *pbstats.ApiResponseCashierMonthlyRevenue) *response.ApiResponseCashierMonthlyRevenue {
	data := make([]*response.CashierMonthlyRevenueResponse, 0, len(res.GetData()))
	for _, row := range res.GetData() {
		data = append(data, &response.CashierMonthlyRevenueResponse{
			Year:         row.GetYear(),
			Month:        row.GetMonth(),
			CashierID:    int(row.GetCashierId()),
			TotalRevenue: row.GetTotalRevenue(),
			OrderCount:   int(row.GetOrderCount()),
		})
	}
	return &response.ApiResponseCashierMonthlyRevenue{Status: res.GetStatus(), Message: res.GetMessage(), Data: data}
}

func (s *statsResponseMapper) ToApiResponseProductMonthlySold(res *pbstats.ApiResponseProductMonthlySold) *response.ApiResponseProductMonthlySold {
	data := make([]*response.ProductMonthlySoldResponse, 0, len(res.GetData()))
	for _, row := range res.GetData() {
		data = append(data, &response.ProductMonthlySoldResponse{
			Month:     row.GetMonth(),
			ProductID: int(row.GetProductId()),
			Quantity:  row.GetQuantity(),
			Subtotal:  row.GetSubtotal(),
		})
	}
	return &response.ApiResponseProductMonthlySold{Status: res.GetStatus(), Message: res.GetMessage(), Data: data}
}

func (s *statsResponseMapper) ToApiResponseCategoryMonthlySold(res *pbstats.ApiResponseCategoryMonthlySold) *response.ApiResponseCategoryMonthlySold {
	data := make([]*response.CategoryMonthlySoldResponse, 0, len(res.GetData()))
	for _, row := range res.GetData() {
		data = append(data, &response.CategoryMonthlySoldResponse{
			Month:      row.GetMonth(),
			CategoryID: int(row.GetCategoryId()),
			Quantity:   row.GetQuantity(),
			Subtotal:   row.GetSubtotal(),
		})
	}
	return &response.ApiResponseCategoryMonthlySold{Status: res.GetStatus(), Message: res.GetMessage(), Data: data}
}

func (s *statsResponseMapper) ToApiResponseTransactionMonthlySuccess(res *pbstats.ApiResponseTransactionMonthlySuccess) *response.ApiResponseTransactionMonthlySuccess {
	data := make([]*response.TransactionMonthlySuccessResponse, 0, len(res.GetData()))
	for _, row := range res.GetData() {
		data = append(data, &response.TransactionMonthlySuccessResponse{
			Month:       row.GetMonth(),
			TotalCount:  row.GetTotalCount(),
			TotalAmount: row.GetTotalAmount(),
		})
	}
	return &response.ApiResponseTransactionMonthlySuccess{Status: res.GetStatus(), Message: res.GetMessage(), Data: data}
}

func (s *statsResponseMapper) ToApiResponseCashierMonthlyOrders(res *pbstats.ApiResponseCashierMonthlyOrders) *response.ApiResponseCashierMonthlyOrders {
	data := make([]*response.CashierMonthlyOrdersResponse, 0, len(res.GetData()))
	for _, row := range res.GetData() {
		data = append(data, &response.CashierMonthlyOrdersResponse{
			Month:       row.GetMonth(),
			CashierID:   int(row.GetCashierId()),
			OrderCount:  row.GetOrderCount(),
			TotalAmount: row.GetTotalAmount(),
		})
	}
	return &response.ApiResponseCashierMonthlyOrders{Status: res.GetStatus(), Message: res.GetMessage(), Data: data}
}
