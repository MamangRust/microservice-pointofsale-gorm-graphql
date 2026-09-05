package handler

import (
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/order"
)

type OrderHandleGrpc interface {
	pb.OrderServiceServer
}
