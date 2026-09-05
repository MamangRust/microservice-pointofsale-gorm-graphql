package handler

import (
	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	"github.com/MamangRust/microservice-point-of-sale-role/service"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/role"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Role pb.RoleServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Role: NewRoleHandleGrpc(
			deps.Service,
			deps.Logger,
		),
	}
}
