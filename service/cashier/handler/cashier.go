package handler

import (
	"context"
	"math"

	"github.com/MamangRust/microservice-point-of-sale-cashier/service"
	"github.com/MamangRust/microservice-point-of-sale-cashier/repository"
	"github.com/MamangRust/microservice-point-of-sale-pkg/database/models"
	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	"github.com/MamangRust/microservice-point-of-sale-shared/convert"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/cashier_errors"
	pbcommon "github.com/MamangRust/microservice-pointofsale-grpc/pb/common"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/cashier"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type cashierHandleGrpc struct {
	pb.UnimplementedCashierServiceServer
	cashierQuery   service.CashierQueryService
	cashierCommand service.CashierCommandService
	logger         logger.LoggerInterface
}

func NewCashierHandleGrpc(
	svc *service.Service,
	logger logger.LoggerInterface,
) pb.CashierServiceServer {
	return &cashierHandleGrpc{
		cashierQuery:   svc.CashierQuery,
		cashierCommand: svc.CashierCommand,
		logger:         logger,
	}
}

func (s *cashierHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllCashierRequest) (*pb.ApiResponsePaginationCashier, error) {
	s.logger.Info("FindAll cashier called", zap.Int32("page", request.GetPage()), zap.Int32("pageSize", request.GetPageSize()))
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	reqService := requests.FindAllCashiers{Search: search, Page: page, PageSize: pageSize}
	cashier, totalRecords, err := s.cashierQuery.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAll cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}
	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	return &pb.ApiResponsePaginationCashier{
		Status: "success", Message: "Successfully fetched cashier",
		Data: mapResponsesCashier(cashier), Pagination: mapPaginationMeta(&pbcommon.PaginationMeta{CurrentPage: int32(page), PageSize: int32(pageSize), TotalPages: int32(totalPages), TotalRecords: int32(*totalRecords)}),
	}, nil
}

func (s *cashierHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdCashierRequest) (*pb.ApiResponseCashier, error) {
	s.logger.Info("FindById cashier called", zap.Int32("id", request.GetId()))
	id := int(request.GetId())
	if id <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}
	cashier, err := s.cashierQuery.FindById(ctx, id)
	if err != nil {
		s.logger.Error("FindById cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}
	return &pb.ApiResponseCashier{Status: "success", Message: "Successfully fetched cashier", Data: mapResponseCashier(cashier)}, nil
}

func (s *cashierHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllCashierRequest) (*pb.ApiResponsePaginationCashierDeleteAt, error) {
	s.logger.Info("FindByActive cashier called", zap.Int32("page", request.GetPage()))
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	reqService := requests.FindAllCashiers{Search: search, Page: page, PageSize: pageSize}
	cashier, totalRecords, err := s.cashierQuery.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByActive cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}
	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	return &pb.ApiResponsePaginationCashierDeleteAt{
		Status: "success", Message: "Successfully fetched active cashier",
		Data: mapResponsesCashierActive(cashier), Pagination: mapPaginationMeta(&pbcommon.PaginationMeta{CurrentPage: int32(page), PageSize: int32(pageSize), TotalPages: int32(totalPages), TotalRecords: int32(*totalRecords)}),
	}, nil
}

func (s *cashierHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllCashierRequest) (*pb.ApiResponsePaginationCashierDeleteAt, error) {
	s.logger.Info("FindByTrashed cashier called")
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	reqService := requests.FindAllCashiers{Search: search, Page: page, PageSize: pageSize}
	users, totalRecords, err := s.cashierQuery.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByTrashed cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}
	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	return &pb.ApiResponsePaginationCashierDeleteAt{
		Status: "success", Message: "Successfully fetched trashed cashier",
		Data: mapResponsesCashierTrashed(users), Pagination: mapPaginationMeta(&pbcommon.PaginationMeta{CurrentPage: int32(page), PageSize: int32(pageSize), TotalPages: int32(totalPages), TotalRecords: int32(*totalRecords)}),
	}, nil
}

func (s *cashierHandleGrpc) FindByMerchant(ctx context.Context, request *pb.FindByMerchantCashierRequest) (*pb.ApiResponsePaginationCashier, error) {
	s.logger.Info("FindByMerchant cashier called", zap.Int32("merchantId", request.GetMerchantId()))
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()
	merchantID := int(request.GetMerchantId())
	if merchantID <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidMerchantId
	}
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	reqService := requests.FindAllCashierMerchant{Search: search, Page: page, PageSize: pageSize, MerchantID: merchantID}
	cashier, totalRecords, err := s.cashierQuery.FindByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByMerchant cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}
	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	return &pb.ApiResponsePaginationCashier{
		Status: "success", Message: "Successfully fetched cashier",
		Data: mapResponsesCashierByMerchant(cashier), Pagination: mapPaginationMeta(&pbcommon.PaginationMeta{CurrentPage: int32(page), PageSize: int32(pageSize), TotalPages: int32(totalPages), TotalRecords: int32(*totalRecords)}),
	}, nil
}

func (s *cashierHandleGrpc) CreateCashier(ctx context.Context, request *pb.CreateCashierRequest) (*pb.ApiResponseCashier, error) {
	s.logger.Info("CreateCashier called", zap.String("name", request.GetName()))
	req := &requests.CreateCashierRequest{Name: request.GetName(), MerchantID: int(request.GetMerchantId()), UserID: int(request.GetUserId())}
	if err := req.Validate(); err != nil {
		return nil, cashier_errors.ErrGrpcValidateCreateCashier
	}
	cashier, err := s.cashierCommand.CreateCashier(ctx, req)
	if err != nil {
		s.logger.Error("CreateCashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}
	return &pb.ApiResponseCashier{Status: "success", Message: "Successfully created cashier", Data: mapResponseCashier(cashier)}, nil
}

func (s *cashierHandleGrpc) UpdateCashier(ctx context.Context, request *pb.UpdateCashierRequest) (*pb.ApiResponseCashier, error) {
	s.logger.Info("UpdateCashier called", zap.Int32("id", request.GetCashierId()))
	id := int(request.GetCashierId())
	if id <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}
	req := &requests.UpdateCashierRequest{CashierID: &id, Name: request.GetName()}
	if err := req.Validate(); err != nil {
		return nil, cashier_errors.ErrGrpcValidateUpdateCashier
	}
	cashier, err := s.cashierCommand.UpdateCashier(ctx, req)
	if err != nil {
		s.logger.Error("UpdateCashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}
	return &pb.ApiResponseCashier{Status: "success", Message: "Successfully updated cashier", Data: mapResponseCashier(cashier)}, nil
}

func (s *cashierHandleGrpc) TrashedCashier(ctx context.Context, request *pb.FindByIdCashierRequest) (*pb.ApiResponseCashierDeleteAt, error) {
	s.logger.Info("TrashedCashier called", zap.Int32("id", request.GetId()))
	id := int(request.GetId())
	if id <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}
	cashier, err := s.cashierCommand.TrashedCashier(ctx, id)
	if err != nil {
		s.logger.Error("TrashedCashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}
	return &pb.ApiResponseCashierDeleteAt{Status: "success", Message: "Successfully trashed cashier", Data: mapResponseCashierDeleteAt(cashier)}, nil
}

func (s *cashierHandleGrpc) RestoreCashier(ctx context.Context, request *pb.FindByIdCashierRequest) (*pb.ApiResponseCashierDeleteAt, error) {
	s.logger.Info("RestoreCashier called", zap.Int32("id", request.GetId()))
	id := int(request.GetId())
	if id <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}
	cashier, err := s.cashierCommand.RestoreCashier(ctx, id)
	if err != nil {
		s.logger.Error("RestoreCashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}
	return &pb.ApiResponseCashierDeleteAt{Status: "success", Message: "Successfully restored cashier", Data: mapResponseCashierDeleteAt(cashier)}, nil
}

func (s *cashierHandleGrpc) DeleteCashierPermanent(ctx context.Context, request *pb.FindByIdCashierRequest) (*pb.ApiResponseCashierDelete, error) {
	s.logger.Info("DeleteCashierPermanent called", zap.Int32("id", request.GetId()))
	id := int(request.GetId())
	if id <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}
	_, err := s.cashierCommand.DeleteCashierPermanent(ctx, id)
	if err != nil {
		s.logger.Error("DeleteCashierPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}
	return &pb.ApiResponseCashierDelete{Status: "success", Message: "Successfully deleted cashier permanently"}, nil
}

func (s *cashierHandleGrpc) RestoreAllCashier(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseCashierAll, error) {
	s.logger.Info("RestoreAllCashier called")
	_, err := s.cashierCommand.RestoreAllCashier(ctx)
	if err != nil {
		s.logger.Error("RestoreAllCashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}
	return &pb.ApiResponseCashierAll{Status: "success", Message: "Successfully restore all cashier"}, nil
}

func (s *cashierHandleGrpc) DeleteAllCashierPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseCashierAll, error) {
	s.logger.Info("DeleteAllCashierPermanent called")
	_, err := s.cashierCommand.DeleteAllCashierPermanent(ctx)
	if err != nil {
		s.logger.Error("DeleteAllCashierPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}
	return &pb.ApiResponseCashierAll{Status: "success", Message: "Successfully delete cashier permanen"}, nil
}

// Map helpers
func mapPaginationMeta(meta *pbcommon.PaginationMeta) *pbcommon.PaginationMeta {
	if meta == nil {
		return nil
	}
	return &pbcommon.PaginationMeta{CurrentPage: meta.CurrentPage, PageSize: meta.PageSize, TotalPages: meta.TotalPages, TotalRecords: meta.TotalRecords}
}

func fmtCashStr(s *string) string {
	return convert.StrVal(s)
}


func mapResponseCashier(cashier *models.Cashier) *pb.CashierResponse {
	if cashier == nil {
		return nil
	}
	return &pb.CashierResponse{
		Id: cashier.CashierID, MerchantId: cashier.MerchantID, Name: cashier.Name,
		CreatedAt: convert.FormatTimePtr(cashier.CreatedAt), UpdatedAt: convert.FormatTimePtr(cashier.UpdatedAt),
	}
}

func mapResponsesCashier(cashiers []*repository.CashierResult) []*pb.CashierResponse {
	var res []*pb.CashierResponse
	for _, c := range cashiers {
		res = append(res, &pb.CashierResponse{
			Id: c.CashierID, MerchantId: c.MerchantID, Name: c.Name,
			CreatedAt: fmtCashStr(c.CreatedAt), UpdatedAt: fmtCashStr(c.UpdatedAt),
		})
	}
	return res
}

func mapResponsesCashierByMerchant(cashiers []*repository.CashierResult) []*pb.CashierResponse {
	return mapResponsesCashier(cashiers)
}

func mapResponseCashierDeleteAt(cashier *models.Cashier) *pb.CashierResponseDeleteAt {
	if cashier == nil {
		return nil
	}
	return &pb.CashierResponseDeleteAt{
		Id: cashier.CashierID, MerchantId: cashier.MerchantID, Name: cashier.Name,
		CreatedAt: convert.FormatTimePtr(cashier.CreatedAt), UpdatedAt: convert.FormatTimePtr(cashier.UpdatedAt),		DeletedAt: convert.TimeToWrappers(cashier.DeletedAt),
	}
}

func mapResponsesCashierActive(cashiers []*repository.CashierResultDeleteAt) []*pb.CashierResponseDeleteAt {
	var res []*pb.CashierResponseDeleteAt
	for _, c := range cashiers {
		res = append(res, &pb.CashierResponseDeleteAt{
			Id: c.CashierID, MerchantId: c.MerchantID, Name: c.Name,
			CreatedAt: fmtCashStr(c.CreatedAt), UpdatedAt: fmtCashStr(c.UpdatedAt),		DeletedAt: convert.StrValToWrappers(c.DeletedAt),
		})
	}
	return res
}

func mapResponsesCashierTrashed(cashiers []*repository.CashierResultDeleteAt) []*pb.CashierResponseDeleteAt {
	return mapResponsesCashierActive(cashiers)
}
