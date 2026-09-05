package handler

import (
	"context"
	"fmt"

	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/stats"
	"github.com/MamangRust/microservice-point-of-sale-stats-reader/repository"
	"go.uber.org/zap"
)

type CashierStatsHandler struct {
	pb.UnimplementedCashierStatsServiceServer
	repo  repository.Repository
	cache *StatsCache
	log   logger.LoggerInterface
}

func NewCashierStatsHandler(repo repository.Repository, cache *StatsCache, log logger.LoggerInterface) *CashierStatsHandler {
	return &CashierStatsHandler{repo: repo, cache: cache, log: log}
}

func (h *CashierStatsHandler) FindMonthlyOrders(ctx context.Context, req *pb.FindCashierStatsRequest) (*pb.ApiResponseCashierMonthlyOrders, error) {
	key := fmt.Sprintf("stats:reader:cashier:monthly-orders:%d", req.GetCashierId())
	if cached, found := CacheGet[pb.ApiResponseCashierMonthlyOrders](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetCashierMonthlyOrders(ctx, int(req.GetCashierId()))
	if err != nil {
		h.log.Error("FindMonthlyOrders failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseCashierMonthlyOrders{
		Status:  "success",
		Message: "Cashier monthly orders retrieved successfully",
		Data:    mapCashierMonthlyOrders(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *CashierStatsHandler) FindMonthlyTotalSales(ctx context.Context, req *pb.FindYearMonthTotalSales) (*pb.ApiResponseCashierMonthlyTotalSales, error) {
	key := fmt.Sprintf("stats:reader:cashier:monthly-total-sales:%d:%d", req.GetYear(), req.GetMonth())
	if cached, found := CacheGet[pb.ApiResponseCashierMonthlyTotalSales](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetCashierMonthlyTotalSales(ctx, int(req.GetYear()), int(req.GetMonth()))
	if err != nil {
		h.log.Error("FindMonthlyTotalSales failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseCashierMonthlyTotalSales{
		Status:  "success",
		Message: "Monthly total sales retrieved successfully",
		Data:    mapCashierTotalSales(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *CashierStatsHandler) FindYearlyTotalSales(ctx context.Context, req *pb.FindYearTotalSales) (*pb.ApiResponseCashierYearlyTotalSales, error) {
	key := fmt.Sprintf("stats:reader:cashier:yearly-total-sales:%d", req.GetYear())
	if cached, found := CacheGet[pb.ApiResponseCashierYearlyTotalSales](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetCashierYearlyTotalSales(ctx, int(req.GetYear()))
	if err != nil {
		h.log.Error("FindYearlyTotalSales failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseCashierYearlyTotalSales{
		Status:  "success",
		Message: "Yearly total sales retrieved successfully",
		Data:    mapCashierYearTotalSales(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *CashierStatsHandler) FindMonthlyTotalSalesById(ctx context.Context, req *pb.FindYearMonthTotalSalesById) (*pb.ApiResponseCashierMonthlyTotalSales, error) {
	key := fmt.Sprintf("stats:reader:cashier:monthly-total-sales-id:%d:%d:%d", req.GetYear(), req.GetMonth(), req.GetCashierId())
	if cached, found := CacheGet[pb.ApiResponseCashierMonthlyTotalSales](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetCashierMonthlyTotalSalesById(ctx, int(req.GetYear()), int(req.GetMonth()), int(req.GetCashierId()))
	if err != nil {
		h.log.Error("FindMonthlyTotalSalesById failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseCashierMonthlyTotalSales{
		Status:  "success",
		Message: "Monthly total sales by ID retrieved successfully",
		Data:    mapCashierTotalSales(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *CashierStatsHandler) FindYearlyTotalSalesById(ctx context.Context, req *pb.FindYearTotalSalesById) (*pb.ApiResponseCashierYearlyTotalSales, error) {
	key := fmt.Sprintf("stats:reader:cashier:yearly-total-sales-id:%d:%d", req.GetYear(), req.GetCashierId())
	if cached, found := CacheGet[pb.ApiResponseCashierYearlyTotalSales](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetCashierYearlyTotalSalesById(ctx, int(req.GetYear()), int(req.GetCashierId()))
	if err != nil {
		h.log.Error("FindYearlyTotalSalesById failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseCashierYearlyTotalSales{
		Status:  "success",
		Message: "Yearly total sales by ID retrieved successfully",
		Data:    mapCashierYearTotalSales(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *CashierStatsHandler) FindMonthlyTotalSalesByMerchant(ctx context.Context, req *pb.FindYearMonthTotalSalesByMerchant) (*pb.ApiResponseCashierMonthlyTotalSales, error) {
	key := fmt.Sprintf("stats:reader:cashier:monthly-total-sales-merchant:%d:%d:%d", req.GetYear(), req.GetMonth(), req.GetMerchantId())
	if cached, found := CacheGet[pb.ApiResponseCashierMonthlyTotalSales](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetCashierMonthlyTotalSalesByMerchant(ctx, int(req.GetYear()), int(req.GetMonth()), int(req.GetMerchantId()))
	if err != nil {
		h.log.Error("FindMonthlyTotalSalesByMerchant failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseCashierMonthlyTotalSales{
		Status:  "success",
		Message: "Monthly total sales by merchant retrieved successfully",
		Data:    mapCashierTotalSales(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *CashierStatsHandler) FindYearlyTotalSalesByMerchant(ctx context.Context, req *pb.FindYearTotalSalesByMerchant) (*pb.ApiResponseCashierYearlyTotalSales, error) {
	key := fmt.Sprintf("stats:reader:cashier:yearly-total-sales-merchant:%d:%d", req.GetYear(), req.GetMerchantId())
	if cached, found := CacheGet[pb.ApiResponseCashierYearlyTotalSales](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetCashierYearlyTotalSalesByMerchant(ctx, int(req.GetYear()), int(req.GetMerchantId()))
	if err != nil {
		h.log.Error("FindYearlyTotalSalesByMerchant failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseCashierYearlyTotalSales{
		Status:  "success",
		Message: "Yearly total sales by merchant retrieved successfully",
		Data:    mapCashierYearTotalSales(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *CashierStatsHandler) FindMonthSales(ctx context.Context, req *pb.FindYearCashier) (*pb.ApiResponseCashierMonthSales, error) {
	key := fmt.Sprintf("stats:reader:cashier:month-sales:%d", req.GetYear())
	if cached, found := CacheGet[pb.ApiResponseCashierMonthSales](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetCashierSalesMonthly(ctx, int(req.GetYear()))
	if err != nil {
		h.log.Error("FindMonthSales failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseCashierMonthSales{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapCashierSalesMonthly(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *CashierStatsHandler) FindYearSales(ctx context.Context, req *pb.FindYearCashier) (*pb.ApiResponseCashierYearSales, error) {
	key := fmt.Sprintf("stats:reader:cashier:year-sales:%d", req.GetYear())
	if cached, found := CacheGet[pb.ApiResponseCashierYearSales](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetCashierSalesYearly(ctx, int(req.GetYear()))
	if err != nil {
		h.log.Error("FindYearSales failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseCashierYearSales{
		Status:  "success",
		Message: "Yearly sales retrieved successfully",
		Data:    mapCashierSalesYearly(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *CashierStatsHandler) FindMonthSalesByMerchant(ctx context.Context, req *pb.FindYearCashierByMerchant) (*pb.ApiResponseCashierMonthSales, error) {
	key := fmt.Sprintf("stats:reader:cashier:month-sales-merchant:%d:%d", req.GetYear(), req.GetMerchantId())
	if cached, found := CacheGet[pb.ApiResponseCashierMonthSales](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetCashierSalesMonthlyByMerchant(ctx, int(req.GetYear()), int(req.GetMerchantId()))
	if err != nil {
		h.log.Error("FindMonthSalesByMerchant failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseCashierMonthSales{
		Status:  "success",
		Message: "Monthly sales by merchant retrieved successfully",
		Data:    mapCashierSalesMonthly(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *CashierStatsHandler) FindYearSalesByMerchant(ctx context.Context, req *pb.FindYearCashierByMerchant) (*pb.ApiResponseCashierYearSales, error) {
	key := fmt.Sprintf("stats:reader:cashier:year-sales-merchant:%d:%d", req.GetYear(), req.GetMerchantId())
	if cached, found := CacheGet[pb.ApiResponseCashierYearSales](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetCashierSalesYearlyByMerchant(ctx, int(req.GetYear()), int(req.GetMerchantId()))
	if err != nil {
		h.log.Error("FindYearSalesByMerchant failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseCashierYearSales{
		Status:  "success",
		Message: "Yearly sales by merchant retrieved successfully",
		Data:    mapCashierSalesYearly(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *CashierStatsHandler) FindMonthSalesById(ctx context.Context, req *pb.FindYearCashierById) (*pb.ApiResponseCashierMonthSales, error) {
	key := fmt.Sprintf("stats:reader:cashier:month-sales-id:%d:%d", req.GetYear(), req.GetCashierId())
	if cached, found := CacheGet[pb.ApiResponseCashierMonthSales](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetCashierSalesMonthlyById(ctx, int(req.GetYear()), int(req.GetCashierId()))
	if err != nil {
		h.log.Error("FindMonthSalesById failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseCashierMonthSales{
		Status:  "success",
		Message: "Monthly sales by ID retrieved successfully",
		Data:    mapCashierSalesMonthly(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *CashierStatsHandler) FindYearSalesById(ctx context.Context, req *pb.FindYearCashierById) (*pb.ApiResponseCashierYearSales, error) {
	key := fmt.Sprintf("stats:reader:cashier:year-sales-id:%d:%d", req.GetYear(), req.GetCashierId())
	if cached, found := CacheGet[pb.ApiResponseCashierYearSales](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetCashierSalesYearlyById(ctx, int(req.GetYear()), int(req.GetCashierId()))
	if err != nil {
		h.log.Error("FindYearSalesById failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseCashierYearSales{
		Status:  "success",
		Message: "Yearly sales by ID retrieved successfully",
		Data:    mapCashierSalesYearly(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

// ── Mappers ─────────────────────────────────────────────────────────────

func mapCashierMonthlyOrders(data []repository.CashierMonthlyOrders) []*pb.CashierMonthlyOrdersResponse {
	var out []*pb.CashierMonthlyOrdersResponse
	for _, d := range data {
		out = append(out, &pb.CashierMonthlyOrdersResponse{
			Month:       d.Month,
			CashierId:   int32(d.CashierID),
			OrderCount:  int64(d.OrderCount),
			TotalAmount: d.TotalAmount,
		})
	}
	return out
}

func mapCashierTotalSales(data []repository.CashierTotalSales) []*pb.CashierResponseMonthTotalSales {
	var out []*pb.CashierResponseMonthTotalSales
	for _, d := range data {
		out = append(out, &pb.CashierResponseMonthTotalSales{
			Year:       d.Year,
			Month:      d.Month,
			TotalSales: int32(d.TotalSales),
		})
	}
	return out
}

func mapCashierYearTotalSales(data []repository.CashierYearTotalSales) []*pb.CashierResponseYearTotalSales {
	var out []*pb.CashierResponseYearTotalSales
	for _, d := range data {
		out = append(out, &pb.CashierResponseYearTotalSales{
			Year:       d.Year,
			TotalSales: int32(d.TotalSales),
		})
	}
	return out
}

func mapCashierSalesMonthly(data []repository.CashierSalesMonthly) []*pb.CashierResponseMonthSales {
	var out []*pb.CashierResponseMonthSales
	for _, d := range data {
		out = append(out, &pb.CashierResponseMonthSales{
			Month:      d.Month,
			CashierId:  int32(d.CashierID),
			OrderCount: int32(d.OrderCount),
			TotalSales: int32(d.TotalSales),
		})
	}
	return out
}

func mapCashierSalesYearly(data []repository.CashierSalesYearly) []*pb.CashierResponseYearSales {
	var out []*pb.CashierResponseYearSales
	for _, d := range data {
		out = append(out, &pb.CashierResponseYearSales{
			Year:       d.Year,
			CashierId:  int32(d.CashierID),
			OrderCount: int32(d.OrderCount),
			TotalSales: int32(d.TotalSales),
		})
	}
	return out
}
