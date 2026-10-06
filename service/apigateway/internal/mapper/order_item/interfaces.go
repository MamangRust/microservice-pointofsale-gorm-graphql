package orderitemgraphqlmapper

import (
	"github.com/MamangRust/microservice-point-of-sale-apigateway/internal/model"
	orderitempb "github.com/MamangRust/microservice-point-of-sale-pb/order_item"
)

type OrderItemGraphqlMapper interface {
	ToGraphqlResponseOrderItem(res *orderitempb.ApiResponseOrderItem) *model.APIResponseOrderItem
	ToGraphqlResponsesOrderItem(res *orderitempb.ApiResponsesOrderItem) *model.APIResponsesOrderItem
	ToGrapqhlResponseOrderItemDelete(res *orderitempb.ApiResponseOrderItemDelete) *model.APIResponseOrderItemDelete
	ToGrapqhlResponseOrderItemAll(res *orderitempb.ApiResponseOrderItemAll) *model.APIResponseOrderItemAll
	ToGraphqlResponsePaginationOrderItem(res *orderitempb.ApiResponsePaginationOrderItem) *model.APIResponsePaginationOrderItem
	ToGraphqlResponsePaginationOrderItemDeleteAt(res *orderitempb.ApiResponsePaginationOrderItemDeleteAt) *model.APIResponsePaginationOrderItemDeleteAt
}
