package handler

import (
	"context"

	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb"
)

type AuthHandleGrpc interface {
	pb.AuthServiceServer
	LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.ApiResponseLogin, error)
	RegisterUser(ctx context.Context, req *pb.RegisterRequest) (*pb.ApiResponseRegister, error)
}
