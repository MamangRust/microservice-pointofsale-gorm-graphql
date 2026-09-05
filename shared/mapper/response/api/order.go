package response_api

import (
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/response"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/order"
	pbstats "github.com/MamangRust/microservice-pointofsale-grpc/pb/stats"
)

type orderResponseMapper struct {
}

func NewOrderResponseMapper() *orderResponseMapper {
	return &orderResponseMapper{}
}

func (o *orderResponseMapper) ToResponseOrder(order *pb.OrderResponse) *response.OrderResponse {
	return &response.OrderResponse{
		ID:         int(order.Id),
		MerchantID: int(order.MerchantId),
		CashierID:  int(order.CashierId),
		TotalPrice: int(order.TotalPrice),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}
}

func (o *orderResponseMapper) ToResponsesOrder(orders []*pb.OrderResponse) []*response.OrderResponse {
	var mappedOrders []*response.OrderResponse

	for _, order := range orders {
		mappedOrders = append(mappedOrders, o.ToResponseOrder(order))
	}

	return mappedOrders
}

func (o *orderResponseMapper) ToResponseOrderDeleteAt(order *pb.OrderResponseDeleteAt) *response.OrderResponseDeleteAt {
	var deletedAt string

	if order.DeletedAt != nil {
		deletedAt = order.DeletedAt.Value
	}

	return &response.OrderResponseDeleteAt{
		ID:         int(order.Id),
		MerchantID: int(order.MerchantId),
		CashierID:  int(order.CashierId),
		TotalPrice: int(order.TotalPrice),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
		DeleteAt:   &deletedAt,
	}
}

func (o *orderResponseMapper) ToResponsesOrderDeleteAt(orders []*pb.OrderResponseDeleteAt) []*response.OrderResponseDeleteAt {
	var mappedOrders []*response.OrderResponseDeleteAt

	for _, order := range orders {
		mappedOrders = append(mappedOrders, o.ToResponseOrderDeleteAt(order))
	}

	return mappedOrders
}

func (s *orderResponseMapper) ToOrderMonthlyPrice(category *pbstats.OrderMonthlyResponse) *response.OrderMonthlyResponse {
	return &response.OrderMonthlyResponse{
		Month:          category.Month,
		OrderCount:     int(category.OrderCount),
		TotalRevenue:   int(category.TotalRevenue),
		TotalItemsSold: int(category.TotalItemsSold),
	}
}

func (s *orderResponseMapper) ToOrderMonthlyPrices(c []*pbstats.OrderMonthlyResponse) []*response.OrderMonthlyResponse {
	var categoryRecords []*response.OrderMonthlyResponse

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.ToOrderMonthlyPrice(category))
	}

	return categoryRecords
}

func (s *orderResponseMapper) ToOrderYearlyPrice(category *pbstats.OrderYearlyResponse) *response.OrderYearlyResponse {
	return &response.OrderYearlyResponse{
		Year:               category.Year,
		OrderCount:         int(category.OrderCount),
		TotalRevenue:       int(category.TotalRevenue),
		TotalItemsSold:     int(category.TotalItemsSold),
		ActiveCashiers:     int(category.ActiveCashiers),
		UniqueProductsSold: int(category.UniqueProductsSold),
	}
}

func (s *orderResponseMapper) ToOrderYearlyPrices(c []*pbstats.OrderYearlyResponse) []*response.OrderYearlyResponse {
	var categoryRecords []*response.OrderYearlyResponse

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.ToOrderYearlyPrice(category))
	}

	return categoryRecords
}

func (s *orderResponseMapper) ToResponseOrderMonthlyTotalRevenue(c *pbstats.OrderMonthlyTotalRevenueResponse) *response.OrderMonthlyTotalRevenueResponse {
	return &response.OrderMonthlyTotalRevenueResponse{
		Year:           c.Year,
		Month:          c.Month,
		TotalRevenue:   int(c.TotalRevenue),
		TotalItemsSold: int(c.TotalItemsSold),
	}
}

func (s *orderResponseMapper) ToResponseOrderMonthlyTotalRevenues(c []*pbstats.OrderMonthlyTotalRevenueResponse) []*response.OrderMonthlyTotalRevenueResponse {
	var orderRecords []*response.OrderMonthlyTotalRevenueResponse

	for _, row := range c {
		orderRecords = append(orderRecords, s.ToResponseOrderMonthlyTotalRevenue(row))
	}

	return orderRecords
}

func (s *orderResponseMapper) ToResponseOrderYearlyTotalRevenue(c *pbstats.OrderYearlyTotalRevenueResponse) *response.OrderYearlyTotalRevenueResponse {
	return &response.OrderYearlyTotalRevenueResponse{
		Year:         c.Year,
		TotalRevenue: int(c.TotalRevenue),
	}
}

func (s *orderResponseMapper) ToResponseOrderYearlyTotalRevenues(c []*pbstats.OrderYearlyTotalRevenueResponse) []*response.OrderYearlyTotalRevenueResponse {
	var orderRecords []*response.OrderYearlyTotalRevenueResponse

	for _, row := range c {
		orderRecords = append(orderRecords, s.ToResponseOrderYearlyTotalRevenue(row))
	}

	return orderRecords
}

func (o *orderResponseMapper) ToApiResponseOrder(pbResponse *pb.ApiResponseOrder) *response.ApiResponseOrder {
	return &response.ApiResponseOrder{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToResponseOrder(pbResponse.Data),
	}
}

func (o *orderResponseMapper) ToApiResponseOrderDeleteAt(pbResponse *pb.ApiResponseOrderDeleteAt) *response.ApiResponseOrderDeleteAt {
	return &response.ApiResponseOrderDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToResponseOrderDeleteAt(pbResponse.Data),
	}
}

func (o *orderResponseMapper) ToApiResponsesOrder(pbResponse *pb.ApiResponsesOrder) *response.ApiResponsesOrder {
	return &response.ApiResponsesOrder{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToResponsesOrder(pbResponse.Data),
	}
}

func (o *orderResponseMapper) ToApiResponseOrderDelete(pbResponse *pb.ApiResponseOrderDelete) *response.ApiResponseOrderDelete {
	return &response.ApiResponseOrderDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (o *orderResponseMapper) ToApiResponseOrderAll(pbResponse *pb.ApiResponseOrderAll) *response.ApiResponseOrderAll {
	return &response.ApiResponseOrderAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (o *orderResponseMapper) ToApiResponsePaginationOrderDeleteAt(pbResponse *pb.ApiResponsePaginationOrderDeleteAt) *response.ApiResponsePaginationOrderDeleteAt {
	return &response.ApiResponsePaginationOrderDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       o.ToResponsesOrderDeleteAt(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}

func (o *orderResponseMapper) ToApiResponsePaginationOrder(pbResponse *pb.ApiResponsePaginationOrder) *response.ApiResponsePaginationOrder {
	return &response.ApiResponsePaginationOrder{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       o.ToResponsesOrder(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}

func (o *orderResponseMapper) ToApiResponseMonthlyOrder(pbResponse *pbstats.ApiResponseOrderMonthly) *response.ApiResponseOrderMonthly {
	return &response.ApiResponseOrderMonthly{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToOrderMonthlyPrices(pbResponse.Data),
	}
}

func (o *orderResponseMapper) ToApiResponseYearlyOrder(pbResponse *pbstats.ApiResponseOrderYearly) *response.ApiResponseOrderYearly {
	return &response.ApiResponseOrderYearly{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToOrderYearlyPrices(pbResponse.Data),
	}
}

func (o *orderResponseMapper) ToApiResponseMonthlyTotalRevenue(pbResponse *pbstats.ApiResponseOrderMonthlyTotalRevenue) *response.ApiResponseOrderMonthlyTotalRevenue {
	return &response.ApiResponseOrderMonthlyTotalRevenue{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToResponseOrderMonthlyTotalRevenues(pbResponse.Data),
	}
}

func (o *orderResponseMapper) ToApiResponseYearlyTotalRevenue(pbResponse *pbstats.ApiResponseOrderYearlyTotalRevenue) *response.ApiResponseOrderYearlyTotalRevenue {
	return &response.ApiResponseOrderYearlyTotalRevenue{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToResponseOrderYearlyTotalRevenues(pbResponse.Data),
	}
}
