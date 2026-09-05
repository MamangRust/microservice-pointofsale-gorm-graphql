package handler

import (
	"context"
	"fmt"

	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/stats"
	"github.com/MamangRust/microservice-point-of-sale-stats-reader/repository"
	"go.uber.org/zap"
)

type OrderStatsHandler struct {
	pb.UnimplementedOrderStatsServiceServer
	repo  repository.Repository
	cache *StatsCache
	log   logger.LoggerInterface
}

func NewOrderStatsHandler(repo repository.Repository, cache *StatsCache, log logger.LoggerInterface) *OrderStatsHandler {
	return &OrderStatsHandler{repo: repo, cache: cache, log: log}
}

func (h *OrderStatsHandler) FindMonthlyRevenue(ctx context.Context, req *pb.FindYearOrder) (*pb.ApiResponseOrderMonthly, error) {
	key := fmt.Sprintf("stats:reader:order:monthly-revenue:%d", req.GetYear())
	if cached, found := CacheGet[pb.ApiResponseOrderMonthly](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetOrderMonthlyDetail(ctx, int(req.GetYear()))
	if err != nil {
		h.log.Error("FindMonthlyRevenue failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseOrderMonthly{
		Status:  "success",
		Message: "Monthly revenue retrieved successfully",
		Data:    mapOrderMonthlyDetail(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *OrderStatsHandler) FindYearlyRevenue(ctx context.Context, req *pb.FindYearOrder) (*pb.ApiResponseOrderYearly, error) {
	key := fmt.Sprintf("stats:reader:order:yearly-revenue:%d", req.GetYear())
	if cached, found := CacheGet[pb.ApiResponseOrderYearly](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetOrderYearlyDetail(ctx, int(req.GetYear()))
	if err != nil {
		h.log.Error("FindYearlyRevenue failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseOrderYearly{
		Status:  "success",
		Message: "Yearly revenue retrieved successfully",
		Data:    mapOrderYearlyDetail(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *OrderStatsHandler) FindMonthlyRevenueByMerchant(ctx context.Context, req *pb.FindYearOrderByMerchant) (*pb.ApiResponseOrderMonthly, error) {
	key := fmt.Sprintf("stats:reader:order:monthly-revenue-merchant:%d:%d", req.GetYear(), req.GetMerchantId())
	if cached, found := CacheGet[pb.ApiResponseOrderMonthly](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetOrderMonthlyDetailByMerchant(ctx, int(req.GetYear()), int(req.GetMerchantId()))
	if err != nil {
		h.log.Error("FindMonthlyRevenueByMerchant failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseOrderMonthly{
		Status:  "success",
		Message: "Monthly revenue by merchant retrieved successfully",
		Data:    mapOrderMonthlyDetail(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *OrderStatsHandler) FindYearlyRevenueByMerchant(ctx context.Context, req *pb.FindYearOrderByMerchant) (*pb.ApiResponseOrderYearly, error) {
	key := fmt.Sprintf("stats:reader:order:yearly-revenue-merchant:%d:%d", req.GetYear(), req.GetMerchantId())
	if cached, found := CacheGet[pb.ApiResponseOrderYearly](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetOrderYearlyDetailByMerchant(ctx, int(req.GetYear()), int(req.GetMerchantId()))
	if err != nil {
		h.log.Error("FindYearlyRevenueByMerchant failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseOrderYearly{
		Status:  "success",
		Message: "Yearly revenue by merchant retrieved successfully",
		Data:    mapOrderYearlyDetail(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *OrderStatsHandler) FindMonthlyTotalRevenue(ctx context.Context, req *pb.FindYearMonthTotalRevenue) (*pb.ApiResponseOrderMonthlyTotalRevenue, error) {
	key := fmt.Sprintf("stats:reader:order:monthly-total-revenue:%d:%d", req.GetYear(), req.GetMonth())
	if cached, found := CacheGet[pb.ApiResponseOrderMonthlyTotalRevenue](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetOrderMonthlyTotalRevenue(ctx, int(req.GetYear()), int(req.GetMonth()))
	if err != nil {
		h.log.Error("FindMonthlyTotalRevenue failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseOrderMonthlyTotalRevenue{
		Status:  "success",
		Message: "Monthly total revenue retrieved successfully",
		Data:    mapOrderTotalRevenue(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *OrderStatsHandler) FindYearlyTotalRevenue(ctx context.Context, req *pb.FindYearTotalRevenue) (*pb.ApiResponseOrderYearlyTotalRevenue, error) {
	key := fmt.Sprintf("stats:reader:order:yearly-total-revenue:%d", req.GetYear())
	if cached, found := CacheGet[pb.ApiResponseOrderYearlyTotalRevenue](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetOrderYearlyTotalRevenue(ctx, int(req.GetYear()))
	if err != nil {
		h.log.Error("FindYearlyTotalRevenue failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseOrderYearlyTotalRevenue{
		Status:  "success",
		Message: "Yearly total revenue retrieved successfully",
		Data:    mapOrderYearTotalRevenue(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *OrderStatsHandler) FindMonthlyTotalRevenueById(ctx context.Context, req *pb.FindYearMonthTotalRevenueById) (*pb.ApiResponseOrderMonthlyTotalRevenue, error) {
	key := fmt.Sprintf("stats:reader:order:monthly-total-revenue-id:%d:%d:%d", req.GetYear(), req.GetMonth(), req.GetOrderId())
	if cached, found := CacheGet[pb.ApiResponseOrderMonthlyTotalRevenue](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetOrderMonthlyTotalRevenueById(ctx, int(req.GetYear()), int(req.GetMonth()), int(req.GetOrderId()))
	if err != nil {
		h.log.Error("FindMonthlyTotalRevenueById failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseOrderMonthlyTotalRevenue{
		Status:  "success",
		Message: "Monthly total revenue by ID retrieved successfully",
		Data:    mapOrderTotalRevenue(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *OrderStatsHandler) FindYearlyTotalRevenueById(ctx context.Context, req *pb.FindYearTotalRevenueById) (*pb.ApiResponseOrderYearlyTotalRevenue, error) {
	key := fmt.Sprintf("stats:reader:order:yearly-total-revenue-id:%d:%d", req.GetYear(), req.GetOrderId())
	if cached, found := CacheGet[pb.ApiResponseOrderYearlyTotalRevenue](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetOrderYearlyTotalRevenueById(ctx, int(req.GetYear()), int(req.GetOrderId()))
	if err != nil {
		h.log.Error("FindYearlyTotalRevenueById failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseOrderYearlyTotalRevenue{
		Status:  "success",
		Message: "Yearly total revenue by ID retrieved successfully",
		Data:    mapOrderYearTotalRevenue(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *OrderStatsHandler) FindMonthlyTotalRevenueByMerchant(ctx context.Context, req *pb.FindYearMonthTotalRevenueByMerchant) (*pb.ApiResponseOrderMonthlyTotalRevenue, error) {
	key := fmt.Sprintf("stats:reader:order:monthly-total-revenue-merchant:%d:%d:%d", req.GetYear(), req.GetMonth(), req.GetMerchantId())
	if cached, found := CacheGet[pb.ApiResponseOrderMonthlyTotalRevenue](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetOrderMonthlyTotalRevenueByMerchant(ctx, int(req.GetYear()), int(req.GetMonth()), int(req.GetMerchantId()))
	if err != nil {
		h.log.Error("FindMonthlyTotalRevenueByMerchant failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseOrderMonthlyTotalRevenue{
		Status:  "success",
		Message: "Monthly total revenue by merchant retrieved successfully",
		Data:    mapOrderTotalRevenue(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *OrderStatsHandler) FindYearlyTotalRevenueByMerchant(ctx context.Context, req *pb.FindYearTotalRevenueByMerchant) (*pb.ApiResponseOrderYearlyTotalRevenue, error) {
	key := fmt.Sprintf("stats:reader:order:yearly-total-revenue-merchant:%d:%d", req.GetYear(), req.GetMerchantId())
	if cached, found := CacheGet[pb.ApiResponseOrderYearlyTotalRevenue](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetOrderYearlyTotalRevenueByMerchant(ctx, int(req.GetYear()), int(req.GetMerchantId()))
	if err != nil {
		h.log.Error("FindYearlyTotalRevenueByMerchant failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseOrderYearlyTotalRevenue{
		Status:  "success",
		Message: "Yearly total revenue by merchant retrieved successfully",
		Data:    mapOrderYearTotalRevenue(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *OrderStatsHandler) FindMonthlyTotalRevenueFromCH(ctx context.Context, req *pb.FindYearMonthStatsRequest) (*pb.ApiResponseOrderMonthlyRevenue, error) {
	key := fmt.Sprintf("stats:reader:order:monthly-revenue-ch:%d:%d", req.GetYear(), req.GetMonth())
	if cached, found := CacheGet[pb.ApiResponseOrderMonthlyRevenue](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetMonthlyTotalRevenue(ctx, int(req.GetYear()), int(req.GetMonth()))
	if err != nil {
		h.log.Error("FindMonthlyTotalRevenueFromCH failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseOrderMonthlyRevenue{
		Status:  "success",
		Message: "Monthly revenue retrieved successfully",
		Data:    mapMonthlyRevenue(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *OrderStatsHandler) FindYearlyTotalRevenueFromCH(ctx context.Context, req *pb.FindYearStatsRequest) (*pb.ApiResponseOrderYearlyRevenue, error) {
	key := fmt.Sprintf("stats:reader:order:yearly-revenue-ch:%d", req.GetYear())
	if cached, found := CacheGet[pb.ApiResponseOrderYearlyRevenue](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetYearlyTotalRevenue(ctx, int(req.GetYear()))
	if err != nil {
		h.log.Error("FindYearlyTotalRevenueFromCH failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseOrderYearlyRevenue{
		Status:  "success",
		Message: "Yearly revenue retrieved successfully",
		Data:    mapYearlyRevenue(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *OrderStatsHandler) FindCashierMonthlyRevenue(ctx context.Context, req *pb.FindCashierStatsRequest) (*pb.ApiResponseCashierMonthlyRevenue, error) {
	key := fmt.Sprintf("stats:reader:order:cashier-monthly-revenue:%d", req.GetCashierId())
	if cached, found := CacheGet[pb.ApiResponseCashierMonthlyRevenue](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetCashierMonthlyRevenue(ctx, int(req.GetCashierId()))
	if err != nil {
		h.log.Error("FindCashierMonthlyRevenue failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseCashierMonthlyRevenue{
		Status:  "success",
		Message: "Cashier monthly revenue retrieved successfully",
		Data:    mapCashierMonthlyRevenue(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

// ── Mappers ─────────────────────────────────────────────────────────────

func mapOrderMonthlyDetail(data []repository.OrderMonthlyDetail) []*pb.OrderMonthlyResponse {
	var out []*pb.OrderMonthlyResponse
	for _, d := range data {
		out = append(out, &pb.OrderMonthlyResponse{
			Month:          d.Month,
			OrderCount:     int32(d.OrderCount),
			TotalRevenue:   int32(d.TotalRev),
			TotalItemsSold: int32(d.ItemsSold),
		})
	}
	return out
}

func mapOrderYearlyDetail(data []repository.OrderYearlyDetail) []*pb.OrderYearlyResponse {
	var out []*pb.OrderYearlyResponse
	for _, d := range data {
		out = append(out, &pb.OrderYearlyResponse{
			Year:              d.Year,
			OrderCount:        int32(d.OrderCount),
			TotalRevenue:      int32(d.TotalRev),
			TotalItemsSold:    int32(d.ItemsSold),
			ActiveCashiers:    int32(d.ActiveCashiers),
			UniqueProductsSold: int32(d.UniqueProds),
		})
	}
	return out
}

func mapOrderTotalRevenue(data []repository.OrderTotalRevenue) []*pb.OrderMonthlyTotalRevenueResponse {
	var out []*pb.OrderMonthlyTotalRevenueResponse
	for _, d := range data {
		out = append(out, &pb.OrderMonthlyTotalRevenueResponse{
			Year:          d.Year,
			Month:         d.Month,
			OrderCount:    int32(d.OrderCount),
			TotalRevenue:  int32(d.TotalRev),
			TotalItemsSold: int32(d.ItemsSold),
		})
	}
	return out
}

func mapOrderYearTotalRevenue(data []repository.OrderYearTotalRevenue) []*pb.OrderYearlyTotalRevenueResponse {
	var out []*pb.OrderYearlyTotalRevenueResponse
	for _, d := range data {
		out = append(out, &pb.OrderYearlyTotalRevenueResponse{
			Year:              d.Year,
			OrderCount:        int32(d.OrderCount),
			TotalRevenue:      int32(d.TotalRev),
			TotalItemsSold:    int32(d.ItemsSold),
			ActiveCashiers:    int32(d.ActiveCashiers),
			UniqueProductsSold: int32(d.UniqueProds),
		})
	}
	return out
}

func mapMonthlyRevenue(data []repository.MonthlyRevenue) []*pb.OrderMonthlyRevenueResponse {
	var out []*pb.OrderMonthlyRevenueResponse
	for _, d := range data {
		out = append(out, &pb.OrderMonthlyRevenueResponse{
			Year:         d.Year,
			Month:        d.Month,
			TotalRevenue: d.TotalRevenue,
			OrderCount:   int32(d.OrderCount),
		})
	}
	return out
}

func mapYearlyRevenue(data []repository.YearlyRevenue) []*pb.OrderYearlyRevenueResponse {
	var out []*pb.OrderYearlyRevenueResponse
	for _, d := range data {
		out = append(out, &pb.OrderYearlyRevenueResponse{
			Year:         d.Year,
			TotalRevenue: d.TotalRevenue,
			OrderCount:   int32(d.OrderCount),
		})
	}
	return out
}

func mapCashierMonthlyRevenue(data []repository.CashierMonthlyRevenue) []*pb.CashierMonthlyRevenueResponse {
	var out []*pb.CashierMonthlyRevenueResponse
	for _, d := range data {
		out = append(out, &pb.CashierMonthlyRevenueResponse{
			Year:         d.Year,
			Month:        d.Month,
			CashierId:    int32(d.CashierID),
			TotalRevenue: d.TotalRevenue,
			OrderCount:   int32(d.OrderCount),
		})
	}
	return out
}
