package handler

import (
	"github.com/MamangRust/microservice-point-of-sale-cashier/service"
	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/cashier"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Cashier pb.CashierServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Cashier: NewCashierHandleGrpc(
			deps.Service,
			deps.Logger,
		),
	}
}
