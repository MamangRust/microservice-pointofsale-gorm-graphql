package cashiergraphqlmapper

import (
	graphqlmapper "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/mapper"
	"github.com/MamangRust/microservice-point-of-sale-apigateway/internal/model"
	cashierpb "github.com/MamangRust/microservice-point-of-sale-pb/cashier"
)

type cashierGraphqlMapper struct {
}

func NewCashierGraphqlMapper() *cashierGraphqlMapper {
	return &cashierGraphqlMapper{}
}

func (c *cashierGraphqlMapper) ToGraphqlResponseCashier(res *cashierpb.ApiResponseCashier) *model.APIResponseCashier {
	return &model.APIResponseCashier{
		Status:  res.Status,
		Message: res.Message,
		Data:    c.mapResponseCashier(res.Data),
	}
}

func (c *cashierGraphqlMapper) ToGraphqlResponsesCashier(res *cashierpb.ApiResponsesCashier) *model.APIResponsesCashier {
	return &model.APIResponsesCashier{
		Status:  res.Status,
		Message: res.Message,
		Data:    c.mapResponsesCashier(res.Data),
	}
}

func (c *cashierGraphqlMapper) ToGraphqlResponseCashierDeleteAt(res *cashierpb.ApiResponseCashierDeleteAt) *model.APIResponseCashierDeleteAt {
	return &model.APIResponseCashierDeleteAt{
		Status:  res.Status,
		Message: res.Message,
		Data:    c.mapResponseCashierDeleteAt(res.Data),
	}
}

func (c *cashierGraphqlMapper) ToGraphqlResponseCashierDelete(res *cashierpb.ApiResponseCashierDelete) *model.APIResponseCashierDelete {
	return &model.APIResponseCashierDelete{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (c *cashierGraphqlMapper) ToGraphqlResponseCashierAll(res *cashierpb.ApiResponseCashierAll) *model.APIResponseCashierAll {
	return &model.APIResponseCashierAll{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (c *cashierGraphqlMapper) ToGraphqlResponsePaginationCashier(res *cashierpb.ApiResponsePaginationCashier) *model.APIResponsePaginationCashier {
	return &model.APIResponsePaginationCashier{
		Status:     res.Status,
		Message:    res.Message,
		Data:       c.mapResponsesCashier(res.Data),
		Pagination: graphqlmapper.MapPaginationMeta(res.Pagination),
	}
}

func (c *cashierGraphqlMapper) ToGraphqlResponsePaginationCashierDeleteAt(res *cashierpb.ApiResponsePaginationCashierDeleteAt) *model.APIResponsePaginationCashierDeleteAt {
	return &model.APIResponsePaginationCashierDeleteAt{
		Status:     res.Status,
		Message:    res.Message,
		Data:       c.mapResponsesCashierDeleteAt(res.Data),
		Pagination: graphqlmapper.MapPaginationMeta(res.Pagination),
	}
}

func (c *cashierGraphqlMapper) ToGraphqlResponseMonthlyTotalSales(res *cashierpb.ApiResponseCashierMonthlyTotalSales) *model.APIResponseCashierMonthlyTotalSales {
	return &model.APIResponseCashierMonthlyTotalSales{
		Status:  res.Status,
		Message: res.Message,
		Data:    c.mapResponsesCashierMonthlyTotalSales(res.Data),
	}
}

func (c *cashierGraphqlMapper) ToGraphqlResponseMonthlySales(res *cashierpb.ApiResponseCashierMonthSales) *model.APIResponseCashierMonthSales {
	return &model.APIResponseCashierMonthSales{
		Status:  res.Status,
		Message: res.Message,
		Data:    c.mapResponsesCashierMonthlySales(res.Data),
	}
}

func (c *cashierGraphqlMapper) ToGraphqlResponseYearlySales(res *cashierpb.ApiResponseCashierYearSales) *model.APIResponseCashierYearSales {
	return &model.APIResponseCashierYearSales{
		Status:  res.Status,
		Message: res.Message,
		Data:    c.mapResponsesCashierYearlySales(res.Data),
	}
}

func (c *cashierGraphqlMapper) ToGraphqlResponseYearlyTotalSales(res *cashierpb.ApiResponseCashierYearlyTotalSales) *model.APIResponseCashierYearlyTotalSales {
	return &model.APIResponseCashierYearlyTotalSales{
		Status:  res.Status,
		Message: res.Message,
		Data:    c.mapResponsesCashierYearlyTotalSales(res.Data),
	}
}

func (c *cashierGraphqlMapper) mapResponseCashier(cashier *cashierpb.CashierResponse) *model.CashierResponse {
	if cashier == nil {
		return nil
	}
	return &model.CashierResponse{
		ID:         int32(cashier.Id),
		MerchantID: int32(cashier.MerchantId),
		Name:       cashier.Name,
		CreatedAt:  &cashier.CreatedAt,
		UpdatedAt:  &cashier.UpdatedAt,
	}
}

func (c *cashierGraphqlMapper) mapResponsesCashier(cashiers []*cashierpb.CashierResponse) []*model.CashierResponse {
	var responses []*model.CashierResponse
	for _, cashier := range cashiers {
		responses = append(responses, c.mapResponseCashier(cashier))
	}
	return responses
}

func (c *cashierGraphqlMapper) mapResponseCashierDeleteAt(cashier *cashierpb.CashierResponseDeleteAt) *model.CashierResponseDeleteAt {
	if cashier == nil {
		return nil
	}
	var deletedAt string
	if cashier.DeletedAt != nil {
		deletedAt = cashier.DeletedAt.Value
	}

	return &model.CashierResponseDeleteAt{
		ID:         int32(cashier.Id),
		MerchantID: int32(cashier.MerchantId),
		Name:       cashier.Name,
		CreatedAt:  &cashier.CreatedAt,
		UpdatedAt:  &cashier.UpdatedAt,
		DeletedAt:  &deletedAt,
	}
}

func (c *cashierGraphqlMapper) mapResponsesCashierDeleteAt(cashiers []*cashierpb.CashierResponseDeleteAt) []*model.CashierResponseDeleteAt {
	var responses []*model.CashierResponseDeleteAt
	for _, cashier := range cashiers {
		responses = append(responses, c.mapResponseCashierDeleteAt(cashier))
	}
	return responses
}

func (c *cashierGraphqlMapper) mapResponseCashierMonthlySales(s *cashierpb.CashierResponseMonthSales) *model.CashierResponseMonthSales {
	if s == nil {
		return nil
	}
	return &model.CashierResponseMonthSales{
		Month:       s.Month,
		CashierID:   s.CashierId,
		CashierName: s.CashierName,
		OrderCount:  s.OrderCount,
		TotalSales:  s.TotalSales,
	}
}

func (c *cashierGraphqlMapper) mapResponsesCashierMonthlySales(s []*cashierpb.CashierResponseMonthSales) []*model.CashierResponseMonthSales {
	var responses []*model.CashierResponseMonthSales
	for _, cashier := range s {
		responses = append(responses, c.mapResponseCashierMonthlySales(cashier))
	}
	return responses
}

func (c *cashierGraphqlMapper) mapResponseCashierYearlySales(s *cashierpb.CashierResponseYearSales) *model.CashierResponseYearSales {
	if s == nil {
		return nil
	}
	return &model.CashierResponseYearSales{
		Year:        s.Year,
		CashierID:   s.CashierId,
		CashierName: s.CashierName,
		OrderCount:  s.OrderCount,
		TotalSales:  s.TotalSales,
	}
}

func (c *cashierGraphqlMapper) mapResponsesCashierYearlySales(s []*cashierpb.CashierResponseYearSales) []*model.CashierResponseYearSales {
	var responses []*model.CashierResponseYearSales
	for _, cashier := range s {
		responses = append(responses, c.mapResponseCashierYearlySales(cashier))
	}
	return responses
}

func (c *cashierGraphqlMapper) mapResponseCashierMonthlyTotalSales(s *cashierpb.CashierResponseMonthTotalSales) *model.CashierResponseMonthTotalSales {
	if s == nil {
		return nil
	}
	return &model.CashierResponseMonthTotalSales{
		Year:       s.Year,
		Month:      s.Month,
		TotalSales: s.TotalSales,
	}
}

func (c *cashierGraphqlMapper) mapResponsesCashierMonthlyTotalSales(s []*cashierpb.CashierResponseMonthTotalSales) []*model.CashierResponseMonthTotalSales {
	var responses []*model.CashierResponseMonthTotalSales
	for _, cashier := range s {
		responses = append(responses, c.mapResponseCashierMonthlyTotalSales(cashier))
	}
	return responses
}

func (c *cashierGraphqlMapper) mapResponseCashierYearlyTotalSales(s *cashierpb.CashierResponseYearTotalSales) *model.CashierResponseYearTotalSales {
	if s == nil {
		return nil
	}
	return &model.CashierResponseYearTotalSales{
		Year:       s.Year,
		TotalSales: s.TotalSales,
	}
}

func (c *cashierGraphqlMapper) mapResponsesCashierYearlyTotalSales(s []*cashierpb.CashierResponseYearTotalSales) []*model.CashierResponseYearTotalSales {
	var responses []*model.CashierResponseYearTotalSales
	for _, cashier := range s {
		responses = append(responses, c.mapResponseCashierYearlyTotalSales(cashier))
	}
	return responses
}
