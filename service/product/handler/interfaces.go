package handler

import (
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/product"
)

type ProductHandleGrpc interface {
	pb.ProductServiceServer
}
