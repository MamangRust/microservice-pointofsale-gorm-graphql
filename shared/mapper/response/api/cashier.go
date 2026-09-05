package response_api

import (
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/response"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/cashier"
	pbstats "github.com/MamangRust/microservice-pointofsale-grpc/pb/stats"
)

type cashierResponseMapper struct{}

func NewCashierResponseMapper() *cashierResponseMapper {
	return &cashierResponseMapper{}
}

func (c *cashierResponseMapper) ToResponseCashier(cashier *pb.CashierResponse) *response.CashierResponse {
	return &response.CashierResponse{
		ID:         int(cashier.Id),
		MerchantID: int(cashier.MerchantId),
		Name:       cashier.Name,
		CreatedAt:  cashier.CreatedAt,
		UpdatedAt:  cashier.UpdatedAt,
	}
}

func (c *cashierResponseMapper) ToResponsesCashier(cashiers []*pb.CashierResponse) []*response.CashierResponse {
	var mappedCashiers []*response.CashierResponse

	for _, cashier := range cashiers {
		mappedCashiers = append(mappedCashiers, c.ToResponseCashier(cashier))
	}

	return mappedCashiers
}

func (c *cashierResponseMapper) ToResponseCashierDeleteAt(cashier *pb.CashierResponseDeleteAt) *response.CashierResponseDeleteAt {
	var deletedAt string
	if cashier.DeletedAt != nil {
		deletedAt = cashier.DeletedAt.Value
	}

	return &response.CashierResponseDeleteAt{
		ID:         int(cashier.Id),
		MerchantID: int(cashier.MerchantId),
		Name:       cashier.Name,
		CreatedAt:  cashier.CreatedAt,
		UpdatedAt:  cashier.UpdatedAt,
		DeletedAt:  &deletedAt,
	}
}

func (c *cashierResponseMapper) ToResponsesCashierDeleteAt(cashiers []*pb.CashierResponseDeleteAt) []*response.CashierResponseDeleteAt {
	var mappedCashiers []*response.CashierResponseDeleteAt

	for _, cashier := range cashiers {
		mappedCashiers = append(mappedCashiers, c.ToResponseCashierDeleteAt(cashier))
	}

	return mappedCashiers
}

func (s *cashierResponseMapper) ToResponseCashierMonthlySale(cashier *pbstats.CashierResponseMonthSales) *response.CashierResponseMonthSales {
	return &response.CashierResponseMonthSales{
		Month:       cashier.Month,
		CashierID:   int(cashier.CashierId),
		CashierName: cashier.CashierName,
		OrderCount:  int(cashier.OrderCount),
		TotalSales:  int(cashier.TotalSales),
	}
}

func (s *cashierResponseMapper) ToResponseCashierMonthlySales(c []*pbstats.CashierResponseMonthSales) []*response.CashierResponseMonthSales {
	var cashierRecords []*response.CashierResponseMonthSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.ToResponseCashierMonthlySale(cashier))
	}

	return cashierRecords
}

func (s *cashierResponseMapper) ToResponseCashierYearlySale(cashier *pbstats.CashierResponseYearSales) *response.CashierResponseYearSales {
	return &response.CashierResponseYearSales{
		Year:        cashier.Year,
		CashierID:   int(cashier.CashierId),
		CashierName: cashier.CashierName,
		OrderCount:  int(cashier.OrderCount),
		TotalSales:  int(cashier.TotalSales),
	}
}

func (s *cashierResponseMapper) ToResponseCashierYearlySales(c []*pbstats.CashierResponseYearSales) []*response.CashierResponseYearSales {
	var cashierRecords []*response.CashierResponseYearSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.ToResponseCashierYearlySale(cashier))
	}

	return cashierRecords
}

func (s *cashierResponseMapper) ToResponseCashierMonthlyTotalSale(c *pbstats.CashierResponseMonthTotalSales) *response.CashierResponseMonthTotalSales {
	return &response.CashierResponseMonthTotalSales{
		Year:       c.Year,
		Month:      c.Month,
		TotalSales: int(c.TotalSales),
	}
}

func (s *cashierResponseMapper) ToResponseCashierMonthlyTotalSales(c []*pbstats.CashierResponseMonthTotalSales) []*response.CashierResponseMonthTotalSales {
	var cashierRecords []*response.CashierResponseMonthTotalSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.ToResponseCashierMonthlyTotalSale(cashier))
	}

	return cashierRecords
}

func (s *cashierResponseMapper) ToResponseCashierYearlyTotalSale(c *pbstats.CashierResponseYearTotalSales) *response.CashierResponseYearTotalSales {
	return &response.CashierResponseYearTotalSales{
		Year:       c.Year,
		TotalSales: int(c.TotalSales),
	}
}

func (s *cashierResponseMapper) ToResponseCashierYearlyTotalSales(c []*pbstats.CashierResponseYearTotalSales) []*response.CashierResponseYearTotalSales {
	var cashierRecords []*response.CashierResponseYearTotalSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.ToResponseCashierYearlyTotalSale(cashier))
	}

	return cashierRecords
}

func (c *cashierResponseMapper) ToApiResponseCashier(pbResponse *pb.ApiResponseCashier) *response.ApiResponseCashier {
	return &response.ApiResponseCashier{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    c.ToResponseCashier(pbResponse.Data),
	}
}

func (c *cashierResponseMapper) ToApiResponsesCashier(pbResponse *pb.ApiResponsesCashier) *response.ApiResponsesCashier {
	return &response.ApiResponsesCashier{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    c.ToResponsesCashier(pbResponse.Data),
	}
}
func (c *cashierResponseMapper) ToApiResponseCashierDeleteAt(pbResponse *pb.ApiResponseCashierDeleteAt) *response.ApiResponseCashierDeleteAt {
	return &response.ApiResponseCashierDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    c.ToResponseCashierDeleteAt(pbResponse.Data),
	}
}

func (c *cashierResponseMapper) ToApiResponseCashierDelete(pbResponse *pb.ApiResponseCashierDelete) *response.ApiResponseCashierDelete {
	return &response.ApiResponseCashierDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (c *cashierResponseMapper) ToApiResponseCashierAll(pbResponse *pb.ApiResponseCashierAll) *response.ApiResponseCashierAll {
	return &response.ApiResponseCashierAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (c *cashierResponseMapper) ToApiResponsePaginationCashierDeleteAt(pbResponse *pb.ApiResponsePaginationCashierDeleteAt) *response.ApiResponsePaginationCashierDeleteAt {
	return &response.ApiResponsePaginationCashierDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       c.ToResponsesCashierDeleteAt(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}

func (c *cashierResponseMapper) ToApiResponsePaginationCashier(pbResponse *pb.ApiResponsePaginationCashier) *response.ApiResponsePaginationCashier {
	return &response.ApiResponsePaginationCashier{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       c.ToResponsesCashier(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}

func (c *cashierResponseMapper) ToApiResponseCashierMonthlySale(pbResponse *pbstats.ApiResponseCashierMonthSales) *response.ApiResponseCashierMonthSales {
	return &response.ApiResponseCashierMonthSales{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    c.ToResponseCashierMonthlySales(pbResponse.Data),
	}
}

func (c *cashierResponseMapper) ToApiResponseCashierYearlySale(pbResponse *pbstats.ApiResponseCashierYearSales) *response.ApiResponseCashierYearSales {
	return &response.ApiResponseCashierYearSales{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    c.ToResponseCashierYearlySales(pbResponse.Data),
	}
}

func (u *cashierResponseMapper) ToApiResponseMonthlyTotalSales(pbResponse *pbstats.ApiResponseCashierMonthlyTotalSales) *response.ApiResponseCashierMonthlyTotalSales {
	return &response.ApiResponseCashierMonthlyTotalSales{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    u.ToResponseCashierMonthlyTotalSales(pbResponse.Data),
	}
}

func (u *cashierResponseMapper) ToApiResponseYearlyTotalSales(pbResponse *pbstats.ApiResponseCashierYearlyTotalSales) *response.ApiResponseCashierYearlyTotalSales {
	return &response.ApiResponseCashierYearlyTotalSales{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    u.ToResponseCashierYearlyTotalSales(pbResponse.Data),
	}
}
