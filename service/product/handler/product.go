package handler

import (
	"context"
	"math"

	"github.com/MamangRust/microservice-point-of-sale-product/repository"
	"github.com/MamangRust/microservice-point-of-sale-pkg/database/models"
	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	"github.com/MamangRust/microservice-point-of-sale-product/service"
	"github.com/MamangRust/microservice-point-of-sale-shared/convert"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/product_errors"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/product"
	pbcommon "github.com/MamangRust/microservice-pointofsale-grpc/pb/common"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type productHandleGrpc struct {
	pb.UnimplementedProductServiceServer
	productQueryService   service.ProductQueryService
	productCommandService service.ProductCommandService
	logger                logger.LoggerInterface
}

func NewProductHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pb.ProductServiceServer {
	return &productHandleGrpc{
		productQueryService:   service.ProductQuery,
		productCommandService: service.ProductCommand,
		logger:                logger,
	}
}

func (s *productHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllProductRequest) (*pb.ApiResponsePaginationProduct, error) {
	s.logger.Info("FindAll products called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllProducts{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	product, totalRecords, err := s.productQueryService.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAll products failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindAll products success")

	return &pb.ApiResponsePaginationProduct{
		Status:     "success",
		Message:    "Successfully fetched product",
		Data:       mapResponsesProduct(product),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *productHandleGrpc) FindByMerchant(ctx context.Context, request *pb.FindAllProductMerchantRequest) (*pb.ApiResponsePaginationProduct, error) {
	s.logger.Info("FindByMerchant products called", zap.Int32("merchantId", request.GetMerchantId()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()
	merchant_id := int(request.GetMerchantId())
	min_price := int(request.GetMinPrice())
	max_price := int(request.GetMaxPrice())

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if min_price <= 0 {
		min_price = 0
	}
	if max_price <= 0 {
		max_price = 0
	}

	reqService := requests.ProductByMerchantRequest{
		MerchantID: merchant_id,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
		MinPrice:   &min_price,
		MaxPrice:   &max_price,
	}

	product, totalRecords, err := s.productQueryService.FindByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByMerchant products failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByMerchant products success")

	return &pb.ApiResponsePaginationProduct{
		Status:     "success",
		Message:    "Successfully fetched product",
		Data:       mapResponsesProductByMerchant(product, int32(merchant_id)),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *productHandleGrpc) FindByCategory(ctx context.Context, request *pb.FindAllProductCategoryRequest) (*pb.ApiResponsePaginationProduct, error) {
	s.logger.Info("FindByCategory products called", zap.String("categoryName", request.GetCategoryName()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()
	category_name := request.GetCategoryName()
	min_price := int(request.GetMinprice())
	max_price := int(request.GetMaxprice())

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if min_price <= 0 {
		min_price = 0
	}
	if max_price <= 0 {
		max_price = 0
	}

	reqService := requests.ProductByCategoryRequest{
		Page:         page,
		PageSize:     pageSize,
		Search:       search,
		CategoryName: category_name,
		MinPrice:     &min_price,
		MaxPrice:     &max_price,
	}

	product, totalRecords, err := s.productQueryService.FindByCategory(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByCategory products failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByCategory products success")

	return &pb.ApiResponsePaginationProduct{
		Status:     "success",
		Message:    "Successfully fetched product",
		Data:       mapResponsesProductByCategory(product),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *productHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdProductRequest) (*pb.ApiResponseProduct, error) {
	s.logger.Info("FindById product called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, product_errors.ErrGrpcInvalidID
	}

	product, err := s.productQueryService.FindById(ctx, id)
	if err != nil {
		s.logger.Error("FindById product failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindById product success")

	return &pb.ApiResponseProduct{
		Status:  "success",
		Message: "Successfully fetched product",
		Data:    mapResponseProduct(product),
	}, nil
}

func (s *productHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllProductRequest) (*pb.ApiResponsePaginationProductDeleteAt, error) {
	s.logger.Info("FindByActive products called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllProducts{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	product, totalRecords, err := s.productQueryService.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByActive products failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByActive products success")

	return &pb.ApiResponsePaginationProductDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active product",
		Data:       mapResponsesProductActive(product),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *productHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllProductRequest) (*pb.ApiResponsePaginationProductDeleteAt, error) {
	s.logger.Info("FindByTrashed products called")

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllProducts{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.productQueryService.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByTrashed products failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByTrashed products success")

	return &pb.ApiResponsePaginationProductDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed product",
		Data:       mapResponsesProductTrashed(users),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *productHandleGrpc) Create(ctx context.Context, request *pb.CreateProductRequest) (*pb.ApiResponseProduct, error) {
	s.logger.Info("Create product called", zap.String("name", request.GetName()))

	req := &requests.CreateProductRequest{
		MerchantID:   int(request.GetMerchantId()),
		CategoryID:   int(request.GetCategoryId()),
		Name:         request.GetName(),
		Description:  request.GetDescription(),
		Price:        int(request.GetPrice()),
		CountInStock: int(request.GetCountInStock()),
		Brand:        request.GetBrand(),
		Weight:       int(request.GetWeight()),
		ImageProduct: request.GetImageProduct(),
	}

	if err := req.Validate(); err != nil {
		return nil, product_errors.ErrGrpcValidateCreateProduct
	}

	product, err := s.productCommandService.CreateProduct(ctx, req)
	if err != nil {
		s.logger.Error("Create product failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Create product success")

	return &pb.ApiResponseProduct{
		Status:  "success",
		Message: "Successfully created product",
		Data:    mapResponseProduct(product),
	}, nil
}

func (s *productHandleGrpc) Update(ctx context.Context, request *pb.UpdateProductRequest) (*pb.ApiResponseProduct, error) {
	s.logger.Info("Update product called", zap.Int32("id", request.GetProductId()))

	id := int(request.GetProductId())
	if id <= 0 {
		return nil, product_errors.ErrGrpcInvalidID
	}

	req := &requests.UpdateProductRequest{
		ProductID:    &id,
		MerchantID:   int(request.GetMerchantId()),
		CategoryID:   int(request.GetCategoryId()),
		Name:         request.GetName(),
		Description:  request.GetDescription(),
		Price:        int(request.GetPrice()),
		CountInStock: int(request.GetCountInStock()),
		Brand:        request.GetBrand(),
		Weight:       int(request.GetWeight()),
		ImageProduct: request.GetImageProduct(),
	}

	if err := req.Validate(); err != nil {
		return nil, product_errors.ErrGrpcValidateUpdateProduct
	}

	product, err := s.productCommandService.UpdateProduct(ctx, req)
	if err != nil {
		s.logger.Error("Update product failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Update product success")

	return &pb.ApiResponseProduct{
		Status:  "success",
		Message: "Successfully updated product",
		Data:    mapResponseProduct(product),
	}, nil
}

func (s *productHandleGrpc) TrashedProduct(ctx context.Context, request *pb.FindByIdProductRequest) (*pb.ApiResponseProductDeleteAt, error) {
	s.logger.Info("TrashedProduct called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, product_errors.ErrGrpcInvalidID
	}

	product, err := s.productCommandService.TrashProduct(ctx, id)
	if err != nil {
		s.logger.Error("TrashedProduct failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("TrashedProduct success")

	return &pb.ApiResponseProductDeleteAt{
		Status:  "success",
		Message: "Successfully trashed product",
		Data:    mapResponseProductDeleteAt(product),
	}, nil
}

func (s *productHandleGrpc) RestoreProduct(ctx context.Context, request *pb.FindByIdProductRequest) (*pb.ApiResponseProductDeleteAt, error) {
	s.logger.Info("RestoreProduct called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, product_errors.ErrGrpcInvalidID
	}

	product, err := s.productCommandService.RestoreProduct(ctx, id)
	if err != nil {
		s.logger.Error("RestoreProduct failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreProduct success")

	return &pb.ApiResponseProductDeleteAt{
		Status:  "success",
		Message: "Successfully restored product",
		Data:    mapResponseProductDeleteAt(product),
	}, nil
}

func (s *productHandleGrpc) DeleteProductPermanent(ctx context.Context, request *pb.FindByIdProductRequest) (*pb.ApiResponseProductDelete, error) {
	s.logger.Info("DeleteProductPermanent called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, product_errors.ErrGrpcInvalidID
	}

	_, err := s.productCommandService.DeleteProductPermanent(ctx, id)
	if err != nil {
		s.logger.Error("DeleteProductPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteProductPermanent success")

	return &pb.ApiResponseProductDelete{
		Status:  "success",
		Message: "Successfully deleted Product permanently",
	}, nil
}

func (s *productHandleGrpc) RestoreAllProduct(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseProductAll, error) {
	s.logger.Info("RestoreAllProduct called")

	_, err := s.productCommandService.RestoreAllProducts(ctx)
	if err != nil {
		s.logger.Error("RestoreAllProduct failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreAllProduct success")

	return &pb.ApiResponseProductAll{
		Status:  "success",
		Message: "Successfully restore all Product",
	}, nil
}

func (s *productHandleGrpc) DeleteAllProductPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseProductAll, error) {
	s.logger.Info("DeleteAllProductPermanent called")

	_, err := s.productCommandService.DeleteAllProductsPermanent(ctx)
	if err != nil {
		s.logger.Error("DeleteAllProductPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteAllProductPermanent success")

	return &pb.ApiResponseProductAll{
		Status:  "success",
		Message: "Successfully delete Product permanen",
	}, nil
}

// Map helpers
func parseStrPointer(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func parseInt32Pointer(i *int32) int32 {
	if i != nil {
		return *i
	}
	return 0
}


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

func mapResponseProduct(product *models.Product) *pb.ProductResponse {
	if product == nil {
		return nil
	}
	var createdAtStr, updatedAtStr string
	if product.CreatedAt != nil {
		createdAtStr = convert.FormatTimePtr(product.CreatedAt)
	}
	if product.UpdatedAt != nil {
		updatedAtStr = convert.FormatTimePtr(product.UpdatedAt)
	}
	return &pb.ProductResponse{
		Id:           int32(product.ProductID),
		MerchantId:   int32(product.MerchantID),
		CategoryId:   int32(product.CategoryID),
		Name:         product.Name,
		Description:  parseStrPointer(product.Description),
		Price:        int32(product.Price),
		CountInStock: int32(product.CountInStock),
		Brand:        parseStrPointer(product.Brand),
		Weight:       parseInt32Pointer(product.Weight),
		SlugProduct:  parseStrPointer(product.SlugProduct),
		ImageProduct: parseStrPointer(product.ImageProduct),
		Barcode:      parseStrPointer(product.Barcode),
		CreatedAt:    createdAtStr,
		UpdatedAt:    updatedAtStr,
	}
}

func mapResponsesProduct(products []*repository.ProductResult) []*pb.ProductResponse {
	var mappedProducts []*pb.ProductResponse
	for _, p := range products {
		if p == nil {
			continue
		}
		var createdAtStr, updatedAtStr string
		if p.CreatedAt != "" {
			createdAtStr = p.CreatedAt
		}
		if p.UpdatedAt != "" {
			updatedAtStr = p.UpdatedAt
		}
		mappedProducts = append(mappedProducts, &pb.ProductResponse{
			Id:           int32(p.ProductID),
			MerchantId:   int32(p.MerchantID),
			CategoryId:   int32(p.CategoryID),
			Name:         p.Name,
			Description:  parseStrPointer(p.Description),
			Price:        int32(p.Price),
			CountInStock: int32(p.CountInStock),
			Brand:        parseStrPointer(p.Brand),
			Weight:       parseInt32Pointer(p.Weight),
			SlugProduct:  parseStrPointer(p.SlugProduct),
			ImageProduct: parseStrPointer(p.ImageProduct),
			Barcode:      parseStrPointer(p.Barcode),
			CreatedAt:    createdAtStr,
			UpdatedAt:    updatedAtStr,
		})
	}
	return mappedProducts
}

func mapResponsesProductByMerchant(products []*repository.ProductByMerchantResult, merchantId int32) []*pb.ProductResponse {
	var mappedProducts []*pb.ProductResponse
	for _, p := range products {
		if p == nil {
			continue
		}
		var createdAtStr string
		if p.CreatedAt != "" {
			createdAtStr = p.CreatedAt
		}
		mappedProducts = append(mappedProducts, &pb.ProductResponse{
			Id:           int32(p.ProductID),
			MerchantId:   merchantId,
			Name:         p.Name,
			Description:  parseStrPointer(p.Description),
			Price:        int32(p.Price),
			CountInStock: int32(p.CountInStock),
			Brand:        parseStrPointer(p.Brand),
			ImageProduct: parseStrPointer(p.ImageProduct),
			CreatedAt:    createdAtStr,
		})
	}
	return mappedProducts
}

func mapResponsesProductByCategory(products []*repository.ProductByCategoryResult) []*pb.ProductResponse {
	var mappedProducts []*pb.ProductResponse
	for _, p := range products {
		if p == nil {
			continue
		}
		var createdAtStr, updatedAtStr string
		if p.CreatedAt != "" {
			createdAtStr = p.CreatedAt
		}
		if p.UpdatedAt != "" {
			updatedAtStr = p.UpdatedAt
		}
		mappedProducts = append(mappedProducts, &pb.ProductResponse{
			Id:           int32(p.ProductID),
			MerchantId:   int32(p.MerchantID),
			CategoryId:   int32(p.CategoryID),
			Name:         p.Name,
			Description:  parseStrPointer(p.Description),
			Price:        int32(p.Price),
			CountInStock: int32(p.CountInStock),
			Brand:        parseStrPointer(p.Brand),
			Weight:       parseInt32Pointer(p.Weight),
			SlugProduct:  parseStrPointer(p.SlugProduct),
			ImageProduct: parseStrPointer(p.ImageProduct),
			Barcode:      parseStrPointer(p.Barcode),
			CreatedAt:    createdAtStr,
			UpdatedAt:    updatedAtStr,
		})
	}
	return mappedProducts
}

func mapResponseProductDeleteAt(product *models.Product) *pb.ProductResponseDeleteAt {
	if product == nil {
		return nil
	}
	var createdAtStr, updatedAtStr string
	if product.CreatedAt != nil {
		createdAtStr = convert.FormatTimePtr(product.CreatedAt)
	}
	if product.UpdatedAt != nil {
		updatedAtStr = convert.FormatTimePtr(product.UpdatedAt)
	}
	deletedAt := convert.TimeToWrappers(product.DeletedAt)

	return &pb.ProductResponseDeleteAt{
		Id:           int32(product.ProductID),
		MerchantId:   int32(product.MerchantID),
		CategoryId:   int32(product.CategoryID),
		Name:         product.Name,
		Description:  parseStrPointer(product.Description),
		Price:        int32(product.Price),
		CountInStock: int32(product.CountInStock),
		Brand:        parseStrPointer(product.Brand),
		Weight:       parseInt32Pointer(product.Weight),
		SlugProduct:  parseStrPointer(product.SlugProduct),
		ImageProduct: parseStrPointer(product.ImageProduct),
		Barcode:      parseStrPointer(product.Barcode),
		CreatedAt:    createdAtStr,
		UpdatedAt:    updatedAtStr,
		DeletedAt:    deletedAt,
	}
}

func mapResponsesProductActive(products []*repository.ProductResultDeleteAt) []*pb.ProductResponseDeleteAt {
	var mappedProducts []*pb.ProductResponseDeleteAt
	for _, p := range products {
		if p == nil {
			continue
		}
		var createdAtStr, updatedAtStr string
		if p.CreatedAt != "" {
			createdAtStr = p.CreatedAt
		}
		if p.UpdatedAt != "" {
			updatedAtStr = p.UpdatedAt
		}
		deletedAt := convert.StrValToWrappers(&p.DeletedAt)

		mappedProducts = append(mappedProducts, &pb.ProductResponseDeleteAt{
			Id:           int32(p.ProductID),
			MerchantId:   int32(p.MerchantID),
			CategoryId:   int32(p.CategoryID),
			Name:         p.Name,
			Description:  parseStrPointer(p.Description),
			Price:        int32(p.Price),
			CountInStock: int32(p.CountInStock),
			Brand:        parseStrPointer(p.Brand),
			Weight:       parseInt32Pointer(p.Weight),
			SlugProduct:  parseStrPointer(p.SlugProduct),
			ImageProduct: parseStrPointer(p.ImageProduct),
			Barcode:      parseStrPointer(p.Barcode),
			CreatedAt:    createdAtStr,
			UpdatedAt:    updatedAtStr,
			DeletedAt:    deletedAt,
		})
	}
	return mappedProducts
}

func mapResponsesProductTrashed(products []*repository.ProductResultDeleteAt) []*pb.ProductResponseDeleteAt {
	var mappedProducts []*pb.ProductResponseDeleteAt
	for _, p := range products {
		if p == nil {
			continue
		}
		var createdAtStr, updatedAtStr string
		if p.CreatedAt != "" {
			createdAtStr = p.CreatedAt
		}
		if p.UpdatedAt != "" {
			updatedAtStr = p.UpdatedAt
		}
		deletedAt := convert.StrValToWrappers(&p.DeletedAt)

		mappedProducts = append(mappedProducts, &pb.ProductResponseDeleteAt{
			Id:           int32(p.ProductID),
			MerchantId:   int32(p.MerchantID),
			CategoryId:   int32(p.CategoryID),
			Name:         p.Name,
			Description:  parseStrPointer(p.Description),
			Price:        int32(p.Price),
			CountInStock: int32(p.CountInStock),
			Brand:        parseStrPointer(p.Brand),
			Weight:       parseInt32Pointer(p.Weight),
			SlugProduct:  parseStrPointer(p.SlugProduct),
			ImageProduct: parseStrPointer(p.ImageProduct),
			Barcode:      parseStrPointer(p.Barcode),
			CreatedAt:    createdAtStr,
			UpdatedAt:    updatedAtStr,
			DeletedAt:    deletedAt,
		})
	}
	return mappedProducts
}
