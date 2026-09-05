package handler

import (
	"context"
	"fmt"

	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/stats"
	"github.com/MamangRust/microservice-point-of-sale-stats-reader/repository"
	"go.uber.org/zap"
)

type TransactionStatsHandler struct {
	pb.UnimplementedTransactionStatsServiceServer
	repo  repository.Repository
	cache *StatsCache
	log   logger.LoggerInterface
}

func NewTransactionStatsHandler(repo repository.Repository, cache *StatsCache, log logger.LoggerInterface) *TransactionStatsHandler {
	return &TransactionStatsHandler{repo: repo, cache: cache, log: log}
}

func (h *TransactionStatsHandler) FindMonthlySuccess(ctx context.Context, req *pb.FindYearMonthStatsRequest) (*pb.ApiResponseTransactionMonthlySuccess, error) {
	key := fmt.Sprintf("stats:reader:transaction:monthly-success:%d:%d", int(req.GetYear()), int(req.GetMonth()))
	if cached, found := CacheGet[pb.ApiResponseTransactionMonthlySuccess](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetTransactionMonthlySuccess(ctx, int(req.GetYear()), int(req.GetMonth()))
	if err != nil {
		h.log.Error("FindMonthlySuccess failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseTransactionMonthlySuccess{
		Status:  "success",
		Message: "Monthly transaction success retrieved successfully",
		Data:    mapTransactionMonthlySuccess(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *TransactionStatsHandler) FindMonthStatusSuccess(ctx context.Context, req *pb.FindMonthlyTransactionStatus) (*pb.ApiResponseTransactionMonthAmountSuccess, error) {
	return h.findStatusMonthly(ctx, int(req.GetYear()), int(req.GetMonth()), "success")
}

func (h *TransactionStatsHandler) FindYearStatusSuccess(ctx context.Context, req *pb.FindYearlyTransactionStatus) (*pb.ApiResponseTransactionYearAmountSuccess, error) {
	return h.findStatusYearly(ctx, int(req.GetYear()), "success")
}

func (h *TransactionStatsHandler) FindMonthStatusFailed(ctx context.Context, req *pb.FindMonthlyTransactionStatus) (*pb.ApiResponseTransactionMonthAmountFailed, error) {
	return h.findStatusMonthlyFailed(ctx, int(req.GetYear()), int(req.GetMonth()), "failed")
}

func (h *TransactionStatsHandler) FindYearStatusFailed(ctx context.Context, req *pb.FindYearlyTransactionStatus) (*pb.ApiResponseTransactionYearAmountFailed, error) {
	return h.findStatusYearlyFailed(ctx, int(req.GetYear()), "failed")
}

func (h *TransactionStatsHandler) FindMonthStatusSuccessByMerchant(ctx context.Context, req *pb.FindMonthlyTransactionStatusByMerchant) (*pb.ApiResponseTransactionMonthAmountSuccess, error) {
	return h.findStatusMonthlyByMerchant(ctx, int(req.GetYear()), int(req.GetMonth()), int(req.GetMerchantId()), "success")
}

func (h *TransactionStatsHandler) FindYearStatusSuccessByMerchant(ctx context.Context, req *pb.FindYearlyTransactionStatusByMerchant) (*pb.ApiResponseTransactionYearAmountSuccess, error) {
	return h.findStatusYearlyByMerchant(ctx, int(req.GetYear()), int(req.GetMerchantId()), "success")
}

func (h *TransactionStatsHandler) FindMonthStatusFailedByMerchant(ctx context.Context, req *pb.FindMonthlyTransactionStatusByMerchant) (*pb.ApiResponseTransactionMonthAmountFailed, error) {
	return h.findStatusMonthlyFailedByMerchant(ctx, int(req.GetYear()), int(req.GetMonth()), int(req.GetMerchantId()), "failed")
}

func (h *TransactionStatsHandler) FindYearStatusFailedByMerchant(ctx context.Context, req *pb.FindYearlyTransactionStatusByMerchant) (*pb.ApiResponseTransactionYearAmountFailed, error) {
	return h.findStatusYearlyFailedByMerchant(ctx, int(req.GetYear()), int(req.GetMerchantId()), "failed")
}

func (h *TransactionStatsHandler) FindMonthMethodSuccess(ctx context.Context, req *pb.MonthTransactionMethod) (*pb.ApiResponseTransactionMonthPaymentMethod, error) {
	return h.findMethodMonthly(ctx, int(req.GetYear()), int(req.GetMonth()), "success")
}

func (h *TransactionStatsHandler) FindYearMethodSuccess(ctx context.Context, req *pb.YearTransactionMethod) (*pb.ApiResponseTransactionYearPaymentmethod, error) {
	return h.findMethodYearly(ctx, int(req.GetYear()), "success")
}

func (h *TransactionStatsHandler) FindMonthMethodByMerchantSuccess(ctx context.Context, req *pb.MonthTransactionMethodByMerchant) (*pb.ApiResponseTransactionMonthPaymentMethod, error) {
	return h.findMethodMonthlyByMerchant(ctx, int(req.GetYear()), int(req.GetMonth()), int(req.GetMerchantId()), "success")
}

func (h *TransactionStatsHandler) FindYearMethodByMerchantSuccess(ctx context.Context, req *pb.YearTransactionMethodByMerchant) (*pb.ApiResponseTransactionYearPaymentmethod, error) {
	return h.findMethodYearlyByMerchant(ctx, int(req.GetYear()), int(req.GetMerchantId()), "success")
}

func (h *TransactionStatsHandler) FindMonthMethodFailed(ctx context.Context, req *pb.MonthTransactionMethod) (*pb.ApiResponseTransactionMonthPaymentMethod, error) {
	return h.findMethodMonthly(ctx, int(req.GetYear()), int(req.GetMonth()), "failed")
}

func (h *TransactionStatsHandler) FindYearMethodFailed(ctx context.Context, req *pb.YearTransactionMethod) (*pb.ApiResponseTransactionYearPaymentmethod, error) {
	return h.findMethodYearly(ctx, int(req.GetYear()), "failed")
}

func (h *TransactionStatsHandler) FindMonthMethodByMerchantFailed(ctx context.Context, req *pb.MonthTransactionMethodByMerchant) (*pb.ApiResponseTransactionMonthPaymentMethod, error) {
	return h.findMethodMonthlyByMerchant(ctx, int(req.GetYear()), int(req.GetMonth()), int(req.GetMerchantId()), "failed")
}

func (h *TransactionStatsHandler) FindYearMethodByMerchantFailed(ctx context.Context, req *pb.YearTransactionMethodByMerchant) (*pb.ApiResponseTransactionYearPaymentmethod, error) {
	return h.findMethodYearlyByMerchant(ctx, int(req.GetYear()), int(req.GetMerchantId()), "failed")
}

// ── Internal helpers ─────────────────────────────────────────────────────

func (h *TransactionStatsHandler) findStatusMonthly(ctx context.Context, year, month int, status string) (*pb.ApiResponseTransactionMonthAmountSuccess, error) {
	key := fmt.Sprintf("stats:reader:txn:status-month:%s:%d:%d", status, year, month)
	if cached, found := CacheGet[pb.ApiResponseTransactionMonthAmountSuccess](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetTransactionStatusMonthly(ctx, year, month, status)
	if err != nil {
		h.log.Error("findStatusMonthly failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseTransactionMonthAmountSuccess{
		Status:  "success",
		Message: "Monthly transaction status retrieved successfully",
		Data:    mapTransactionStatusMonthly(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *TransactionStatsHandler) findStatusYearly(ctx context.Context, year int, status string) (*pb.ApiResponseTransactionYearAmountSuccess, error) {
	key := fmt.Sprintf("stats:reader:txn:status-year:%s:%d", status, year)
	if cached, found := CacheGet[pb.ApiResponseTransactionYearAmountSuccess](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetTransactionStatusYearly(ctx, year, status)
	if err != nil {
		h.log.Error("findStatusYearly failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseTransactionYearAmountSuccess{
		Status:  "success",
		Message: "Yearly transaction status retrieved successfully",
		Data:    mapTransactionStatusYearly(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *TransactionStatsHandler) findStatusMonthlyFailed(ctx context.Context, year, month int, status string) (*pb.ApiResponseTransactionMonthAmountFailed, error) {
	key := fmt.Sprintf("stats:reader:txn:status-month-fail:%d:%d", year, month)
	if cached, found := CacheGet[pb.ApiResponseTransactionMonthAmountFailed](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetTransactionStatusMonthly(ctx, year, month, status)
	if err != nil {
		h.log.Error("findStatusMonthlyFailed failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseTransactionMonthAmountFailed{
		Status:  "success",
		Message: "Monthly transaction failed status retrieved successfully",
		Data:    mapTransactionStatusMonthlyFailed(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *TransactionStatsHandler) findStatusYearlyFailed(ctx context.Context, year int, status string) (*pb.ApiResponseTransactionYearAmountFailed, error) {
	key := fmt.Sprintf("stats:reader:txn:status-year-fail:%d", year)
	if cached, found := CacheGet[pb.ApiResponseTransactionYearAmountFailed](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetTransactionStatusYearly(ctx, year, status)
	if err != nil {
		h.log.Error("findStatusYearlyFailed failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseTransactionYearAmountFailed{
		Status:  "success",
		Message: "Yearly transaction failed status retrieved successfully",
		Data:    mapTransactionStatusYearlyFailed(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *TransactionStatsHandler) findStatusMonthlyByMerchant(ctx context.Context, year, month, merchantID int, status string) (*pb.ApiResponseTransactionMonthAmountSuccess, error) {
	key := fmt.Sprintf("stats:reader:txn:status-month-merchant:%s:%d:%d:%d", status, year, month, merchantID)
	if cached, found := CacheGet[pb.ApiResponseTransactionMonthAmountSuccess](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetTransactionStatusMonthlyByMerchant(ctx, year, month, merchantID, status)
	if err != nil {
		h.log.Error("findStatusMonthlyByMerchant failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseTransactionMonthAmountSuccess{
		Status:  "success",
		Message: "Monthly transaction status by merchant retrieved successfully",
		Data:    mapTransactionStatusMonthly(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *TransactionStatsHandler) findStatusYearlyByMerchant(ctx context.Context, year, merchantID int, status string) (*pb.ApiResponseTransactionYearAmountSuccess, error) {
	key := fmt.Sprintf("stats:reader:txn:status-year-merchant:%s:%d:%d", status, year, merchantID)
	if cached, found := CacheGet[pb.ApiResponseTransactionYearAmountSuccess](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetTransactionStatusYearlyByMerchant(ctx, year, merchantID, status)
	if err != nil {
		h.log.Error("findStatusYearlyByMerchant failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseTransactionYearAmountSuccess{
		Status:  "success",
		Message: "Yearly transaction status by merchant retrieved successfully",
		Data:    mapTransactionStatusYearly(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *TransactionStatsHandler) findStatusMonthlyFailedByMerchant(ctx context.Context, year, month, merchantID int, status string) (*pb.ApiResponseTransactionMonthAmountFailed, error) {
	key := fmt.Sprintf("stats:reader:txn:status-month-fail-merchant:%d:%d:%d", year, month, merchantID)
	if cached, found := CacheGet[pb.ApiResponseTransactionMonthAmountFailed](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetTransactionStatusMonthlyByMerchant(ctx, year, month, merchantID, status)
	if err != nil {
		h.log.Error("findStatusMonthlyFailedByMerchant failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseTransactionMonthAmountFailed{
		Status:  "success",
		Message: "Monthly transaction failed by merchant retrieved successfully",
		Data:    mapTransactionStatusMonthlyFailed(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *TransactionStatsHandler) findStatusYearlyFailedByMerchant(ctx context.Context, year, merchantID int, status string) (*pb.ApiResponseTransactionYearAmountFailed, error) {
	key := fmt.Sprintf("stats:reader:txn:status-year-fail-merchant:%d:%d", year, merchantID)
	if cached, found := CacheGet[pb.ApiResponseTransactionYearAmountFailed](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetTransactionStatusYearlyByMerchant(ctx, year, merchantID, status)
	if err != nil {
		h.log.Error("findStatusYearlyFailedByMerchant failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseTransactionYearAmountFailed{
		Status:  "success",
		Message: "Yearly transaction failed by merchant retrieved successfully",
		Data:    mapTransactionStatusYearlyFailed(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *TransactionStatsHandler) findMethodMonthly(ctx context.Context, year, month int, status string) (*pb.ApiResponseTransactionMonthPaymentMethod, error) {
	key := fmt.Sprintf("stats:reader:txn:method-month:%s:%d:%d", status, year, month)
	if cached, found := CacheGet[pb.ApiResponseTransactionMonthPaymentMethod](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetTransactionMethodMonthly(ctx, year, month, status)
	if err != nil {
		h.log.Error("findMethodMonthly failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseTransactionMonthPaymentMethod{
		Status:  "success",
		Message: "Monthly transaction method retrieved successfully",
		Data:    mapTransactionMethodMonthly(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *TransactionStatsHandler) findMethodYearly(ctx context.Context, year int, status string) (*pb.ApiResponseTransactionYearPaymentmethod, error) {
	key := fmt.Sprintf("stats:reader:txn:method-year:%s:%d", status, year)
	if cached, found := CacheGet[pb.ApiResponseTransactionYearPaymentmethod](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetTransactionMethodYearly(ctx, year, status)
	if err != nil {
		h.log.Error("findMethodYearly failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseTransactionYearPaymentmethod{
		Status:  "success",
		Message: "Yearly transaction method retrieved successfully",
		Data:    mapTransactionMethodYearly(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *TransactionStatsHandler) findMethodMonthlyByMerchant(ctx context.Context, year, month, merchantID int, status string) (*pb.ApiResponseTransactionMonthPaymentMethod, error) {
	key := fmt.Sprintf("stats:reader:txn:method-month-merchant:%s:%d:%d:%d", status, year, month, merchantID)
	if cached, found := CacheGet[pb.ApiResponseTransactionMonthPaymentMethod](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetTransactionMethodMonthlyByMerchant(ctx, year, month, merchantID, status)
	if err != nil {
		h.log.Error("findMethodMonthlyByMerchant failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseTransactionMonthPaymentMethod{
		Status:  "success",
		Message: "Monthly transaction method by merchant retrieved successfully",
		Data:    mapTransactionMethodMonthly(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

func (h *TransactionStatsHandler) findMethodYearlyByMerchant(ctx context.Context, year, merchantID int, status string) (*pb.ApiResponseTransactionYearPaymentmethod, error) {
	key := fmt.Sprintf("stats:reader:txn:method-year-merchant:%s:%d:%d", status, year, merchantID)
	if cached, found := CacheGet[pb.ApiResponseTransactionYearPaymentmethod](ctx, h.cache, key); found {
		return cached, nil
	}
	data, err := h.repo.GetTransactionMethodYearlyByMerchant(ctx, year, merchantID, status)
	if err != nil {
		h.log.Error("findMethodYearlyByMerchant failed", zap.Error(err))
		return nil, err
	}
	resp := &pb.ApiResponseTransactionYearPaymentmethod{
		Status:  "success",
		Message: "Yearly transaction method by merchant retrieved successfully",
		Data:    mapTransactionMethodYearly(data),
	}
	CacheSet(ctx, h.cache, key, resp)
	return resp, nil
}

// ── Mappers ─────────────────────────────────────────────────────────────

func mapTransactionMonthlySuccess(data []repository.TransactionMonthlySuccess) []*pb.TransactionMonthlySuccessResponse {
	var out []*pb.TransactionMonthlySuccessResponse
	for _, d := range data {
		out = append(out, &pb.TransactionMonthlySuccessResponse{
			Month:       d.Month,
			TotalCount:  int64(d.TotalCount),
			TotalAmount: d.TotalAmount,
		})
	}
	return out
}

func mapTransactionStatusMonthly(data []repository.TransactionStatusMonthly) []*pb.TransactionMonthlyAmountSuccess {
	var out []*pb.TransactionMonthlyAmountSuccess
	for _, d := range data {
		out = append(out, &pb.TransactionMonthlyAmountSuccess{
			Year:         d.Year,
			Month:        d.Month,
			TotalSuccess: int32(d.Count),
			TotalAmount:  int32(d.Amount),
		})
	}
	return out
}

func mapTransactionStatusYearly(data []repository.TransactionStatusYearly) []*pb.TransactionYearlyAmountSuccess {
	var out []*pb.TransactionYearlyAmountSuccess
	for _, d := range data {
		out = append(out, &pb.TransactionYearlyAmountSuccess{
			Year:         d.Year,
			TotalSuccess: int32(d.Count),
			TotalAmount:  int32(d.Amount),
		})
	}
	return out
}

func mapTransactionStatusMonthlyFailed(data []repository.TransactionStatusMonthly) []*pb.TransactionMonthlyAmountFailed {
	var out []*pb.TransactionMonthlyAmountFailed
	for _, d := range data {
		out = append(out, &pb.TransactionMonthlyAmountFailed{
			Year:        d.Year,
			Month:       d.Month,
			TotalFailed: int32(d.Count),
			TotalAmount: int32(d.Amount),
		})
	}
	return out
}

func mapTransactionStatusYearlyFailed(data []repository.TransactionStatusYearly) []*pb.TransactionYearlyAmountFailed {
	var out []*pb.TransactionYearlyAmountFailed
	for _, d := range data {
		out = append(out, &pb.TransactionYearlyAmountFailed{
			Year:        d.Year,
			TotalFailed: int32(d.Count),
			TotalAmount: int32(d.Amount),
		})
	}
	return out
}

func mapTransactionMethodMonthly(data []repository.TransactionMethodMonthly) []*pb.TransactionMonthlyMethod {
	var out []*pb.TransactionMonthlyMethod
	for _, d := range data {
		out = append(out, &pb.TransactionMonthlyMethod{
			Month:             d.Month,
			PaymentMethod:     d.PaymentMethod,
			TotalTransactions: int32(d.Transactions),
			TotalAmount:       int32(d.Amount),
		})
	}
	return out
}

func mapTransactionMethodYearly(data []repository.TransactionMethodYearly) []*pb.TransactionYearlyMethod {
	var out []*pb.TransactionYearlyMethod
	for _, d := range data {
		out = append(out, &pb.TransactionYearlyMethod{
			Year:              d.Year,
			PaymentMethod:     d.PaymentMethod,
			TotalTransactions: int32(d.Transactions),
			TotalAmount:       int32(d.Amount),
		})
	}
	return out
}
