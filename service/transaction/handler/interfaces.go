package handler

import (
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/transaction"
)

type TransactionHandleGrpc interface {
	pb.TransactionServiceServer
}
