package protomapper

import (
	pbcommon "github.com/MamangRust/microservice-point-of-sale-pb/common"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/response"

	pb "github.com/MamangRust/microservice-point-of-sale-pb/transaction"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type transactionProtoMapper struct{}

func NewTransactionProtoMapper() *transactionProtoMapper {
	return &transactionProtoMapper{}
}

func (t *transactionProtoMapper) ToProtoResponseTransaction(status string, message string, trans *response.TransactionResponse) *pb.ApiResponseTransaction {
	return &pb.ApiResponseTransaction{
		Status:  status,
		Message: message,
		Data:    t.mapResponseTransaction(trans),
	}
}

func (t *transactionProtoMapper) ToProtoResponsesTransaction(status string, message string, transList []*response.TransactionResponse) *pb.ApiResponsesTransaction {
	return &pb.ApiResponsesTransaction{
		Status:  status,
		Message: message,
		Data:    t.mapResponsesTransaction(transList),
	}
}

func (t *transactionProtoMapper) ToProtoResponseTransactionDeleteAt(status string, message string, trans *response.TransactionResponseDeleteAt) *pb.ApiResponseTransactionDeleteAt {
	return &pb.ApiResponseTransactionDeleteAt{
		Status:  status,
		Message: message,
		Data:    t.mapResponseTransactionDeleteAt(trans),
	}
}

func (t *transactionProtoMapper) ToProtoResponseTransactionDelete(status string, message string) *pb.ApiResponseTransactionDelete {
	return &pb.ApiResponseTransactionDelete{
		Status:  status,
		Message: message,
	}
}

func (t *transactionProtoMapper) ToProtoResponseTransactionAll(status string, message string) *pb.ApiResponseTransactionAll {
	return &pb.ApiResponseTransactionAll{
		Status:  status,
		Message: message,
	}
}

func (t *transactionProtoMapper) ToProtoResponsePaginationTransactionDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, transactions []*response.TransactionResponseDeleteAt) *pb.ApiResponsePaginationTransactionDeleteAt {
	return &pb.ApiResponsePaginationTransactionDeleteAt{
		Status:     status,
		Message:    message,
		Data:       t.mapResponsesTransactionDeleteAt(transactions),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (t *transactionProtoMapper) ToProtoResponsePaginationTransaction(pagination *pbcommon.PaginationMeta, status string, message string, transactions []*response.TransactionResponse) *pb.ApiResponsePaginationTransaction {
	return &pb.ApiResponsePaginationTransaction{
		Status:     status,
		Message:    message,
		Data:       t.mapResponsesTransaction(transactions),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (t *transactionProtoMapper) mapResponseTransaction(transaction *response.TransactionResponse) *pb.TransactionResponse {
	return &pb.TransactionResponse{
		Id:            int32(transaction.ID),
		OrderId:       int32(transaction.OrderID),
		MerchantId:    int32(transaction.MerchantID),
		PaymentMethod: transaction.PaymentMethod,
		Amount:        int32(transaction.Amount),
		ChangeAmount:  int32(transaction.ChangeAmount),
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
	}
}

func (t *transactionProtoMapper) mapResponsesTransaction(transactions []*response.TransactionResponse) []*pb.TransactionResponse {
	var mappedTransactions []*pb.TransactionResponse

	for _, transaction := range transactions {
		mappedTransactions = append(mappedTransactions, t.mapResponseTransaction(transaction))
	}

	return mappedTransactions
}

func (t *transactionProtoMapper) mapResponseTransactionDeleteAt(transaction *response.TransactionResponseDeleteAt) *pb.TransactionResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if transaction.DeletedAt != nil {
		deletedAt = wrapperspb.String(*transaction.DeletedAt)
	}

	return &pb.TransactionResponseDeleteAt{
		Id:            int32(transaction.ID),
		OrderId:       int32(transaction.OrderID),
		MerchantId:    int32(transaction.MerchantID),
		PaymentMethod: transaction.PaymentMethod,
		Amount:        int32(transaction.Amount),
		ChangeAmount:  int32(transaction.ChangeAmount),
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
		DeletedAt:     deletedAt,
	}
}

func (t *transactionProtoMapper) mapResponsesTransactionDeleteAt(transactions []*response.TransactionResponseDeleteAt) []*pb.TransactionResponseDeleteAt {
	var mappedTransactions []*pb.TransactionResponseDeleteAt

	for _, transaction := range transactions {
		mappedTransactions = append(mappedTransactions, t.mapResponseTransactionDeleteAt(transaction))
	}

	return mappedTransactions
}
