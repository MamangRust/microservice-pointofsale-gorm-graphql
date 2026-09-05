package handler

import (
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/category"
)

type CategoryHandleGrpc interface {
	pb.CategoryServiceServer
}
