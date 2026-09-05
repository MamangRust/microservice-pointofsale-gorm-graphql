package handler

import (
	"context"
	"math"

	"github.com/MamangRust/microservice-point-of-sale-category/service"
	"github.com/MamangRust/microservice-point-of-sale-category/repository"
	"github.com/MamangRust/microservice-point-of-sale-pkg/database/models"
	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	"github.com/MamangRust/microservice-point-of-sale-shared/convert"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/category_errors"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/category"
	pbcommon "github.com/MamangRust/microservice-pointofsale-grpc/pb/common"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type categoryHandleGrpc struct {
	pb.UnimplementedCategoryServiceServer
	categoryQuery           service.CategoryQueryService
	categoryCommand         service.CategoryCommandService
	logger                  logger.LoggerInterface
}

func NewCategoryHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pb.CategoryServiceServer {
	return &categoryHandleGrpc{
		categoryQuery:           service.CategoryQuery,
		categoryCommand:         service.CategoryCommand,
		logger:                  logger,
	}
}

func (s *categoryHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllCategoryRequest) (*pb.ApiResponsePaginationCategory, error) {
	s.logger.Info("FindAll categories called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCategory{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	category, totalRecords, err := s.categoryQuery.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAll categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindAll categories success", zap.Int("count", len(category)))

	return &pb.ApiResponsePaginationCategory{
		Status:     "success",
		Message:    "Successfully fetched categories",
		Data:       mapResponsesCategory(category),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *categoryHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdCategoryRequest) (*pb.ApiResponseCategory, error) {
	s.logger.Info("FindById category called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	category, err := s.categoryQuery.FindById(ctx, id)
	if err != nil {
		s.logger.Error("FindById category failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindById category success", zap.Int("id", id))

	return &pb.ApiResponseCategory{
		Status:  "success",
		Message: "Successfully fetched category",
		Data:    mapResponseCategory(category),
	}, nil
}

func (s *categoryHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllCategoryRequest) (*pb.ApiResponsePaginationCategoryDeleteAt, error) {
	s.logger.Info("FindByActive categories called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCategory{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	categories, totalRecords, err := s.categoryQuery.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByActive categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByActive categories success")

	return &pb.ApiResponsePaginationCategoryDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active categories",
		Data:       mapResponsesCategoryActive(categories),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *categoryHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllCategoryRequest) (*pb.ApiResponsePaginationCategoryDeleteAt, error) {
	s.logger.Info("FindByTrashed categories called")

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCategory{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	categories, totalRecords, err := s.categoryQuery.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByTrashed categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByTrashed categories success")

	return &pb.ApiResponsePaginationCategoryDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed categories",
		Data:       mapResponsesCategoryTrashed(categories),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *categoryHandleGrpc) Create(ctx context.Context, request *pb.CreateCategoryRequest) (*pb.ApiResponseCategory, error) {
	s.logger.Info("Create category called", zap.String("name", request.GetName()))

	req := &requests.CreateCategoryRequest{
		Name:        request.GetName(),
		Description: request.GetDescription(),
	}

	if err := req.Validate(); err != nil {
		return nil, category_errors.ErrGrpcValidateCreateCategory
	}

	category, err := s.categoryCommand.CreateCategory(ctx, req)
	if err != nil {
		s.logger.Error("Create category failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Create category success")

	return &pb.ApiResponseCategory{
		Status:  "success",
		Message: "Successfully created category",
		Data:    mapResponseCategory(category),
	}, nil
}

func (s *categoryHandleGrpc) Update(ctx context.Context, request *pb.UpdateCategoryRequest) (*pb.ApiResponseCategory, error) {
	s.logger.Info("Update category called", zap.Int32("id", request.GetCategoryId()))

	id := int(request.GetCategoryId())
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	req := &requests.UpdateCategoryRequest{
		CategoryID:  &id,
		Name:        request.GetName(),
		Description: request.GetDescription(),
	}

	if err := req.Validate(); err != nil {
		return nil, category_errors.ErrGrpcValidateUpdateCategory
	}

	category, err := s.categoryCommand.UpdateCategory(ctx, req)
	if err != nil {
		s.logger.Error("Update category failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Update category success")

	return &pb.ApiResponseCategory{
		Status:  "success",
		Message: "Successfully updated category",
		Data:    mapResponseCategory(category),
	}, nil
}

func (s *categoryHandleGrpc) TrashedCategory(ctx context.Context, request *pb.FindByIdCategoryRequest) (*pb.ApiResponseCategoryDeleteAt, error) {
	s.logger.Info("TrashedCategory called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	category, err := s.categoryCommand.TrashedCategory(ctx, id)
	if err != nil {
		s.logger.Error("TrashedCategory failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("TrashedCategory success")

	return &pb.ApiResponseCategoryDeleteAt{
		Status:  "success",
		Message: "Successfully trashed category",
		Data:    mapResponseCategoryDeleteAt(category),
	}, nil
}

func (s *categoryHandleGrpc) RestoreCategory(ctx context.Context, request *pb.FindByIdCategoryRequest) (*pb.ApiResponseCategoryDeleteAt, error) {
	s.logger.Info("RestoreCategory called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	category, err := s.categoryCommand.RestoreCategory(ctx, id)
	if err != nil {
		s.logger.Error("RestoreCategory failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreCategory success")

	return &pb.ApiResponseCategoryDeleteAt{
		Status:  "success",
		Message: "Successfully restored category",
		Data:    mapResponseCategoryDeleteAt(category),
	}, nil
}

func (s *categoryHandleGrpc) DeleteCategoryPermanent(ctx context.Context, request *pb.FindByIdCategoryRequest) (*pb.ApiResponseCategoryDelete, error) {
	s.logger.Info("DeleteCategoryPermanent called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	_, err := s.categoryCommand.DeleteCategoryPermanent(ctx, id)
	if err != nil {
		s.logger.Error("DeleteCategoryPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteCategoryPermanent success")

	return &pb.ApiResponseCategoryDelete{
		Status:  "success",
		Message: "Successfully deleted category permanently",
	}, nil
}

func (s *categoryHandleGrpc) RestoreAllCategory(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseCategoryAll, error) {
	s.logger.Info("RestoreAllCategory called")

	_, err := s.categoryCommand.RestoreAllCategories(ctx)
	if err != nil {
		s.logger.Error("RestoreAllCategory failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreAllCategory success")

	return &pb.ApiResponseCategoryAll{
		Status:  "success",
		Message: "Successfully restore all category",
	}, nil
}

func (s *categoryHandleGrpc) DeleteAllCategoryPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseCategoryAll, error) {
	s.logger.Info("DeleteAllCategoryPermanent called")

	_, err := s.categoryCommand.DeleteAllCategoriesPermanent(ctx)
	if err != nil {
		s.logger.Error("DeleteAllCategoryPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteAllCategoryPermanent success")

	return &pb.ApiResponseCategoryAll{
		Status:  "success",
		Message: "Successfully delete category permanen",
	}, nil
}

// Internal map helpers
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

func fmtCatStr(s *string) string {
	return convert.StrVal(s)
}


func mapResponseCategory(category *models.Category) *pb.CategoryResponse {
	if category == nil {
		return nil
	}
	return &pb.CategoryResponse{
		Id:           category.CategoryID,
		Name:         category.Name,
		Description:  fmtCatStr(category.Description),
		SlugCategory: fmtCatStr(category.SlugCategory),
		CreatedAt:    convert.FormatTimePtr(category.CreatedAt),
		UpdatedAt:    convert.FormatTimePtr(category.UpdatedAt),
	}
}

func mapResponseGetCategory(category *repository.CategoryResult) *pb.CategoryResponse {
	if category == nil {
		return nil
	}
	return &pb.CategoryResponse{
		Id:           category.CategoryID,
		Name:         category.Name,
		Description:  fmtCatStr(category.Description),
		SlugCategory: fmtCatStr(category.SlugCategory),
		CreatedAt:    fmtCatStr(category.CreatedAt),
		UpdatedAt:    fmtCatStr(category.UpdatedAt),
	}
}

func mapResponsesCategory(categories []*repository.CategoryResult) []*pb.CategoryResponse {
	var res []*pb.CategoryResponse
	for _, c := range categories {
		res = append(res, mapResponseGetCategory(c))
	}
	return res
}

func mapResponseCategoryDeleteAt(category *models.Category) *pb.CategoryResponseDeleteAt {
	if category == nil {
		return nil
	}
	return &pb.CategoryResponseDeleteAt{
		Id:           category.CategoryID,
		Name:         category.Name,
		Description:  fmtCatStr(category.Description),
		SlugCategory: fmtCatStr(category.SlugCategory),
		CreatedAt:    convert.FormatTimePtr(category.CreatedAt),
		UpdatedAt:    convert.FormatTimePtr(category.UpdatedAt),
		DeletedAt:    convert.TimeToWrappers(category.DeletedAt),
	}
}

func mapResponseGetCategoryActive(category *repository.CategoryResultDeleteAt) *pb.CategoryResponseDeleteAt {
	if category == nil {
		return nil
	}
	return &pb.CategoryResponseDeleteAt{
		Id:           category.CategoryID,
		Name:         category.Name,
		Description:  fmtCatStr(category.Description),
		SlugCategory: fmtCatStr(category.SlugCategory),
		CreatedAt:    fmtCatStr(category.CreatedAt),
		UpdatedAt:    fmtCatStr(category.UpdatedAt),
		DeletedAt:    convert.StrValToWrappers(category.DeletedAt),
	}
}

func mapResponsesCategoryActive(categories []*repository.CategoryResultDeleteAt) []*pb.CategoryResponseDeleteAt {
	var res []*pb.CategoryResponseDeleteAt
	for _, c := range categories {
		res = append(res, mapResponseGetCategoryActive(c))
	}
	return res
}

func mapResponseGetCategoryTrashed(category *repository.CategoryResultDeleteAt) *pb.CategoryResponseDeleteAt {
	if category == nil {
		return nil
	}
	return &pb.CategoryResponseDeleteAt{
		Id:           category.CategoryID,
		Name:         category.Name,
		Description:  fmtCatStr(category.Description),
		SlugCategory: fmtCatStr(category.SlugCategory),
		CreatedAt:    fmtCatStr(category.CreatedAt),
		UpdatedAt:    fmtCatStr(category.UpdatedAt),
		DeletedAt:    convert.StrValToWrappers(category.DeletedAt),
	}
}

func mapResponsesCategoryTrashed(categories []*repository.CategoryResultDeleteAt) []*pb.CategoryResponseDeleteAt {
	var res []*pb.CategoryResponseDeleteAt
	for _, c := range categories {
		res = append(res, mapResponseGetCategoryTrashed(c))
	}
	return res
}
























