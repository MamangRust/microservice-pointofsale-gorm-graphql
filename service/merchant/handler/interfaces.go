package handler

import (
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/merchant"
)

type MerchantDocumentHandleGrpc interface {
	pb.MerchantDocumentServiceServer
}

type MerchantHandleGrpc interface {
	pb.MerchantServiceServer
}
