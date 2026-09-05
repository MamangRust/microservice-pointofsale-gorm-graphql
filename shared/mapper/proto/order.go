package protomapper

import (
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/response"
	pbcommon "github.com/MamangRust/microservice-pointofsale-grpc/pb/common"
	pbstats "github.com/MamangRust/microservice-pointofsale-grpc/pb/stats"

	"google.golang.org/protobuf/types/known/wrapperspb"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/order"
)

type orderProtoMapper struct{}

func NewOrderProtoMapper() *orderProtoMapper {
	return &orderProtoMapper{}
}

func (o *orderProtoMapper) ToProtoResponseOrder(status string, message string, pbResponse *response.OrderResponse) *pb.ApiResponseOrder {
	return &pb.ApiResponseOrder{
		Status:  status,
		Message: message,
		Data:    o.mapResponseOrder(pbResponse),
	}
}

func (o *orderProtoMapper) ToProtoResponsesOrder(status string, message string, pbResponse []*response.OrderResponse) *pb.ApiResponsesOrder {
	return &pb.ApiResponsesOrder{
		Status:  status,
		Message: message,
		Data:    o.mapResponsesOrder(pbResponse),
	}
}

func (o *orderProtoMapper) ToProtoResponseOrderDeleteAt(status string, message string, pbResponse *response.OrderResponseDeleteAt) *pb.ApiResponseOrderDeleteAt {
	return &pb.ApiResponseOrderDeleteAt{
		Status:  status,
		Message: message,
		Data:    o.mapResponseOrderDeleteAt(pbResponse),
	}
}

func (o *orderProtoMapper) ToProtoResponseOrderDelete(status string, message string) *pb.ApiResponseOrderDelete {
	return &pb.ApiResponseOrderDelete{
		Status:  status,
		Message: message,
	}
}

func (o *orderProtoMapper) ToProtoResponseOrderAll(status string, message string) *pb.ApiResponseOrderAll {
	return &pb.ApiResponseOrderAll{
		Status:  status,
		Message: message,
	}
}

func (o *orderProtoMapper) ToProtoResponsePaginationOrderDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, orders []*response.OrderResponseDeleteAt) *pb.ApiResponsePaginationOrderDeleteAt {
	return &pb.ApiResponsePaginationOrderDeleteAt{
		Status:     status,
		Message:    message,
		Data:       o.mapResponsesOrderDeleteAt(orders),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (o *orderProtoMapper) ToProtoResponsePaginationOrder(pagination *pbcommon.PaginationMeta, status string, message string, orders []*response.OrderResponse) *pb.ApiResponsePaginationOrder {
	return &pb.ApiResponsePaginationOrder{
		Status:     status,
		Message:    message,
		Data:       o.mapResponsesOrder(orders),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (o *orderProtoMapper) ToProtoResponseMonthlyRevenue(status string, message string, row []*response.OrderMonthlyResponse) *pbstats.ApiResponseOrderMonthly {
	return &pbstats.ApiResponseOrderMonthly{
		Status:  status,
		Message: message,
		Data:    o.mapResponsesOrderMonthlyPrices(row),
	}
}

func (o *orderProtoMapper) ToProtoResponseYearlyRevenue(status string, message string, row []*response.OrderYearlyResponse) *pbstats.ApiResponseOrderYearly {
	return &pbstats.ApiResponseOrderYearly{
		Status:  status,
		Message: message,
		Data:    o.mapResponsesOrderYearlyPrices(row),
	}
}

func (o *orderProtoMapper) ToProtoResponseMonthlyTotalRevenue(status string, message string, row []*response.OrderMonthlyTotalRevenueResponse) *pbstats.ApiResponseOrderMonthlyTotalRevenue {
	return &pbstats.ApiResponseOrderMonthlyTotalRevenue{
		Status:  status,
		Message: message,
		Data:    o.mapResponseOrderMonthlyTotalRevenues(row),
	}
}

func (o *orderProtoMapper) ToProtoResponseYearlyTotalRevenue(status string, message string, row []*response.OrderYearlyTotalRevenueResponse) *pbstats.ApiResponseOrderYearlyTotalRevenue {
	return &pbstats.ApiResponseOrderYearlyTotalRevenue{
		Status:  status,
		Message: message,
		Data:    o.mapResponseOrderYearlyTotalRevenues(row),
	}
}

func (o *orderProtoMapper) mapResponseOrder(order *response.OrderResponse) *pb.OrderResponse {
	return &pb.OrderResponse{
		Id:         int32(order.ID),
		MerchantId: int32(order.MerchantID),
		CashierId:  int32(order.CashierID),
		TotalPrice: int32(order.TotalPrice),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}
}

func (o *orderProtoMapper) mapResponsesOrder(orders []*response.OrderResponse) []*pb.OrderResponse {
	var mappedOrders []*pb.OrderResponse

	for _, order := range orders {
		mappedOrders = append(mappedOrders, o.mapResponseOrder(order))
	}

	return mappedOrders
}

func (o *orderProtoMapper) mapResponseOrderDeleteAt(order *response.OrderResponseDeleteAt) *pb.OrderResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue

	if order.DeleteAt != nil {
		deletedAt = wrapperspb.String(*order.DeleteAt)
	}

	return &pb.OrderResponseDeleteAt{
		Id:         int32(order.ID),
		MerchantId: int32(order.MerchantID),
		CashierId:  int32(order.CashierID),
		TotalPrice: int32(order.TotalPrice),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
		DeletedAt:  deletedAt,
	}
}

func (o *orderProtoMapper) mapResponsesOrderDeleteAt(orders []*response.OrderResponseDeleteAt) []*pb.OrderResponseDeleteAt {
	var mappedOrders []*pb.OrderResponseDeleteAt

	for _, order := range orders {
		mappedOrders = append(mappedOrders, o.mapResponseOrderDeleteAt(order))
	}

	return mappedOrders
}

func (s *orderProtoMapper) mapResponseOrderMonthlyPrice(category *response.OrderMonthlyResponse) *pbstats.OrderMonthlyResponse {
	return &pbstats.OrderMonthlyResponse{
		Month:          category.Month,
		OrderCount:     int32(category.OrderCount),
		TotalRevenue:   int32(category.TotalRevenue),
		TotalItemsSold: int32(category.TotalItemsSold),
	}
}

func (s *orderProtoMapper) mapResponsesOrderMonthlyPrices(c []*response.OrderMonthlyResponse) []*pbstats.OrderMonthlyResponse {
	var categoryRecords []*pbstats.OrderMonthlyResponse

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.mapResponseOrderMonthlyPrice(category))
	}

	return categoryRecords
}

func (s *orderProtoMapper) mapResponseOrderYearlyPrice(category *response.OrderYearlyResponse) *pbstats.OrderYearlyResponse {
	return &pbstats.OrderYearlyResponse{
		Year:               category.Year,
		OrderCount:         int32(category.OrderCount),
		TotalRevenue:       int32(category.TotalRevenue),
		TotalItemsSold:     int32(category.TotalItemsSold),
		ActiveCashiers:     int32(category.ActiveCashiers),
		UniqueProductsSold: int32(category.UniqueProductsSold),
	}
}

func (s *orderProtoMapper) mapResponsesOrderYearlyPrices(c []*response.OrderYearlyResponse) []*pbstats.OrderYearlyResponse {
	var categoryRecords []*pbstats.OrderYearlyResponse

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.mapResponseOrderYearlyPrice(category))
	}

	return categoryRecords
}

func (s *orderProtoMapper) mapResponseOrderMonthlyTotalRevenue(c *response.OrderMonthlyTotalRevenueResponse) *pbstats.OrderMonthlyTotalRevenueResponse {
	return &pbstats.OrderMonthlyTotalRevenueResponse{
		Year:           c.Year,
		Month:          c.Month,
		TotalRevenue:   int32(c.TotalRevenue),
		TotalItemsSold: int32(c.TotalItemsSold),
	}
}

func (s *orderProtoMapper) mapResponseOrderMonthlyTotalRevenues(c []*response.OrderMonthlyTotalRevenueResponse) []*pbstats.OrderMonthlyTotalRevenueResponse {
	var orderRecords []*pbstats.OrderMonthlyTotalRevenueResponse

	for _, row := range c {
		orderRecords = append(orderRecords, s.mapResponseOrderMonthlyTotalRevenue(row))
	}

	return orderRecords
}

func (s *orderProtoMapper) mapResponseOrderYearlyTotalRevenue(c *response.OrderYearlyTotalRevenueResponse) *pbstats.OrderYearlyTotalRevenueResponse {
	return &pbstats.OrderYearlyTotalRevenueResponse{
		Year:         c.Year,
		TotalRevenue: int32(c.TotalRevenue),
	}
}

func (s *orderProtoMapper) mapResponseOrderYearlyTotalRevenues(c []*response.OrderYearlyTotalRevenueResponse) []*pbstats.OrderYearlyTotalRevenueResponse {
	var orderRecords []*pbstats.OrderYearlyTotalRevenueResponse

	for _, row := range c {
		orderRecords = append(orderRecords, s.mapResponseOrderYearlyTotalRevenue(row))
	}

	return orderRecords
}
