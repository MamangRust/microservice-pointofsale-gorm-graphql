package protomapper

import (
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/response"
	pbcommon "github.com/MamangRust/microservice-pointofsale-grpc/pb/common"

	"google.golang.org/protobuf/types/known/wrapperspb"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/transaction"
	pbstats "github.com/MamangRust/microservice-pointofsale-grpc/pb/stats"
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

func (t *transactionProtoMapper) ToProtoResponseMonthAmountSuccess(status string, message string, row []*response.TransactionMonthlyAmountSuccessResponse) *pbstats.ApiResponseTransactionMonthAmountSuccess {
	return &pbstats.ApiResponseTransactionMonthAmountSuccess{
		Status:  status,
		Message: message,
		Data:    t.mapResponsesTransactionMonthlyAmountSuccess(row),
	}
}

func (t *transactionProtoMapper) ToProtoResponseYearAmountSuccess(status string, message string, row []*response.TransactionYearlyAmountSuccessResponse) *pbstats.ApiResponseTransactionYearAmountSuccess {
	return &pbstats.ApiResponseTransactionYearAmountSuccess{
		Status:  status,
		Message: message,
		Data:    t.mapResponsesTransactionYearlyAmountSuccess(row),
	}
}

func (t *transactionProtoMapper) ToProtoResponseMonthAmountFailed(status string, message string, row []*response.TransactionMonthlyAmountFailedResponse) *pbstats.ApiResponseTransactionMonthAmountFailed {
	return &pbstats.ApiResponseTransactionMonthAmountFailed{
		Status:  status,
		Message: message,
		Data:    t.mapResponsesTransactionMonthlyAmountFailed(row),
	}
}

func (t *transactionProtoMapper) ToProtoResponseYearAmountFailed(status string, message string, row []*response.TransactionYearlyAmountFailedResponse) *pbstats.ApiResponseTransactionYearAmountFailed {
	return &pbstats.ApiResponseTransactionYearAmountFailed{
		Status:  status,
		Message: message,
		Data:    t.mapResponsesTransactionYearlyAmountFailed(row),
	}
}

func (t *transactionProtoMapper) ToProtoResponseMonthMethod(status string, message string, row []*response.TransactionMonthlyMethodResponse) *pbstats.ApiResponseTransactionMonthPaymentMethod {
	return &pbstats.ApiResponseTransactionMonthPaymentMethod{
		Status:  status,
		Message: message,
		Data:    t.mapResponsesTransactionMonthlyMethod(row),
	}
}

func (t *transactionProtoMapper) ToProtoResponseYearMethod(status string, message string, row []*response.TransactionYearlyMethodResponse) *pbstats.ApiResponseTransactionYearPaymentmethod {
	return &pbstats.ApiResponseTransactionYearPaymentmethod{
		Status:  status,
		Message: message,
		Data:    t.mapResponsesTransactionYearlyMethod(row),
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

func (s *transactionProtoMapper) mapResponseTransactionMonthAmountSuccess(row *response.TransactionMonthlyAmountSuccessResponse) *pbstats.TransactionMonthlyAmountSuccess {
	return &pbstats.TransactionMonthlyAmountSuccess{
		Year:         row.Year,
		Month:        row.Month,
		TotalSuccess: int32(row.TotalSuccess),
		TotalAmount:  int32(row.TotalAmount),
	}
}

func (s *transactionProtoMapper) mapResponsesTransactionMonthlyAmountSuccess(rows []*response.TransactionMonthlyAmountSuccessResponse) []*pbstats.TransactionMonthlyAmountSuccess {
	var transaction []*pbstats.TransactionMonthlyAmountSuccess

	for _, row := range rows {
		transaction = append(transaction, s.mapResponseTransactionMonthAmountSuccess(row))
	}

	return transaction
}

func (s *transactionProtoMapper) mapResponseTransactionYearAmountSuccess(row *response.TransactionYearlyAmountSuccessResponse) *pbstats.TransactionYearlyAmountSuccess {
	return &pbstats.TransactionYearlyAmountSuccess{
		Year:         row.Year,
		TotalSuccess: int32(row.TotalSuccess),
		TotalAmount:  int32(row.TotalAmount),
	}
}

func (s *transactionProtoMapper) mapResponsesTransactionYearlyAmountSuccess(rows []*response.TransactionYearlyAmountSuccessResponse) []*pbstats.TransactionYearlyAmountSuccess {
	var transaction []*pbstats.TransactionYearlyAmountSuccess

	for _, row := range rows {
		transaction = append(transaction, s.mapResponseTransactionYearAmountSuccess(row))
	}

	return transaction
}

func (s *transactionProtoMapper) mapResponseTransactionMonthAmountFailed(row *response.TransactionMonthlyAmountFailedResponse) *pbstats.TransactionMonthlyAmountFailed {
	return &pbstats.TransactionMonthlyAmountFailed{
		Year:        row.Year,
		Month:       row.Month,
		TotalFailed: int32(row.TotalFailed),
		TotalAmount: int32(row.TotalAmount),
	}
}

func (s *transactionProtoMapper) mapResponsesTransactionMonthlyAmountFailed(rows []*response.TransactionMonthlyAmountFailedResponse) []*pbstats.TransactionMonthlyAmountFailed {
	var transaction []*pbstats.TransactionMonthlyAmountFailed

	for _, row := range rows {
		transaction = append(transaction, s.mapResponseTransactionMonthAmountFailed(row))
	}

	return transaction
}

func (s *transactionProtoMapper) mapResponseTransactionYearAmountFailed(row *response.TransactionYearlyAmountFailedResponse) *pbstats.TransactionYearlyAmountFailed {
	return &pbstats.TransactionYearlyAmountFailed{
		Year:        row.Year,
		TotalFailed: int32(row.TotalFailed),
		TotalAmount: int32(row.TotalAmount),
	}
}

func (s *transactionProtoMapper) mapResponsesTransactionYearlyAmountFailed(rows []*response.TransactionYearlyAmountFailedResponse) []*pbstats.TransactionYearlyAmountFailed {
	var transaction []*pbstats.TransactionYearlyAmountFailed

	for _, row := range rows {
		transaction = append(transaction, s.mapResponseTransactionYearAmountFailed(row))
	}

	return transaction
}

func (s *transactionProtoMapper) mapResponseTransactionMonthMethod(row *response.TransactionMonthlyMethodResponse) *pbstats.TransactionMonthlyMethod {
	return &pbstats.TransactionMonthlyMethod{
		Month:             row.Month,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func (s *transactionProtoMapper) mapResponsesTransactionMonthlyMethod(rows []*response.TransactionMonthlyMethodResponse) []*pbstats.TransactionMonthlyMethod {
	var transaction []*pbstats.TransactionMonthlyMethod

	for _, row := range rows {
		transaction = append(transaction, s.mapResponseTransactionMonthMethod(row))
	}

	return transaction
}

func (s *transactionProtoMapper) mapResponseTransactionYearMethod(row *response.TransactionYearlyMethodResponse) *pbstats.TransactionYearlyMethod {
	return &pbstats.TransactionYearlyMethod{
		Year:              row.Year,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func (s *transactionProtoMapper) mapResponsesTransactionYearlyMethod(rows []*response.TransactionYearlyMethodResponse) []*pbstats.TransactionYearlyMethod {
	var transaction []*pbstats.TransactionYearlyMethod

	for _, row := range rows {
		transaction = append(transaction, s.mapResponseTransactionYearMethod(row))
	}

	return transaction
}
