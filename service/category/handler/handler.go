package handler

import (
	"github.com/MamangRust/microservice-point-of-sale-category/service"
	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/category"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Category pb.CategoryServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Category: NewCategoryHandleGrpc(
			deps.Service,
			deps.Logger,
		),
	}
}
