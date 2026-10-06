// Package cashier reads cashier sales statistics from ClickHouse's order_daily
// table. It replaces the cashier slice of the old single stats-reader
// Repository god interface.
package cashier

import "context"

// Result models mirror the ClickHouse aggregations in clickhouse.go. Counts and
// amounts are int64/uint64: ClickHouse UInt64/Int64 columns scan into them
// directly (clickhouse-go v2 rejects implicit UInt64→int64 conversion).

type CashierMonthTotalSales struct {
	Year       string
	Month      string
	TotalSales int64
}

type CashierYearTotalSales struct {
	Year       string
	TotalSales int64
}

type CashierMonthSales struct {
	Month      string
	CashierID  uint64
	OrderCount uint64
	TotalSales int64
}

type CashierYearSales struct {
	Year       string
	CashierID  uint64
	OrderCount uint64
	TotalSales int64
}

// Repository serves the cashier rpcs of CashierStatsService,
// CashierStatsByMerchantService and CashierStatsByIdService. cashier_name is
// not stored in ClickHouse — the front end resolves it via the cashier domain
// service.
type Repository interface {
	GetMonthlyTotalSales(ctx context.Context, year, month int) ([]CashierMonthTotalSales, error)
	GetYearlyTotalSales(ctx context.Context, year int) ([]CashierYearTotalSales, error)
	GetMonthSales(ctx context.Context, year int) ([]CashierMonthSales, error)
	GetYearSales(ctx context.Context, year int) ([]CashierYearSales, error)
	GetMonthlyTotalSalesByMerchant(ctx context.Context, year, month, merchantID int) ([]CashierMonthTotalSales, error)
	GetYearlyTotalSalesByMerchant(ctx context.Context, year, merchantID int) ([]CashierYearTotalSales, error)
	GetMonthSalesByMerchant(ctx context.Context, year, merchantID int) ([]CashierMonthSales, error)
	GetYearSalesByMerchant(ctx context.Context, year, merchantID int) ([]CashierYearSales, error)
	GetMonthlyTotalSalesById(ctx context.Context, year, month, cashierID int) ([]CashierMonthTotalSales, error)
	GetYearlyTotalSalesById(ctx context.Context, year, cashierID int) ([]CashierYearTotalSales, error)
	GetMonthSalesById(ctx context.Context, year, cashierID int) ([]CashierMonthSales, error)
	GetYearSalesById(ctx context.Context, year, cashierID int) ([]CashierYearSales, error)
}
