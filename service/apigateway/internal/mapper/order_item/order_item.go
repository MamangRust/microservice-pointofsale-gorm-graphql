package orderitemgraphqlmapper

import (
	graphqlmapper "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/mapper"
	"github.com/MamangRust/microservice-point-of-sale-apigateway/internal/model"
	orderitempb "github.com/MamangRust/microservice-point-of-sale-pb/order_item"
)

type orderItemGraphqlMapper struct {
}

func NewOrderItemGraphqlMapper() *orderItemGraphqlMapper {
	return &orderItemGraphqlMapper{}
}

func (o *orderItemGraphqlMapper) ToGraphqlResponseOrderItem(res *orderitempb.ApiResponseOrderItem) *model.APIResponseOrderItem {
	return &model.APIResponseOrderItem{
		Status:  res.Status,
		Message: res.Message,
		Data:    o.mapResponseOrderItem(res.Data),
	}
}

func (o *orderItemGraphqlMapper) ToGraphqlResponsesOrderItem(res *orderitempb.ApiResponsesOrderItem) *model.APIResponsesOrderItem {
	return &model.APIResponsesOrderItem{
		Status:  res.Status,
		Message: res.Message,
		Data:    o.mapResponsesOrderItem(res.Data),
	}
}

func (o *orderItemGraphqlMapper) ToGrapqhlResponseOrderItemDelete(res *orderitempb.ApiResponseOrderItemDelete) *model.APIResponseOrderItemDelete {
	return &model.APIResponseOrderItemDelete{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (o *orderItemGraphqlMapper) ToGrapqhlResponseOrderItemAll(res *orderitempb.ApiResponseOrderItemAll) *model.APIResponseOrderItemAll {
	return &model.APIResponseOrderItemAll{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (o *orderItemGraphqlMapper) ToGraphqlResponsePaginationOrderItem(res *orderitempb.ApiResponsePaginationOrderItem) *model.APIResponsePaginationOrderItem {
	return &model.APIResponsePaginationOrderItem{
		Status:     res.Status,
		Message:    res.Message,
		Data:       o.mapResponsesOrderItem(res.Data),
		Pagination: graphqlmapper.MapPaginationMeta(res.Pagination),
	}
}

func (o *orderItemGraphqlMapper) ToGraphqlResponsePaginationOrderItemDeleteAt(res *orderitempb.ApiResponsePaginationOrderItemDeleteAt) *model.APIResponsePaginationOrderItemDeleteAt {
	return &model.APIResponsePaginationOrderItemDeleteAt{
		Status:     res.Status,
		Message:    res.Message,
		Data:       o.mapResponsesOrderItemDeleteAt(res.Data),
		Pagination: graphqlmapper.MapPaginationMeta(res.Pagination),
	}
}

func (o *orderItemGraphqlMapper) mapResponseOrderItem(orderItem *orderitempb.OrderItemResponse) *model.OrderItemResponse {
	if orderItem == nil {
		return nil
	}
	return &model.OrderItemResponse{
		ID:        int32(orderItem.Id),
		OrderID:   int32(orderItem.OrderId),
		ProductID: int32(orderItem.ProductId),
		Quantity:  int32(orderItem.Quantity),
		Price:     int32(orderItem.Price),
		CreatedAt: &orderItem.CreatedAt,
		UpdatedAt: &orderItem.UpdatedAt,
	}
}

func (o *orderItemGraphqlMapper) mapResponsesOrderItem(orderItems []*orderitempb.OrderItemResponse) []*model.OrderItemResponse {
	var responses []*model.OrderItemResponse
	for _, orderitem := range orderItems {
		responses = append(responses, o.mapResponseOrderItem(orderitem))
	}
	return responses
}

func (o *orderItemGraphqlMapper) mapResponseOrderItemDeleteAt(orderItem *orderitempb.OrderItemResponseDeleteAt) *model.OrderItemResponseDeleteAt {
	if orderItem == nil {
		return nil
	}
	var deletedAt string
	if orderItem.DeletedAt != nil {
		deletedAt = orderItem.DeletedAt.Value
	}

	return &model.OrderItemResponseDeleteAt{
		ID:        int32(orderItem.Id),
		OrderID:   int32(orderItem.OrderId),
		ProductID: int32(orderItem.ProductId),
		Quantity:  int32(orderItem.Quantity),
		Price:     int32(orderItem.Price),
		CreatedAt: &orderItem.CreatedAt,
		UpdatedAt: &orderItem.UpdatedAt,
		DeletedAt: &deletedAt,
	}
}

func (o *orderItemGraphqlMapper) mapResponsesOrderItemDeleteAt(orderItems []*orderitempb.OrderItemResponseDeleteAt) []*model.OrderItemResponseDeleteAt {
	var responses []*model.OrderItemResponseDeleteAt
	for _, orderitem := range orderItems {
		responses = append(responses, o.mapResponseOrderItemDeleteAt(orderitem))
	}
	return responses
}
