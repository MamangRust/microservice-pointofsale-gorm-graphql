package handler

import (
	"context"
	"fmt"

	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/stats"
	"github.com/MamangRust/microservice-point-of-sale-stats-reader/repository"
	"go.uber.org/zap"
)

type CategoryStatsHandler struct {
	pb.UnimplementedCategoryStatsServiceServer
	repo  repository.Repository
	cache *StatsCache
	log   logger.LoggerInterface
}

func NewCategoryStatsHandler(repo repository.Repository, cache *StatsCache, log logger.LoggerInterface) *CategoryStatsHandler {
	return &CategoryStatsHandler{repo: repo, cache: cache, log: log}
}

func (h *CategoryStatsHandler) FindMonthlySold(ctx context.Context, req *pb.FindYearMonthStatsRequest) (*pb.ApiResponseCategoryMonthlySold, error) {
	key := fmt.Sprintf("stats:reader:category:monthly-sold:%d:%d", req.GetYear(), req.GetMonth())
	if cached, found := CacheGet[pb.ApiResponseCategoryMonthlySold](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetCategoryMonthlySold(ctx, int(req.GetYear()), int(req.GetMonth()))
	if err != nil {
		h.log.Error("FindMonthlySold failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseCategoryMonthlySold{
		Status:  "success",
		Message: "Monthly category sold retrieved successfully",
		Data:    mapCategoryMonthlySold(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *CategoryStatsHandler) FindMonthPrice(ctx context.Context, req *pb.FindYearCategory) (*pb.ApiResponseCategoryMonthPrice, error) {
	key := fmt.Sprintf("stats:reader:category:month-price:%d", req.GetYear())
	if cached, found := CacheGet[pb.ApiResponseCategoryMonthPrice](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetCategoryMonthPrice(ctx, int(req.GetYear()))
	if err != nil {
		h.log.Error("FindMonthPrice failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseCategoryMonthPrice{
		Status:  "success",
		Message: "Monthly category price retrieved successfully",
		Data:    mapCategoryMonthPrice(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *CategoryStatsHandler) FindYearPrice(ctx context.Context, req *pb.FindYearCategory) (*pb.ApiResponseCategoryYearPrice, error) {
	key := fmt.Sprintf("stats:reader:category:year-price:%d", req.GetYear())
	if cached, found := CacheGet[pb.ApiResponseCategoryYearPrice](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetCategoryYearPrice(ctx, int(req.GetYear()))
	if err != nil {
		h.log.Error("FindYearPrice failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseCategoryYearPrice{
		Status:  "success",
		Message: "Yearly category price retrieved successfully",
		Data:    mapCategoryYearPrice(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *CategoryStatsHandler) FindMonthPriceByMerchant(ctx context.Context, req *pb.FindYearCategoryByMerchant) (*pb.ApiResponseCategoryMonthPrice, error) {
	key := fmt.Sprintf("stats:reader:category:month-price-merchant:%d:%d", req.GetYear(), req.GetMerchantId())
	if cached, found := CacheGet[pb.ApiResponseCategoryMonthPrice](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetCategoryMonthPriceByMerchant(ctx, int(req.GetYear()), int(req.GetMerchantId()))
	if err != nil {
		h.log.Error("FindMonthPriceByMerchant failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseCategoryMonthPrice{
		Status:  "success",
		Message: "Monthly category price by merchant retrieved successfully",
		Data:    mapCategoryMonthPrice(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *CategoryStatsHandler) FindYearPriceByMerchant(ctx context.Context, req *pb.FindYearCategoryByMerchant) (*pb.ApiResponseCategoryYearPrice, error) {
	key := fmt.Sprintf("stats:reader:category:year-price-merchant:%d:%d", req.GetYear(), req.GetMerchantId())
	if cached, found := CacheGet[pb.ApiResponseCategoryYearPrice](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetCategoryYearPriceByMerchant(ctx, int(req.GetYear()), int(req.GetMerchantId()))
	if err != nil {
		h.log.Error("FindYearPriceByMerchant failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseCategoryYearPrice{
		Status:  "success",
		Message: "Yearly category price by merchant retrieved successfully",
		Data:    mapCategoryYearPrice(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *CategoryStatsHandler) FindMonthPriceById(ctx context.Context, req *pb.FindYearCategoryById) (*pb.ApiResponseCategoryMonthPrice, error) {
	key := fmt.Sprintf("stats:reader:category:month-price-id:%d:%d", req.GetYear(), req.GetCategoryId())
	if cached, found := CacheGet[pb.ApiResponseCategoryMonthPrice](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetCategoryMonthPriceById(ctx, int(req.GetYear()), int(req.GetCategoryId()))
	if err != nil {
		h.log.Error("FindMonthPriceById failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseCategoryMonthPrice{
		Status:  "success",
		Message: "Monthly category price by ID retrieved successfully",
		Data:    mapCategoryMonthPrice(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *CategoryStatsHandler) FindYearPriceById(ctx context.Context, req *pb.FindYearCategoryById) (*pb.ApiResponseCategoryYearPrice, error) {
	key := fmt.Sprintf("stats:reader:category:year-price-id:%d:%d", req.GetYear(), req.GetCategoryId())
	if cached, found := CacheGet[pb.ApiResponseCategoryYearPrice](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetCategoryYearPriceById(ctx, int(req.GetYear()), int(req.GetCategoryId()))
	if err != nil {
		h.log.Error("FindYearPriceById failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseCategoryYearPrice{
		Status:  "success",
		Message: "Yearly category price by ID retrieved successfully",
		Data:    mapCategoryYearPrice(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *CategoryStatsHandler) FindMonthlyTotalPrices(ctx context.Context, req *pb.FindYearMonthTotalPrices) (*pb.ApiResponseCategoryMonthlyTotalPrice, error) {
	key := fmt.Sprintf("stats:reader:category:monthly-total-prices:%d:%d", req.GetYear(), req.GetMonth())
	if cached, found := CacheGet[pb.ApiResponseCategoryMonthlyTotalPrice](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetCategoryMonthlyTotalPrices(ctx, int(req.GetYear()), int(req.GetMonth()))
	if err != nil {
		h.log.Error("FindMonthlyTotalPrices failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseCategoryMonthlyTotalPrice{
		Status:  "success",
		Message: "Monthly total prices retrieved successfully",
		Data:    mapCategoryTotalPrice(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *CategoryStatsHandler) FindYearlyTotalPrices(ctx context.Context, req *pb.FindYearTotalPrices) (*pb.ApiResponseCategoryYearlyTotalPrice, error) {
	key := fmt.Sprintf("stats:reader:category:yearly-total-prices:%d", req.GetYear())
	if cached, found := CacheGet[pb.ApiResponseCategoryYearlyTotalPrice](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetCategoryYearlyTotalPrices(ctx, int(req.GetYear()))
	if err != nil {
		h.log.Error("FindYearlyTotalPrices failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseCategoryYearlyTotalPrice{
		Status:  "success",
		Message: "Yearly total prices retrieved successfully",
		Data:    mapCategoryYearTotalPrice(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *CategoryStatsHandler) FindMonthlyTotalPricesById(ctx context.Context, req *pb.FindYearMonthTotalPriceById) (*pb.ApiResponseCategoryMonthlyTotalPrice, error) {
	key := fmt.Sprintf("stats:reader:category:monthly-total-prices-id:%d:%d:%d", req.GetYear(), req.GetMonth(), req.GetCategoryId())
	if cached, found := CacheGet[pb.ApiResponseCategoryMonthlyTotalPrice](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetCategoryMonthlyTotalPricesById(ctx, int(req.GetYear()), int(req.GetMonth()), int(req.GetCategoryId()))
	if err != nil {
		h.log.Error("FindMonthlyTotalPricesById failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseCategoryMonthlyTotalPrice{
		Status:  "success",
		Message: "Monthly total prices by ID retrieved successfully",
		Data:    mapCategoryTotalPrice(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *CategoryStatsHandler) FindYearlyTotalPricesById(ctx context.Context, req *pb.FindYearTotalPriceById) (*pb.ApiResponseCategoryYearlyTotalPrice, error) {
	key := fmt.Sprintf("stats:reader:category:yearly-total-prices-id:%d:%d", req.GetYear(), req.GetCategoryId())
	if cached, found := CacheGet[pb.ApiResponseCategoryYearlyTotalPrice](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetCategoryYearlyTotalPricesById(ctx, int(req.GetYear()), int(req.GetCategoryId()))
	if err != nil {
		h.log.Error("FindYearlyTotalPricesById failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseCategoryYearlyTotalPrice{
		Status:  "success",
		Message: "Yearly total prices by ID retrieved successfully",
		Data:    mapCategoryYearTotalPrice(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *CategoryStatsHandler) FindMonthlyTotalPricesByMerchant(ctx context.Context, req *pb.FindYearMonthTotalPriceByMerchant) (*pb.ApiResponseCategoryMonthlyTotalPrice, error) {
	key := fmt.Sprintf("stats:reader:category:monthly-total-prices-merchant:%d:%d:%d", req.GetYear(), req.GetMonth(), req.GetMerchantId())
	if cached, found := CacheGet[pb.ApiResponseCategoryMonthlyTotalPrice](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetCategoryMonthlyTotalPricesByMerchant(ctx, int(req.GetYear()), int(req.GetMonth()), int(req.GetMerchantId()))
	if err != nil {
		h.log.Error("FindMonthlyTotalPricesByMerchant failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseCategoryMonthlyTotalPrice{
		Status:  "success",
		Message: "Monthly total prices by merchant retrieved successfully",
		Data:    mapCategoryTotalPrice(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *CategoryStatsHandler) FindYearlyTotalPricesByMerchant(ctx context.Context, req *pb.FindYearTotalPriceByMerchant) (*pb.ApiResponseCategoryYearlyTotalPrice, error) {
	key := fmt.Sprintf("stats:reader:category:yearly-total-prices-merchant:%d:%d", req.GetYear(), req.GetMerchantId())
	if cached, found := CacheGet[pb.ApiResponseCategoryYearlyTotalPrice](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetCategoryYearlyTotalPricesByMerchant(ctx, int(req.GetYear()), int(req.GetMerchantId()))
	if err != nil {
		h.log.Error("FindYearlyTotalPricesByMerchant failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseCategoryYearlyTotalPrice{
		Status:  "success",
		Message: "Yearly total prices by merchant retrieved successfully",
		Data:    mapCategoryYearTotalPrice(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

// ── Mappers ─────────────────────────────────────────────────────────────

func mapCategoryMonthlySold(data []repository.CategoryMonthlySold) []*pb.CategoryMonthlySoldResponse {
	var out []*pb.CategoryMonthlySoldResponse
	for _, d := range data {
		out = append(out, &pb.CategoryMonthlySoldResponse{
			Month:      d.Month,
			CategoryId: int32(d.CategoryID),
			Quantity:   int64(d.Quantity),
			Subtotal:   d.Subtotal,
		})
	}
	return out
}

func mapCategoryMonthPrice(data []repository.CategoryMonthPrice) []*pb.CategoryMonthPriceResponse {
	var out []*pb.CategoryMonthPriceResponse
	for _, d := range data {
		out = append(out, &pb.CategoryMonthPriceResponse{
			Month:      d.Month,
			CategoryId: int32(d.CategoryID),
			OrderCount: int32(d.OrderCount),
			ItemsSold:  int32(d.ItemsSold),
			TotalRevenue: int32(d.TotalRev),
		})
	}
	return out
}

func mapCategoryYearPrice(data []repository.CategoryYearPrice) []*pb.CategoryYearPriceResponse {
	var out []*pb.CategoryYearPriceResponse
	for _, d := range data {
		out = append(out, &pb.CategoryYearPriceResponse{
			Year:               d.Year,
			CategoryId:         int32(d.CategoryID),
			OrderCount:         int32(d.OrderCount),
			ItemsSold:          int32(d.ItemsSold),
			TotalRevenue:       int32(d.TotalRev),
			UniqueProductsSold: int32(d.UniqueProd),
		})
	}
	return out
}

func mapCategoryTotalPrice(data []repository.CategoryTotalPrice) []*pb.CategoriesMonthlyTotalPriceResponse {
	var out []*pb.CategoriesMonthlyTotalPriceResponse
	for _, d := range data {
		out = append(out, &pb.CategoriesMonthlyTotalPriceResponse{
			Year:        d.Year,
			Month:       d.Month,
			TotalRevenue: int32(d.TotalRev),
		})
	}
	return out
}

func mapCategoryYearTotalPrice(data []repository.CategoryYearTotalPrice) []*pb.CategoriesYearlyTotalPriceResponse {
	var out []*pb.CategoriesYearlyTotalPriceResponse
	for _, d := range data {
		out = append(out, &pb.CategoriesYearlyTotalPriceResponse{
			Year:        d.Year,
			TotalRevenue: int32(d.TotalRev),
		})
	}
	return out
}
