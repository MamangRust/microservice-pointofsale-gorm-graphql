// Package transaction reads transaction status and payment-method statistics
// from ClickHouse's transaction_daily table. It replaces the transaction slice
// of the old single stats-reader Repository god interface.
package transaction

import "context"

// Result models mirror the ClickHouse aggregations in clickhouse.go. Counts and
// amounts are int64/uint64: ClickHouse UInt64/Int64 columns scan into them
// directly (clickhouse-go v2 rejects implicit UInt64→int64 conversion).

type TransactionMonthAmount struct {
	Year        string
	Month       string
	TotalCount  uint64
	TotalAmount int64
}

type TransactionYearAmount struct {
	Year        string
	TotalCount  uint64
	TotalAmount int64
}

type TransactionMonthMethod struct {
	Month             string
	PaymentMethod     string
	TotalTransactions uint64
	TotalAmount       int64
}

type TransactionYearMethod struct {
	Year              string
	PaymentMethod     string
	TotalTransactions uint64
	TotalAmount       int64
}

// Repository serves the transaction rpcs of TransactionStatsStatusService and
// TransactionStatsMethodService.
type Repository interface {
	GetMonthStatusSuccess(ctx context.Context, year, month int) ([]TransactionMonthAmount, error)
	GetYearStatusSuccess(ctx context.Context, year int) ([]TransactionYearAmount, error)
	GetMonthStatusFailed(ctx context.Context, year, month int) ([]TransactionMonthAmount, error)
	GetYearStatusFailed(ctx context.Context, year int) ([]TransactionYearAmount, error)
	GetMonthStatusSuccessByMerchant(ctx context.Context, year, month, merchantID int) ([]TransactionMonthAmount, error)
	GetYearStatusSuccessByMerchant(ctx context.Context, year, merchantID int) ([]TransactionYearAmount, error)
	GetMonthStatusFailedByMerchant(ctx context.Context, year, month, merchantID int) ([]TransactionMonthAmount, error)
	GetYearStatusFailedByMerchant(ctx context.Context, year, merchantID int) ([]TransactionYearAmount, error)
	GetMonthMethodSuccess(ctx context.Context, year, month int) ([]TransactionMonthMethod, error)
	GetYearMethodSuccess(ctx context.Context, year int) ([]TransactionYearMethod, error)
	GetMonthMethodFailed(ctx context.Context, year, month int) ([]TransactionMonthMethod, error)
	GetYearMethodFailed(ctx context.Context, year int) ([]TransactionYearMethod, error)
	GetMonthMethodByMerchantSuccess(ctx context.Context, year, month, merchantID int) ([]TransactionMonthMethod, error)
	GetYearMethodByMerchantSuccess(ctx context.Context, year, merchantID int) ([]TransactionYearMethod, error)
	GetMonthMethodByMerchantFailed(ctx context.Context, year, month, merchantID int) ([]TransactionMonthMethod, error)
	GetYearMethodByMerchantFailed(ctx context.Context, year, merchantID int) ([]TransactionYearMethod, error)
}
