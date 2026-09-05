package handler

import (
	"context"
	"math"

	"github.com/MamangRust/microservice-point-of-sale-merchant/service"
	"github.com/MamangRust/microservice-point-of-sale-merchant/repository"
	"github.com/MamangRust/microservice-point-of-sale-pkg/database/models"
	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/merchant_errors"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/merchant"
	pbcommon "github.com/MamangRust/microservice-pointofsale-grpc/pb/common"
	"time"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantHandleGrpc struct {
	pb.UnimplementedMerchantServiceServer
	merchantQuery   service.MerchantQueryService
	merchantCommand service.MerchantCommandService
	logger          logger.LoggerInterface
}

func NewMerchantHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pb.MerchantServiceServer {
	return &merchantHandleGrpc{
		merchantQuery:   service.MerchantQuery,
		merchantCommand: service.MerchantCommand,
		logger:          logger,
	}
}

func mapSqlNullString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func mapSqlNullTime(t *time.Time) string {
	if t != nil {
		return t.Format("2006-01-02 15:04:05")
	}
	return ""
}

func (s *merchantHandleGrpc) FindAll(ctx context.Context, req *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchant, error) {
	s.logger.Info("FindAll merchant called", zap.Int32("page", req.GetPage()))

	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchants{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchants, totalRecords, err := s.merchantQuery.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAll merchant failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindAll merchant success")

	return &pb.ApiResponsePaginationMerchant{
		Status:     "success",
		Message:    "Successfully fetched merchant record",
		Data:       mapResponsesGetMerchantsRow(merchants),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *merchantHandleGrpc) FindByActive(ctx context.Context, req *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantDeleteAt, error) {
	s.logger.Info("FindByActive merchants called", zap.Int32("page", req.GetPage()))

	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchants{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.merchantQuery.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByActive merchants failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByActive merchants success")

	return &pb.ApiResponsePaginationMerchantDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched merchant record",
		Data:       mapResponsesGetMerchantsActiveRow(res),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *merchantHandleGrpc) FindByTrashed(ctx context.Context, req *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantDeleteAt, error) {
	s.logger.Info("FindByTrashed merchants called")

	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchants{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.merchantQuery.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByTrashed merchants failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByTrashed merchants success")

	return &pb.ApiResponsePaginationMerchantDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched merchant record",
		Data:       mapResponsesGetMerchantsTrashedRow(res),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *merchantHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchant, error) {
	s.logger.Info("FindById merchant called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, merchant_errors.ErrGrpcInvalidID
	}

	merchant, err := s.merchantQuery.FindById(ctx, id)
	if err != nil {
		s.logger.Error("FindById merchant failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindById merchant success")

	return &pb.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully fetched merchant",
		Data:    mapResponseMerchant(merchant),
	}, nil
}

func (s *merchantHandleGrpc) Create(ctx context.Context, request *pb.CreateMerchantRequest) (*pb.ApiResponseMerchant, error) {
	s.logger.Info("Create merchant called", zap.String("name", request.GetName()))

	req := &requests.CreateMerchantRequest{
		UserID:       int(request.GetUserId()),
		Name:         request.GetName(),
		Description:  request.GetDescription(),
		Address:      request.GetAddress(),
		ContactEmail: request.GetContactEmail(),
		ContactPhone: request.GetContactPhone(),
		Status:       request.GetStatus(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchant_errors.ErrGrpcValidateCreateMerchant
	}

	merchant, err := s.merchantCommand.CreateMerchant(ctx, req)
	if err != nil {
		s.logger.Error("Create merchant failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Create merchant success")

	return &pb.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully created merchant",
		Data:    mapResponseMerchant(merchant),
	}, nil
}

func (s *merchantHandleGrpc) Update(ctx context.Context, request *pb.UpdateMerchantRequest) (*pb.ApiResponseMerchant, error) {
	s.logger.Info("Update merchant called", zap.Int32("id", request.GetMerchantId()))

	id := int(request.GetMerchantId())
	if id <= 0 {
		return nil, merchant_errors.ErrGrpcInvalidID
	}

	req := &requests.UpdateMerchantRequest{
		MerchantID:   &id,
		UserID:       int(request.GetUserId()),
		Name:         request.GetName(),
		Description:  request.GetDescription(),
		Address:      request.GetAddress(),
		ContactEmail: request.GetContactEmail(),
		ContactPhone: request.GetContactPhone(),
		Status:       request.GetStatus(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchant_errors.ErrGrpcValidateUpdateMerchant
	}

	merchant, err := s.merchantCommand.UpdateMerchant(ctx, req)
	if err != nil {
		s.logger.Error("Update merchant failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Update merchant success")

	return &pb.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully updated merchant",
		Data:    mapResponseMerchant(merchant),
	}, nil
}

func (s *merchantHandleGrpc) UpdateMerchantStatus(ctx context.Context, req *pb.UpdateMerchantStatusRequest) (*pb.ApiResponseMerchant, error) {
	s.logger.Info("UpdateMerchantStatus merchant called", zap.Int32("id", req.GetMerchantId()))

	id := int(req.GetMerchantId())
	if id <= 0 {
		return nil, merchant_errors.ErrGrpcInvalidID
	}

	request := requests.UpdateMerchantStatusRequest{
		MerchantID: &id,
		Status:     req.GetStatus(),
	}

	if err := request.Validate(); err != nil {
		return nil, merchant_errors.ErrGrpcValidateUpdateMerchantStatus
	}

	merchant, err := s.merchantCommand.UpdateMerchantStatus(ctx, &request)
	if err != nil {
		s.logger.Error("UpdateMerchantStatus merchant failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("UpdateMerchantStatus merchant success")

	return &pb.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully updated merchant status",
		Data:    mapResponseMerchant(merchant),
	}, nil
}

func (s *merchantHandleGrpc) TrashedMerchant(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDeleteAt, error) {
	s.logger.Info("TrashedMerchant called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, merchant_errors.ErrGrpcInvalidID
	}

	merchant, err := s.merchantCommand.TrashedMerchant(ctx, id)
	if err != nil {
		s.logger.Error("TrashedMerchant failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("TrashedMerchant success")

	return &pb.ApiResponseMerchantDeleteAt{
		Status:  "success",
		Message: "Successfully trashed merchant",
		Data:    mapResponseMerchantDeleteAt(merchant),
	}, nil
}

func (s *merchantHandleGrpc) RestoreMerchant(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchant, error) {
	s.logger.Info("RestoreMerchant called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, merchant_errors.ErrGrpcInvalidID
	}

	merchant, err := s.merchantCommand.RestoreMerchant(ctx, id)
	if err != nil {
		s.logger.Error("RestoreMerchant failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreMerchant success")

	return &pb.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully restored merchant",
		Data:    mapResponseMerchant(merchant),
	}, nil
}

func (s *merchantHandleGrpc) DeleteMerchantPermanent(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDelete, error) {
	s.logger.Info("DeleteMerchantPermanent called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, merchant_errors.ErrGrpcInvalidID
	}

	_, err := s.merchantCommand.DeleteMerchantPermanent(ctx, id)
	if err != nil {
		s.logger.Error("DeleteMerchantPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteMerchantPermanent success")

	return &pb.ApiResponseMerchantDelete{
		Status:  "success",
		Message: "Successfully deleted merchant permanently",
	}, nil
}

func (s *merchantHandleGrpc) RestoreAllMerchant(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	s.logger.Info("RestoreAllMerchant called")

	_, err := s.merchantCommand.RestoreAllMerchant(ctx)
	if err != nil {
		s.logger.Error("RestoreAllMerchant failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreAllMerchant success")

	return &pb.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully restore all merchant",
	}, nil
}

func (s *merchantHandleGrpc) DeleteAllMerchantPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	s.logger.Info("DeleteAllMerchantPermanent called")

	_, err := s.merchantCommand.DeleteAllMerchantPermanent(ctx)
	if err != nil {
		s.logger.Error("DeleteAllMerchantPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteAllMerchantPermanent success")

	return &pb.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully delete merchant permanen",
	}, nil
}

// Map helpers
func mapPaginationMeta(meta *pbcommon.PaginationMeta) *pbcommon.PaginationMeta {
	if meta == nil {
		return nil
	}
	return &pbcommon.PaginationMeta{
		CurrentPage:  meta.CurrentPage,
		PageSize:     meta.PageSize,
		TotalPages:   meta.TotalPages,
		TotalRecords: meta.TotalRecords,
	}
}

func mapResponseMerchant(merchant *models.Merchant) *pb.MerchantResponse {
	if merchant == nil {
		return nil
	}
	return &pb.MerchantResponse{
		Id:           int32(merchant.MerchantID),
		UserId:       int32(merchant.UserID),
		Name:         merchant.Name,
		Description:  mapSqlNullString(merchant.Description),
		Address:      mapSqlNullString(merchant.Address),
		ContactEmail: mapSqlNullString(merchant.ContactEmail),
		ContactPhone: mapSqlNullString(merchant.ContactPhone),
		Status:       merchant.Status,
		UpdatedAt:    mapSqlNullTime(merchant.UpdatedAt),
	}
}

func mapResponsesGetMerchantsRow(merchants []*repository.MerchantResult) []*pb.MerchantResponse {
	var mapped []*pb.MerchantResponse
	for _, m := range merchants {
		mapped = append(mapped, &pb.MerchantResponse{
			Id:           int32(m.MerchantID),
			UserId:       int32(m.UserID),
			Name:         m.Name,
			Description:  mapSqlNullString(m.Description),
			Address:      mapSqlNullString(m.Address),
			ContactEmail: mapSqlNullString(m.ContactEmail),
			ContactPhone: mapSqlNullString(m.ContactPhone),
			Status:       m.Status,
			CreatedAt:    m.CreatedAt,
			UpdatedAt:    m.UpdatedAt,
		})
	}
	return mapped
}

func mapResponseMerchantDeleteAt(merchant *models.Merchant) *pb.MerchantResponseDeleteAt {
	if merchant == nil {
		return nil
	}
	return &pb.MerchantResponseDeleteAt{
		Id:           int32(merchant.MerchantID),
		UserId:       int32(merchant.UserID),
		Name:         merchant.Name,
		Description:  mapSqlNullString(merchant.Description),
		Address:      mapSqlNullString(merchant.Address),
		ContactEmail: mapSqlNullString(merchant.ContactEmail),
		ContactPhone: mapSqlNullString(merchant.ContactPhone),
		Status:       merchant.Status,
		UpdatedAt:    mapSqlNullTime(merchant.UpdatedAt),
		DeletedAt:    mapSqlNullTime(merchant.DeletedAt),
	}
}

func mapResponsesGetMerchantsActiveRow(merchants []*repository.MerchantResultDeleteAt) []*pb.MerchantResponseDeleteAt {
	var mapped []*pb.MerchantResponseDeleteAt
	for _, m := range merchants {
		mapped = append(mapped, &pb.MerchantResponseDeleteAt{
			Id:           int32(m.MerchantID),
			UserId:       int32(m.UserID),
			Name:         m.Name,
			Description:  mapSqlNullString(m.Description),
			Address:      mapSqlNullString(m.Address),
			ContactEmail: mapSqlNullString(m.ContactEmail),
			ContactPhone: mapSqlNullString(m.ContactPhone),
			Status:       m.Status,
			CreatedAt:    m.CreatedAt,
			UpdatedAt:    m.UpdatedAt,
			DeletedAt:    m.DeletedAt,
		})
	}
	return mapped
}

func mapResponsesGetMerchantsTrashedRow(merchants []*repository.MerchantResultDeleteAt) []*pb.MerchantResponseDeleteAt {
	var mapped []*pb.MerchantResponseDeleteAt
	for _, m := range merchants {
		mapped = append(mapped, &pb.MerchantResponseDeleteAt{
			Id:           int32(m.MerchantID),
			UserId:       int32(m.UserID),
			Name:         m.Name,
			Description:  mapSqlNullString(m.Description),
			Address:      mapSqlNullString(m.Address),
			ContactEmail: mapSqlNullString(m.ContactEmail),
			ContactPhone: mapSqlNullString(m.ContactPhone),
			Status:       m.Status,
			CreatedAt:    m.CreatedAt,
			UpdatedAt:    m.UpdatedAt,
			DeletedAt:    m.DeletedAt,
		})
	}
	return mapped
}
