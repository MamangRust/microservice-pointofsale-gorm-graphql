package protomapper

import (
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/response"
	pbcommon "github.com/MamangRust/microservice-pointofsale-grpc/pb/common"
	"google.golang.org/protobuf/types/known/wrapperspb"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/merchant"
)

type merchantDocumentProtoMapper struct{}

func NewMerchantDocumentProtoMapper() *merchantDocumentProtoMapper {
	return &merchantDocumentProtoMapper{}
}

func (m *merchantDocumentProtoMapper) ToProtoResponseMerchantDocument(status string, message string, doc *response.MerchantDocumentResponse) *pb.ApiResponseMerchantDocument {
	return &pb.ApiResponseMerchantDocument{
		Status:  status,
		Message: message,
		Data:    m.mapMerchantDocument(doc),
	}
}

func (m *merchantDocumentProtoMapper) ToProtoResponsesMerchantDocument(status string, message string, docs []*response.MerchantDocumentResponse) *pb.ApiResponsesMerchantDocument {
	return &pb.ApiResponsesMerchantDocument{
		Status:  status,
		Message: message,
		Data:    m.mapMerchantDocuments(docs),
	}
}

func (m *merchantDocumentProtoMapper) ToProtoResponsePaginationMerchantDocument(pagination *pbcommon.PaginationMeta, status string, message string, docs []*response.MerchantDocumentResponse) *pb.ApiResponsePaginationMerchantDocument {
	return &pb.ApiResponsePaginationMerchantDocument{
		Status:     status,
		Message:    message,
		Data:       m.mapMerchantDocuments(docs),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *merchantDocumentProtoMapper) ToProtoResponsePaginationMerchantDocumentDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, docs []*response.MerchantDocumentResponseDeleteAt) *pb.ApiResponsePaginationMerchantDocumentAt {
	return &pb.ApiResponsePaginationMerchantDocumentAt{
		Status:     status,
		Message:    message,
		Data:       m.mapMerchantDocumentsDeleteAt(docs),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *merchantDocumentProtoMapper) ToProtoResponseMerchantDocumentDelete(status string, message string) *pb.ApiResponseMerchantDocumentDelete {
	return &pb.ApiResponseMerchantDocumentDelete{
		Status:  status,
		Message: message,
	}
}

func (m *merchantDocumentProtoMapper) ToProtoResponseMerchantDocumentAll(status string, message string) *pb.ApiResponseMerchantDocumentAll {
	return &pb.ApiResponseMerchantDocumentAll{
		Status:  status,
		Message: message,
	}
}

func (m *merchantDocumentProtoMapper) mapMerchantDocument(doc *response.MerchantDocumentResponse) *pb.MerchantDocument {
	return &pb.MerchantDocument{
		DocumentId:   int32(doc.ID),
		MerchantId:   int32(doc.MerchantID),
		DocumentType: doc.DocumentType,
		DocumentUrl:  doc.DocumentURL,
		Status:       doc.Status,
		Note:         doc.Note,
		UploadedAt:   doc.CreatedAt,
		UpdatedAt:    doc.UpdatedAt,
	}
}

func (m *merchantDocumentProtoMapper) mapMerchantDocuments(docs []*response.MerchantDocumentResponse) []*pb.MerchantDocument {
	var res []*pb.MerchantDocument
	for _, doc := range docs {
		res = append(res, m.mapMerchantDocument(doc))
	}
	return res
}

func (m *merchantDocumentProtoMapper) mapMerchantDocumentDeleteAt(doc *response.MerchantDocumentResponseDeleteAt) *pb.MerchantDocumentDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if doc.DeletedAt != nil {
		deletedAt = wrapperspb.String(*doc.DeletedAt)
	}

	return &pb.MerchantDocumentDeleteAt{
		DocumentId:   int32(doc.ID),
		MerchantId:   int32(doc.MerchantID),
		DocumentType: doc.DocumentType,
		DocumentUrl:  doc.DocumentURL,
		Status:       doc.Status,
		Note:         doc.Note,
		UploadedAt:   doc.CreatedAt,
		UpdatedAt:    doc.UpdatedAt,
		DeletedAt:    deletedAt,
	}
}

func (m *merchantDocumentProtoMapper) mapMerchantDocumentsDeleteAt(docs []*response.MerchantDocumentResponseDeleteAt) []*pb.MerchantDocumentDeleteAt {
	var res []*pb.MerchantDocumentDeleteAt
	for _, doc := range docs {
		res = append(res, m.mapMerchantDocumentDeleteAt(doc))
	}
	return res
}
