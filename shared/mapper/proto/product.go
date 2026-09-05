package protomapper

import (
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/response"
	pbcommon "github.com/MamangRust/microservice-pointofsale-grpc/pb/common"

	"google.golang.org/protobuf/types/known/wrapperspb"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/product"
)

type productProtoMapper struct{}

func NewProductProtoMapper() *productProtoMapper {
	return &productProtoMapper{}
}

func (p *productProtoMapper) ToProtoResponseProduct(status string, message string, pbResponse *response.ProductResponse) *pb.ApiResponseProduct {
	return &pb.ApiResponseProduct{
		Status:  status,
		Message: message,
		Data:    p.mapResponseProduct(pbResponse),
	}
}

func (p *productProtoMapper) ToProtoResponsesProduct(status string, message string, pbResponse []*response.ProductResponse) *pb.ApiResponsesProduct {
	return &pb.ApiResponsesProduct{
		Status:  status,
		Message: message,
		Data:    p.mapResponsesProduct(pbResponse),
	}
}

func (p *productProtoMapper) ToProtoResponseProductDeleteAt(status string, message string, pbResponse *response.ProductResponseDeleteAt) *pb.ApiResponseProductDeleteAt {
	return &pb.ApiResponseProductDeleteAt{
		Status:  status,
		Message: message,
		Data:    p.mapResponseProductDeleteAt(pbResponse),
	}
}

func (p *productProtoMapper) ToProtoResponseProductDelete(status string, message string) *pb.ApiResponseProductDelete {
	return &pb.ApiResponseProductDelete{
		Status:  status,
		Message: message,
	}
}

func (p *productProtoMapper) ToProtoResponseProductAll(status string, message string) *pb.ApiResponseProductAll {
	return &pb.ApiResponseProductAll{
		Status:  status,
		Message: message,
	}
}

func (p *productProtoMapper) ToProtoResponsePaginationProductDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, products []*response.ProductResponseDeleteAt) *pb.ApiResponsePaginationProductDeleteAt {
	return &pb.ApiResponsePaginationProductDeleteAt{
		Status:     status,
		Message:    message,
		Data:       p.mapResponsesProductDeleteAt(products),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (p *productProtoMapper) ToProtoResponsePaginationProduct(pagination *pbcommon.PaginationMeta, status string, message string, products []*response.ProductResponse) *pb.ApiResponsePaginationProduct {
	return &pb.ApiResponsePaginationProduct{
		Status:     status,
		Message:    message,
		Data:       p.mapResponsesProduct(products),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (p *productProtoMapper) mapResponseProduct(product *response.ProductResponse) *pb.ProductResponse {
	return &pb.ProductResponse{
		Id:           int32(product.ID),
		MerchantId:   int32(product.MerchantID),
		CategoryId:   int32(product.CategoryID),
		Name:         product.Name,
		Description:  product.Description,
		Price:        int32(product.Price),
		CountInStock: int32(product.CountInStock),
		Brand:        product.Brand,
		Weight:       int32(product.Weight),
		SlugProduct:  product.SlugProduct,
		ImageProduct: product.ImageProduct,
		Barcode:      product.Barcode,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}
}

func (p *productProtoMapper) mapResponsesProduct(products []*response.ProductResponse) []*pb.ProductResponse {
	var mappedProducts []*pb.ProductResponse

	for _, product := range products {
		mappedProducts = append(mappedProducts, p.mapResponseProduct(product))
	}

	return mappedProducts
}

func (p *productProtoMapper) mapResponseProductDeleteAt(product *response.ProductResponseDeleteAt) *pb.ProductResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if product.DeleteAt != nil {
		deletedAt = wrapperspb.String(*product.DeleteAt)
	}

	return &pb.ProductResponseDeleteAt{
		Id:           int32(product.ID),
		MerchantId:   int32(product.MerchantID),
		CategoryId:   int32(product.CategoryID),
		Name:         product.Name,
		Description:  product.Description,
		Price:        int32(product.Price),
		CountInStock: int32(product.CountInStock),
		Brand:        product.Brand,
		Weight:       int32(product.Weight),
		SlugProduct:  product.SlugProduct,
		ImageProduct: product.ImageProduct,
		Barcode:      product.Barcode,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
		DeletedAt:    deletedAt,
	}
}

func (p *productProtoMapper) mapResponsesProductDeleteAt(products []*response.ProductResponseDeleteAt) []*pb.ProductResponseDeleteAt {
	var mappedProducts []*pb.ProductResponseDeleteAt

	for _, product := range products {
		mappedProducts = append(mappedProducts, p.mapResponseProductDeleteAt(product))
	}

	return mappedProducts
}
