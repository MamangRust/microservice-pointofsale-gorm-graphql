package protomapper

import (
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/response"
	pbcommon "github.com/MamangRust/microservice-pointofsale-grpc/pb/common"
	pbstats "github.com/MamangRust/microservice-pointofsale-grpc/pb/stats"

	"google.golang.org/protobuf/types/known/wrapperspb"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/cashier"
)

type cashierProtoMapper struct {
}

func NewCashierProtoMapper() *cashierProtoMapper {
	return &cashierProtoMapper{}
}

func (c *cashierProtoMapper) ToProtoResponseCashier(status string, message string, pbResponse *response.CashierResponse) *pb.ApiResponseCashier {
	return &pb.ApiResponseCashier{
		Status:  status,
		Message: message,
		Data:    c.mapResponseCashier(pbResponse),
	}
}

func (c *cashierProtoMapper) ToProtoResponsesCashier(status string, message string, pbResponse []*response.CashierResponse) *pb.ApiResponsesCashier {
	return &pb.ApiResponsesCashier{
		Status:  status,
		Message: message,
		Data:    c.mapResponsesCashier(pbResponse),
	}
}

func (c *cashierProtoMapper) ToProtoResponseCashierDeleteAt(status string, message string, pbResponse *response.CashierResponseDeleteAt) *pb.ApiResponseCashierDeleteAt {
	return &pb.ApiResponseCashierDeleteAt{
		Status:  status,
		Message: message,
		Data:    c.mapResponseCashierDeleteAt(pbResponse),
	}
}

func (c *cashierProtoMapper) ToProtoResponseCashierDelete(status string, message string) *pb.ApiResponseCashierDelete {
	return &pb.ApiResponseCashierDelete{
		Status:  status,
		Message: message,
	}
}

func (u *cashierProtoMapper) ToProtoResponseCashierAll(status string, message string) *pb.ApiResponseCashierAll {
	return &pb.ApiResponseCashierAll{
		Status:  status,
		Message: message,
	}
}

func (u *cashierProtoMapper) ToProtoResponsePaginationCashierDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, users []*response.CashierResponseDeleteAt) *pb.ApiResponsePaginationCashierDeleteAt {
	return &pb.ApiResponsePaginationCashierDeleteAt{
		Status:     status,
		Message:    message,
		Data:       u.mapResponsesCashierDeleteAt(users),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (u *cashierProtoMapper) ToProtoResponsePaginationCashier(pagination *pbcommon.PaginationMeta, status string, message string, users []*response.CashierResponse) *pb.ApiResponsePaginationCashier {
	return &pb.ApiResponsePaginationCashier{
		Status:     status,
		Message:    message,
		Data:       u.mapResponsesCashier(users),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (u *cashierProtoMapper) ToProtoResponseMonthlyTotalSales(status, message string, row []*response.CashierResponseMonthSales) *pbstats.ApiResponseCashierMonthSales {
	return &pbstats.ApiResponseCashierMonthSales{
		Status:  status,
		Message: message,
		Data:    u.mapResponsesCashierMonthlySales(row),
	}
}

func (u *cashierProtoMapper) ToProtoResponseYearlyTotalSales(status, message string, row []*response.CashierResponseYearSales) *pbstats.ApiResponseCashierYearSales {
	return &pbstats.ApiResponseCashierYearSales{
		Status:  status,
		Message: message,
		Data:    u.mapResponsesCashierYearlySales(row),
	}
}

func (u *cashierProtoMapper) ToProtoMonthlyTotalSales(status, message string, row []*response.CashierResponseMonthTotalSales) *pbstats.ApiResponseCashierMonthlyTotalSales {
	return &pbstats.ApiResponseCashierMonthlyTotalSales{
		Status:  status,
		Message: message,
		Data:    u.mapResponseCashierMonthlyTotalSales(row),
	}
}

func (u *cashierProtoMapper) ToProtoYearlyTotalSales(status, message string, row []*response.CashierResponseYearTotalSales) *pbstats.ApiResponseCashierYearlyTotalSales {
	return &pbstats.ApiResponseCashierYearlyTotalSales{
		Status:  status,
		Message: message,
		Data:    u.mapResponseCashierYearlyTotalSales(row),
	}
}

func (c *cashierProtoMapper) mapResponseCashier(cashier *response.CashierResponse) *pb.CashierResponse {
	return &pb.CashierResponse{
		Id:         int32(cashier.ID),
		MerchantId: int32(cashier.MerchantID),
		Name:       cashier.Name,
		CreatedAt:  cashier.CreatedAt,
		UpdatedAt:  cashier.UpdatedAt,
	}
}

func (c *cashierProtoMapper) mapResponsesCashier(cashiers []*response.CashierResponse) []*pb.CashierResponse {
	var mappedCashiers []*pb.CashierResponse

	for _, cashier := range cashiers {
		mappedCashiers = append(mappedCashiers, c.mapResponseCashier(cashier))
	}

	return mappedCashiers
}

func (c *cashierProtoMapper) mapResponseCashierDeleteAt(cashier *response.CashierResponseDeleteAt) *pb.CashierResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if cashier.DeletedAt != nil {
		deletedAt = wrapperspb.String(*cashier.DeletedAt)
	}

	return &pb.CashierResponseDeleteAt{
		Id:         int32(cashier.ID),
		MerchantId: int32(cashier.MerchantID),
		Name:       cashier.Name,
		CreatedAt:  cashier.CreatedAt,
		UpdatedAt:  cashier.UpdatedAt,
		DeletedAt:  deletedAt,
	}
}

func (c *cashierProtoMapper) mapResponsesCashierDeleteAt(cashiers []*response.CashierResponseDeleteAt) []*pb.CashierResponseDeleteAt {
	var mappedCashiers []*pb.CashierResponseDeleteAt

	for _, cashier := range cashiers {
		mappedCashiers = append(mappedCashiers, c.mapResponseCashierDeleteAt(cashier))
	}

	return mappedCashiers
}

func (s *cashierProtoMapper) mapResponseCashierMonthlySale(cashier *response.CashierResponseMonthSales) *pbstats.CashierResponseMonthSales {
	return &pbstats.CashierResponseMonthSales{
		Month:       cashier.Month,
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func (s *cashierProtoMapper) mapResponsesCashierMonthlySales(c []*response.CashierResponseMonthSales) []*pbstats.CashierResponseMonthSales {
	var cashierRecords []*pbstats.CashierResponseMonthSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.mapResponseCashierMonthlySale(cashier))
	}

	return cashierRecords
}

func (s *cashierProtoMapper) mapResponseCashierYearlySale(cashier *response.CashierResponseYearSales) *pbstats.CashierResponseYearSales {
	return &pbstats.CashierResponseYearSales{
		Year:        cashier.Year,
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func (s *cashierProtoMapper) mapResponsesCashierYearlySales(c []*response.CashierResponseYearSales) []*pbstats.CashierResponseYearSales {
	var cashierRecords []*pbstats.CashierResponseYearSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.mapResponseCashierYearlySale(cashier))
	}

	return cashierRecords
}

func (s *cashierProtoMapper) mapResponseCashierMonthlyTotalSale(c *response.CashierResponseMonthTotalSales) *pbstats.CashierResponseMonthTotalSales {
	return &pbstats.CashierResponseMonthTotalSales{
		Year:       c.Year,
		Month:      c.Month,
		TotalSales: int32(c.TotalSales),
	}
}

func (s *cashierProtoMapper) mapResponseCashierMonthlyTotalSales(c []*response.CashierResponseMonthTotalSales) []*pbstats.CashierResponseMonthTotalSales {
	var cashierRecords []*pbstats.CashierResponseMonthTotalSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.mapResponseCashierMonthlyTotalSale(cashier))
	}

	return cashierRecords
}

func (s *cashierProtoMapper) mapResponseCashierYearlyTotalSale(c *response.CashierResponseYearTotalSales) *pbstats.CashierResponseYearTotalSales {
	return &pbstats.CashierResponseYearTotalSales{
		Year:       c.Year,
		TotalSales: int32(c.TotalSales),
	}
}

func (s *cashierProtoMapper) mapResponseCashierYearlyTotalSales(c []*response.CashierResponseYearTotalSales) []*pbstats.CashierResponseYearTotalSales {
	var cashierRecords []*pbstats.CashierResponseYearTotalSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.mapResponseCashierYearlyTotalSale(cashier))
	}

	return cashierRecords
}
