package ordergraphqlmapper

import (
	"github.com/MamangRust/microservice-point-of-sale-apigateway/internal/model"
	orderpb "github.com/MamangRust/microservice-point-of-sale-pb/order"
)

type OrderGraphqlMapper interface {
	ToGraphqlResponseOrder(res *orderpb.ApiResponseOrder) *model.APIResponseOrder
	ToGraphqlResponsesOrder(res *orderpb.ApiResponsesOrder) *model.APIResponsesOrder
	ToGraphqlResponseOrderDeleteAt(res *orderpb.ApiResponseOrderDeleteAt) *model.APIResponseOrderDeleteAt
	ToGraphqlResponseOrderDelete(res *orderpb.ApiResponseOrderDelete) *model.APIResponseOrderDelete
	ToGraphqlResponseOrderAll(res *orderpb.ApiResponseOrderAll) *model.APIResponseOrderAll
	ToGraphqlResponsePaginationOrder(res *orderpb.ApiResponsePaginationOrder) *model.APIResponsePaginationOrder
	ToGraphqlResponsePaginationOrderDeleteAt(res *orderpb.ApiResponsePaginationOrderDeleteAt) *model.APIResponsePaginationOrderDeleteAt
	ToGraphqlResponseMonthlyRevenue(res *orderpb.ApiResponseOrderMonthly) *model.APIResponseOrderMonthly
	ToGraphqlResponseYearlyRevenue(res *orderpb.ApiResponseOrderYearly) *model.APIResponseOrderYearly
	ToGraphqlResponseMonthlyTotalRevenue(res *orderpb.ApiResponseOrderMonthlyTotalRevenue) *model.APIResponseOrderMonthlyTotalRevenue
	ToGraphqlResponseYearlyTotalRevenue(res *orderpb.ApiResponseOrderYearlyTotalRevenue) *model.APIResponseOrderYearlyTotalRevenue
}
