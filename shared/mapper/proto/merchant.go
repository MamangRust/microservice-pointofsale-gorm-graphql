package protomapper

import (
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/response"
	pbcommon "github.com/MamangRust/microservice-pointofsale-grpc/pb/common"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/merchant"
)

type merchantProtoMapper struct{}

func NewMerchantProtoMaper() *merchantProtoMapper {
	return &merchantProtoMapper{}
}

func (m *merchantProtoMapper) ToProtoResponseMerchant(status string, message string, pbResponse *response.MerchantResponse) *pb.ApiResponseMerchant {
	return &pb.ApiResponseMerchant{
		Status:  status,
		Message: message,
		Data:    m.mapResponseMerchant(pbResponse),
	}
}

func (m *merchantProtoMapper) ToProtoResponsesMerchant(status string, message string, pbResponse []*response.MerchantResponse) *pb.ApiResponsesMerchant {
	return &pb.ApiResponsesMerchant{
		Status:  status,
		Message: message,
		Data:    m.mapResponsesMerchant(pbResponse),
	}
}

func (m *merchantProtoMapper) ToProtoResponseMerchantDeleteAt(status string, message string, pbResponse *response.MerchantResponseDeleteAt) *pb.ApiResponseMerchantDeleteAt {
	return &pb.ApiResponseMerchantDeleteAt{
		Status:  status,
		Message: message,
		Data:    m.mapResponseMerchantDeleteAt(pbResponse),
	}
}

func (m *merchantProtoMapper) ToProtoResponseMerchantDelete(status string, message string) *pb.ApiResponseMerchantDelete {
	return &pb.ApiResponseMerchantDelete{
		Status:  status,
		Message: message,
	}
}

func (m *merchantProtoMapper) ToProtoResponseMerchantAll(status string, message string) *pb.ApiResponseMerchantAll {
	return &pb.ApiResponseMerchantAll{
		Status:  status,
		Message: message,
	}
}

func (m *merchantProtoMapper) ToProtoResponsePaginationMerchantDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, merchants []*response.MerchantResponseDeleteAt) *pb.ApiResponsePaginationMerchantDeleteAt {
	return &pb.ApiResponsePaginationMerchantDeleteAt{
		Status:     status,
		Message:    message,
		Data:       m.mapResponsesMerchantDeleteAt(merchants),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *merchantProtoMapper) ToProtoResponsePaginationMerchant(pagination *pbcommon.PaginationMeta, status string, message string, merchants []*response.MerchantResponse) *pb.ApiResponsePaginationMerchant {
	return &pb.ApiResponsePaginationMerchant{
		Status:     status,
		Message:    message,
		Data:       m.mapResponsesMerchant(merchants),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *merchantProtoMapper) mapResponseMerchant(merchant *response.MerchantResponse) *pb.MerchantResponse {
	return &pb.MerchantResponse{
		Id:           int32(merchant.ID),
		UserId:       int32(merchant.UserID),
		Name:         merchant.Name,
		Description:  merchant.Description,
		Address:      merchant.Address,
		ContactEmail: merchant.ContactEmail,
		ContactPhone: merchant.ContactPhone,
		Status:       merchant.Status,
		CreatedAt:    merchant.CreatedAt,
		UpdatedAt:    merchant.UpdatedAt,
	}
}

func (m *merchantProtoMapper) mapResponsesMerchant(merchants []*response.MerchantResponse) []*pb.MerchantResponse {
	var mappedMerchants []*pb.MerchantResponse

	for _, merchant := range merchants {
		mappedMerchants = append(mappedMerchants, m.mapResponseMerchant(merchant))
	}

	return mappedMerchants
}

func (m *merchantProtoMapper) mapResponseMerchantDeleteAt(merchant *response.MerchantResponseDeleteAt) *pb.MerchantResponseDeleteAt {
	return &pb.MerchantResponseDeleteAt{
		Id:           int32(merchant.ID),
		UserId:       int32(merchant.UserID),
		Name:         merchant.Name,
		Description:  merchant.Description,
		Address:      merchant.Address,
		ContactEmail: merchant.ContactEmail,
		ContactPhone: merchant.ContactPhone,
		Status:       merchant.Status,
		CreatedAt:    merchant.CreatedAt,
		UpdatedAt:    merchant.UpdatedAt,
		DeletedAt:    merchant.DeletedAt,
	}
}

func (m *merchantProtoMapper) mapResponsesMerchantDeleteAt(merchants []*response.MerchantResponseDeleteAt) []*pb.MerchantResponseDeleteAt {
	var mappedMerchants []*pb.MerchantResponseDeleteAt

	for _, merchant := range merchants {
		mappedMerchants = append(mappedMerchants, m.mapResponseMerchantDeleteAt(merchant))
	}

	return mappedMerchants
}
