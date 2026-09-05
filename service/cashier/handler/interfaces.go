package handler

import pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/cashier"

type CashierHandleGrpc interface {
	pb.CashierServiceServer
}
