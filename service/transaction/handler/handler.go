package handler

import (
	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/transaction"
	"github.com/MamangRust/microservice-point-of-sale-transacton/service"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Transaction pb.TransactionServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Transaction: NewTransactionHandleGrpc(
			deps.Service,
			deps.Logger,
		),
	}
}
