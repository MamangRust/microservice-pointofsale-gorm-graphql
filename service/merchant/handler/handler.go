package handler

import (
	"github.com/MamangRust/microservice-point-of-sale-merchant/service"
	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/merchant"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Merchant         pb.MerchantServiceServer
	MerchantDocument pb.MerchantDocumentServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Merchant: NewMerchantHandleGrpc(
			deps.Service,
			deps.Logger,
		),
		MerchantDocument: NewMerchantDocumentHandleGrpc(
			deps.Service,
			deps.Logger,
		),
	}
}
