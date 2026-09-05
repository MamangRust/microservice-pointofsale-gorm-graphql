package handler

import (
	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/user"
	"github.com/MamangRust/microservice-point-of-sale-user/service"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	User pb.UserServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		User: NewUserHandleGrpc(
			deps.Service,
			deps.Logger,
		),
	}
}
