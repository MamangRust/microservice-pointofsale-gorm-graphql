package handler

import (
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/order_item"
)

type OrderItemHandlerGrpc interface {
	pb.OrderItemServiceServer
}
