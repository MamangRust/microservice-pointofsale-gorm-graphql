package handler

import (
	"github.com/MamangRust/microservice-point-of-sale-order-item/service"
	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/order_item"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	OrderItem pb.OrderItemServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		OrderItem: NewOrderItemHandleGrpc(
			deps.Service.OrderItemQuery,
			deps.Logger,
		),
	}
}
