package handler

import (
	"context"
	"math"

	"github.com/MamangRust/microservice-point-of-sale-order/service"
	"github.com/MamangRust/microservice-point-of-sale-order/repository"
	"github.com/MamangRust/microservice-point-of-sale-pkg/database/models"
	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	"github.com/MamangRust/microservice-point-of-sale-shared/convert"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/order_errors"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/order"
	pbcommon "github.com/MamangRust/microservice-pointofsale-grpc/pb/common"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type orderHandleGrpc struct {
	pb.UnimplementedOrderServiceServer
	orderQuery           service.OrderQueryService
	orderCommand         service.OrderCommandService
	logger               logger.LoggerInterface
}

func NewOrderHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pb.OrderServiceServer {
	return &orderHandleGrpc{
		orderQuery:           service.OrderQuery,
		orderCommand:         service.OrderCommand,
		logger:               logger,
	}
}

func (s *orderHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllOrderRequest) (*pb.ApiResponsePaginationOrder, error) {
	s.logger.Info("FindAll orders called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrders{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchant, totalRecords, err := s.orderQuery.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAll orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindAll orders success")

	return &pb.ApiResponsePaginationOrder{
		Status:     "success",
		Message:    "Successfully fetched order",
		Data:       mapResponsesOrder(merchant),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *orderHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdOrderRequest) (*pb.ApiResponseOrder, error) {
	s.logger.Info("FindById order called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	order, err := s.orderQuery.FindById(ctx, id)
	if err != nil {
		s.logger.Error("FindById order failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindById order success")

	return &pb.ApiResponseOrder{
		Status:  "success",
		Message: "Successfully fetched order",
		Data:    mapResponseOrderFromModel(order),
	}, nil
}

func (s *orderHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllOrderRequest) (*pb.ApiResponsePaginationOrderDeleteAt, error) {
	s.logger.Info("FindByActive orders called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrders{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchant, totalRecords, err := s.orderQuery.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByActive orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByActive orders success")

	return &pb.ApiResponsePaginationOrderDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active order",
		Data:       mapResponsesOrderActive(merchant),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *orderHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllOrderRequest) (*pb.ApiResponsePaginationOrderDeleteAt, error) {
	s.logger.Info("FindByTrashed orders called")

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrders{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.orderQuery.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByTrashed orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByTrashed orders success")

	return &pb.ApiResponsePaginationOrderDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed order",
		Data:       mapResponsesOrderTrashed(users),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *orderHandleGrpc) Create(ctx context.Context, request *pb.CreateOrderRequest) (*pb.ApiResponseOrder, error) {
	s.logger.Info("Create order called", zap.Int32("merchantId", request.GetMerchantId()))

	req := &requests.CreateOrderRequest{
		MerchantID: int(request.GetMerchantId()),
		CashierID:  int(request.GetCashierId()),
	}

	for _, item := range request.GetItems() {
		req.Items = append(req.Items, requests.CreateOrderItemRequest{
			ProductID: int(item.GetProductId()),
			Quantity:  int(item.GetQuantity()),
		})
	}

	if err := req.Validate(); err != nil {
		return nil, order_errors.ErrGrpcValidateCreateOrder
	}

	order, err := s.orderCommand.CreateOrder(ctx, req)
	if err != nil {
		s.logger.Error("Create order failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Create order success")

	return &pb.ApiResponseOrder{
		Status:  "success",
		Message: "Successfully created order",
		Data:    mapResponseOrderFromModel(order),
	}, nil
}

func (s *orderHandleGrpc) Update(ctx context.Context, request *pb.UpdateOrderRequest) (*pb.ApiResponseOrder, error) {
	s.logger.Info("Update order called", zap.Int32("id", request.GetOrderId()))

	id := int(request.GetOrderId())
	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	req := &requests.UpdateOrderRequest{
		OrderID: &id,
	}

	for _, item := range request.GetItems() {
		req.Items = append(req.Items, requests.UpdateOrderItemRequest{
			OrderItemID: int(item.GetOrderItemId()),
			ProductID:   int(item.GetProductId()),
			Quantity:    int(item.GetQuantity()),
		})
	}

	if err := req.Validate(); err != nil {
		return nil, order_errors.ErrGrpcValidateUpdateOrder
	}

	order, err := s.orderCommand.UpdateOrder(ctx, req)
	if err != nil {
		s.logger.Error("Update order failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Update order success")

	return &pb.ApiResponseOrder{
		Status:  "success",
		Message: "Successfully updated order",
		Data:    mapResponseOrderFromModel(order),
	}, nil
}

func (s *orderHandleGrpc) TrashedOrder(ctx context.Context, request *pb.FindByIdOrderRequest) (*pb.ApiResponseOrderDeleteAt, error) {
	s.logger.Info("TrashedOrder called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	merchant, err := s.orderCommand.TrashedOrder(ctx, id)
	if err != nil {
		s.logger.Error("TrashedOrder failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("TrashedOrder success")

	return &pb.ApiResponseOrderDeleteAt{
		Status:  "success",
		Message: "Successfully trashed order",
		Data:    mapResponseOrderDeleteAtFromModel(merchant),
	}, nil
}

func (s *orderHandleGrpc) RestoreOrder(ctx context.Context, request *pb.FindByIdOrderRequest) (*pb.ApiResponseOrderDeleteAt, error) {
	s.logger.Info("RestoreOrder called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	merchant, err := s.orderCommand.RestoreOrder(ctx, id)
	if err != nil {
		s.logger.Error("RestoreOrder failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreOrder success")

	return &pb.ApiResponseOrderDeleteAt{
		Status:  "success",
		Message: "Successfully restored order",
		Data:    mapResponseOrderDeleteAtFromModel(merchant),
	}, nil
}

func (s *orderHandleGrpc) DeleteOrderPermanent(ctx context.Context, request *pb.FindByIdOrderRequest) (*pb.ApiResponseOrderDelete, error) {
	s.logger.Info("DeleteOrderPermanent called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	_, err := s.orderCommand.DeleteOrderPermanent(ctx, id)
	if err != nil {
		s.logger.Error("DeleteOrderPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteOrderPermanent success")

	return &pb.ApiResponseOrderDelete{
		Status:  "success",
		Message: "Successfully deleted order permanently",
	}, nil
}

func (s *orderHandleGrpc) RestoreAllOrder(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseOrderAll, error) {
	s.logger.Info("RestoreAllOrder called")

	_, err := s.orderCommand.RestoreAllOrder(ctx)
	if err != nil {
		s.logger.Error("RestoreAllOrder failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreAllOrder success")

	return &pb.ApiResponseOrderAll{
		Status:  "success",
		Message: "Successfully restore all order",
	}, nil
}

func (s *orderHandleGrpc) DeleteAllOrderPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseOrderAll, error) {
	s.logger.Info("DeleteAllOrderPermanent called")

	_, err := s.orderCommand.DeleteAllOrderPermanent(ctx)
	if err != nil {
		s.logger.Error("DeleteAllOrderPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteAllOrderPermanent success")

	return &pb.ApiResponseOrderAll{
		Status:  "success",
		Message: "Successfully delete all order permanen",
	}, nil
}

// Map helpers

func mapPaginationMeta(meta *pbcommon.PaginationMeta) *pbcommon.PaginationMeta {
	if meta == nil {
		return nil
	}
	return &pbcommon.PaginationMeta{
		CurrentPage:  meta.CurrentPage,
		PageSize:     meta.PageSize,
		TotalPages:   meta.TotalPages,
		TotalRecords: meta.TotalRecords,
	}
}


func mapResponseOrderFromModel(order *models.Order) *pb.OrderResponse {
	if order == nil {
		return nil
	}
	return &pb.OrderResponse{
		Id:         order.OrderID,
		MerchantId: order.MerchantID,
		CashierId:  order.CashierID,
		TotalPrice: int32(order.TotalPrice),
		CreatedAt:  convert.FormatTimePtr(order.CreatedAt),
		UpdatedAt:  convert.FormatTimePtr(order.UpdatedAt),
	}
}

func mapResponseOrderDeleteAtFromModel(order *models.Order) *pb.OrderResponseDeleteAt {
	if order == nil {
		return nil
	}
	return &pb.OrderResponseDeleteAt{
		Id:         order.OrderID,
		MerchantId: order.MerchantID,
		CashierId:  order.CashierID,
		TotalPrice: int32(order.TotalPrice),
		CreatedAt:  convert.FormatTimePtr(order.CreatedAt),
		UpdatedAt:  convert.FormatTimePtr(order.UpdatedAt),
		DeletedAt:  convert.TimeToWrappers(order.DeletedAt),
	}
}

func mapResponsesOrder(orders []*repository.OrderResult) []*pb.OrderResponse {
	var mappedOrders []*pb.OrderResponse
	for _, order := range orders {
		if order == nil {
			continue
		}
		mappedOrders = append(mappedOrders, &pb.OrderResponse{
			Id:         order.OrderID,
			MerchantId: order.MerchantID,
			CashierId:  order.CashierID,
			TotalPrice: int32(order.TotalPrice),
			CreatedAt:  order.CreatedAt,
			UpdatedAt:  order.UpdatedAt,
		})
	}
	return mappedOrders
}

func mapResponsesOrderActive(orders []*repository.OrderResultDeleteAt) []*pb.OrderResponseDeleteAt {
	var mappedOrders []*pb.OrderResponseDeleteAt
	for _, order := range orders {
		if order == nil {
			continue
		}
		deletedAt := convert.StrValToWrappers(&order.DeletedAt)
		mappedOrders = append(mappedOrders, &pb.OrderResponseDeleteAt{
			Id:         order.OrderID,
			MerchantId: order.MerchantID,
			CashierId:  order.CashierID,
			TotalPrice: int32(order.TotalPrice),
			CreatedAt:  order.CreatedAt,
			UpdatedAt:  order.UpdatedAt,
			DeletedAt:  deletedAt,
		})
	}
	return mappedOrders
}

func mapResponsesOrderTrashed(orders []*repository.OrderResultDeleteAt) []*pb.OrderResponseDeleteAt {
	return mapResponsesOrderActive(orders)
}












