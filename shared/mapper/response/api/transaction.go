package response_api

import (
	pb "github.com/MamangRust/microservice-point-of-sale-pb/transaction"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/response"
)

type transactionResponseMapper struct {
}

func NewTransactionResponseMapper() *transactionResponseMapper {
	return &transactionResponseMapper{}
}

func (t *transactionResponseMapper) ToResponseTransaction(transaction *pb.TransactionResponse) *response.TransactionResponse {
	return &response.TransactionResponse{
		ID:            int(transaction.Id),
		OrderID:       int(transaction.OrderId),
		MerchantID:    int(transaction.MerchantId),
		PaymentMethod: transaction.PaymentMethod,
		Amount:        int(transaction.Amount),
		ChangeAmount:  int(transaction.ChangeAmount),
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
	}
}

func (t *transactionResponseMapper) ToResponsesTransaction(transactions []*pb.TransactionResponse) []*response.TransactionResponse {
	var mappedTransactions []*response.TransactionResponse

	for _, transaction := range transactions {
		mappedTransactions = append(mappedTransactions, t.ToResponseTransaction(transaction))
	}

	return mappedTransactions
}

func (t *transactionResponseMapper) ToResponseTransactionDeleteAt(transaction *pb.TransactionResponseDeleteAt) *response.TransactionResponseDeleteAt {
	var deletedAt string
	if transaction.DeletedAt != nil {
		deletedAt = transaction.DeletedAt.Value
	}

	return &response.TransactionResponseDeleteAt{
		ID:            int(transaction.Id),
		OrderID:       int(transaction.OrderId),
		MerchantID:    int(transaction.MerchantId),
		PaymentMethod: transaction.PaymentMethod,
		Amount:        int(transaction.Amount),
		ChangeAmount:  int(transaction.ChangeAmount),
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
		DeletedAt:     &deletedAt,
	}
}

func (t *transactionResponseMapper) ToResponsesTransactionDeleteAt(transactions []*pb.TransactionResponseDeleteAt) []*response.TransactionResponseDeleteAt {
	var mappedTransactions []*response.TransactionResponseDeleteAt

	for _, transaction := range transactions {
		mappedTransactions = append(mappedTransactions, t.ToResponseTransactionDeleteAt(transaction))
	}

	return mappedTransactions
}

func (t *transactionResponseMapper) ToApiResponseTransaction(pbResponse *pb.ApiResponseTransaction) *response.ApiResponseTransaction {
	return &response.ApiResponseTransaction{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    t.ToResponseTransaction(pbResponse.Data),
	}
}

func (t *transactionResponseMapper) ToApiResponseTransactionDeleteAt(pbResponse *pb.ApiResponseTransactionDeleteAt) *response.ApiResponseTransactionDeleteAt {
	return &response.ApiResponseTransactionDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    t.ToResponseTransactionDeleteAt(pbResponse.Data),
	}
}

func (t *transactionResponseMapper) ToApiResponsesTransaction(pbResponse *pb.ApiResponsesTransaction) *response.ApiResponsesTransaction {
	return &response.ApiResponsesTransaction{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    t.ToResponsesTransaction(pbResponse.Data),
	}
}

func (t *transactionResponseMapper) ToApiResponseTransactionDelete(pbResponse *pb.ApiResponseTransactionDelete) *response.ApiResponseTransactionDelete {
	return &response.ApiResponseTransactionDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (t *transactionResponseMapper) ToApiResponseTransactionAll(pbResponse *pb.ApiResponseTransactionAll) *response.ApiResponseTransactionAll {
	return &response.ApiResponseTransactionAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (t *transactionResponseMapper) ToApiResponsePaginationTransactionDeleteAt(pbResponse *pb.ApiResponsePaginationTransactionDeleteAt) *response.ApiResponsePaginationTransactionDeleteAt {
	return &response.ApiResponsePaginationTransactionDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       t.ToResponsesTransactionDeleteAt(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}

func (t *transactionResponseMapper) ToApiResponsePaginationTransaction(pbResponse *pb.ApiResponsePaginationTransaction) *response.ApiResponsePaginationTransaction {
	return &response.ApiResponsePaginationTransaction{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       t.ToResponsesTransaction(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}
