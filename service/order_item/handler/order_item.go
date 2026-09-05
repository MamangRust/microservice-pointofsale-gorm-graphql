package handler

import (
	"context"
	"math"

	"github.com/MamangRust/microservice-point-of-sale-order-item/service"
	"github.com/MamangRust/microservice-point-of-sale-order-item/repository"
	"github.com/MamangRust/microservice-point-of-sale-pkg/database/models"
	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	"github.com/MamangRust/microservice-point-of-sale-shared/convert"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors"
	orderitem_errors "github.com/MamangRust/microservice-point-of-sale-shared/errors/order_item_errors"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/order_item"
	pbcommon "github.com/MamangRust/microservice-pointofsale-grpc/pb/common"
	"go.uber.org/zap"
)

type orderItemHandleGrpc struct {
	pb.UnimplementedOrderItemServiceServer
	orderItemService service.OrderItemQueryService
	logger           logger.LoggerInterface
}

func NewOrderItemHandleGrpc(
	orderItemService service.OrderItemQueryService,
	logger logger.LoggerInterface,
) pb.OrderItemServiceServer {
	return &orderItemHandleGrpc{
		orderItemService: orderItemService,
		logger:           logger,
	}
}

func (s *orderItemHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllOrderItemRequest) (*pb.ApiResponsePaginationOrderItem, error) {
	s.logger.Info("FindAll order items called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrderItems{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orderItems, totalRecords, err := s.orderItemService.FindAllOrderItems(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAll order items failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindAll order items success")

	return &pb.ApiResponsePaginationOrderItem{
		Status:     "success",
		Message:    "Successfully fetched order items",
		Data:       mapResponsesOrderItem(orderItems),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *orderItemHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllOrderItemRequest) (*pb.ApiResponsePaginationOrderItemDeleteAt, error) {
	s.logger.Info("FindByActive order items called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrderItems{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orderItems, totalRecords, err := s.orderItemService.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByActive order items failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByActive order items success")

	return &pb.ApiResponsePaginationOrderItemDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active order items",
		Data:       mapResponsesOrderItemActive(orderItems),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *orderItemHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllOrderItemRequest) (*pb.ApiResponsePaginationOrderItemDeleteAt, error) {
	s.logger.Info("FindByTrashed order items called")

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrderItems{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orderItems, totalRecords, err := s.orderItemService.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByTrashed order items failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByTrashed order items success")

	return &pb.ApiResponsePaginationOrderItemDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed order items",
		Data:       mapResponsesOrderItemTrashed(orderItems),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *orderItemHandleGrpc) FindOrderItemByOrder(ctx context.Context, request *pb.FindByIdOrderItemRequest) (*pb.ApiResponsesOrderItem, error) {
	s.logger.Info("FindOrderItemByOrder called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, orderitem_errors.ErrGrpcInvalidID
	}

	orderItems, err := s.orderItemService.FindOrderItemByOrder(ctx, id)
	if err != nil {
		s.logger.Error("FindOrderItemByOrder failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindOrderItemByOrder success")

	return &pb.ApiResponsesOrderItem{
		Status:  "success",
		Message: "Successfully fetched order items by order",
		Data:    mapResponsesOrderItemFromModel(orderItems),
	}, nil
}

// Map helpers

// Map helpers
func mapPaginationMeta(meta *pbcommon.PaginationMeta) *pbcommon.PaginationMeta {
	if meta == nil { return nil }
	return &pbcommon.PaginationMeta{
		CurrentPage: meta.CurrentPage, PageSize: meta.PageSize,
		TotalPages: meta.TotalPages, TotalRecords: meta.TotalRecords,
	}
}

func fmtOIStr(s *string) string {
	return convert.StrVal(s)
}


func mapOrderItemResponse(item *repository.OrderItemResult) *pb.OrderItemResponse {
	if item == nil { return nil }
	return &pb.OrderItemResponse{
		Id: item.OrderItemID, OrderId: item.OrderID, ProductId: item.ProductID,
		Quantity: item.Quantity, Price: int32(item.Price),
		CreatedAt: fmtOIStr(item.CreatedAt), UpdatedAt: fmtOIStr(item.UpdatedAt),
	}
}

func mapResponsesOrderItem(items []*repository.OrderItemResult) []*pb.OrderItemResponse {
	var res []*pb.OrderItemResponse
	for _, i := range items { res = append(res, mapOrderItemResponse(i)) }
	return res
}

func mapOrderItemDeleteAt(item *repository.OrderItemResultDeleteAt) *pb.OrderItemResponseDeleteAt {
	if item == nil { return nil }
	return &pb.OrderItemResponseDeleteAt{
		Id: item.OrderItemID, OrderId: item.OrderID, ProductId: item.ProductID,
		Quantity: item.Quantity, Price: int32(item.Price),
		CreatedAt: fmtOIStr(item.CreatedAt), UpdatedAt: fmtOIStr(item.UpdatedAt),
		DeletedAt: convert.StrValToWrappers(item.DeletedAt),
	}
}

func mapResponsesOrderItemActive(items []*repository.OrderItemResultDeleteAt) []*pb.OrderItemResponseDeleteAt {
	var res []*pb.OrderItemResponseDeleteAt
	for _, i := range items { res = append(res, mapOrderItemDeleteAt(i)) }
	return res
}

func mapResponsesOrderItemTrashed(items []*repository.OrderItemResultDeleteAt) []*pb.OrderItemResponseDeleteAt {
	return mapResponsesOrderItemActive(items)
}

func mapOrderItemModel(item *models.OrderItem) *pb.OrderItemResponse {
	if item == nil { return nil }
	return &pb.OrderItemResponse{
		Id: item.OrderItemID, OrderId: item.OrderID, ProductId: item.ProductID,
		Quantity: item.Quantity, Price: int32(item.Price),
		CreatedAt: convert.FormatTimePtr(item.CreatedAt), UpdatedAt: convert.FormatTimePtr(item.UpdatedAt),
	}
}

func mapResponsesOrderItemFromModel(items []*models.OrderItem) []*pb.OrderItemResponse {
	var res []*pb.OrderItemResponse
	for _, i := range items { res = append(res, mapOrderItemModel(i)) }
	return res
}
