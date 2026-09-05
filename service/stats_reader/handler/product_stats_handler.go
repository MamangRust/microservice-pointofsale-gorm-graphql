package handler

import (
	"context"
	"fmt"

	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/stats"
	"github.com/MamangRust/microservice-point-of-sale-stats-reader/repository"
	"go.uber.org/zap"
)

// ProductStatsHandler serves per-product sold aggregates from ClickHouse
// (order_item_daily), cached in Redis for 5 minutes.
type ProductStatsHandler struct {
	pb.UnimplementedProductStatsServiceServer
	repo  repository.Repository
	cache *StatsCache
	log   logger.LoggerInterface
}

func NewProductStatsHandler(repo repository.Repository, cache *StatsCache, log logger.LoggerInterface) *ProductStatsHandler {
	return &ProductStatsHandler{repo: repo, cache: cache, log: log}
}

func (h *ProductStatsHandler) FindMonthlySold(ctx context.Context, req *pb.FindYearMonthStatsRequest) (*pb.ApiResponseProductMonthlySold, error) {
	key := fmt.Sprintf("stats:reader:product:monthly-sold:%d:%d", req.GetYear(), req.GetMonth())
	if cached, found := CacheGet[pb.ApiResponseProductMonthlySold](ctx, h.cache, key); found {
		return cached, nil
	}

	data, err := h.repo.GetProductMonthlySold(ctx, int(req.GetYear()), int(req.GetMonth()))
	if err != nil {
		h.log.Error("FindMonthlySold failed", zap.Error(err))
		return nil, err
	}

	resp := &pb.ApiResponseProductMonthlySold{
		Status:  "success",
		Message: "Monthly product sold retrieved successfully",
		Data:    mapProductMonthlySold(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func mapProductMonthlySold(data []repository.ProductMonthlySold) []*pb.ProductMonthlySoldResponse {
	var out []*pb.ProductMonthlySoldResponse
	for _, d := range data {
		out = append(out, &pb.ProductMonthlySoldResponse{
			Month:     d.Month,
			ProductId: int32(d.ProductID),
			Quantity:  int64(d.Quantity),
			Subtotal:  d.Subtotal,
		})
	}
	return out
}
