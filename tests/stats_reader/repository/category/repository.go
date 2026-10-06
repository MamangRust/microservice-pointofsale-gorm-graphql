// Package category reads category price statistics from ClickHouse's
// order_item_daily table. It replaces the category slice of the old single
// stats-reader Repository god interface.
package category

import "context"

// Result models mirror the ClickHouse aggregations in clickhouse.go. Counts and
// amounts are int64/uint64: ClickHouse UInt64/Int64 columns scan into them
// directly (clickhouse-go v2 rejects implicit UInt64→int64 conversion).

type CategoryMonthTotalPrice struct {
	Year         string
	Month        string
	TotalRevenue int64
}

type CategoryYearTotalPrice struct {
	Year         string
	TotalRevenue int64
}

type CategoryMonthPrice struct {
	Month        string
	CategoryID   uint64
	OrderCount   uint64
	ItemsSold    uint64
	TotalRevenue int64
}

type CategoryYearPrice struct {
	Year               string
	CategoryID         uint64
	OrderCount         uint64
	ItemsSold          uint64
	TotalRevenue       int64
	UniqueProductsSold uint64
}

// Repository serves the category rpcs of CategoryStatsService,
// CategoryStatsByMerchantService and CategoryStatsByIdService. category_name is
// not stored in ClickHouse — the front end resolves it via the category domain
// service.
type Repository interface {
	GetMonthlyTotalPrices(ctx context.Context, year, month int) ([]CategoryMonthTotalPrice, error)
	GetYearlyTotalPrices(ctx context.Context, year int) ([]CategoryYearTotalPrice, error)
	GetMonthPrice(ctx context.Context, year int) ([]CategoryMonthPrice, error)
	GetYearPrice(ctx context.Context, year int) ([]CategoryYearPrice, error)
	GetMonthlyTotalPricesByMerchant(ctx context.Context, year, month, merchantID int) ([]CategoryMonthTotalPrice, error)
	GetYearlyTotalPricesByMerchant(ctx context.Context, year, merchantID int) ([]CategoryYearTotalPrice, error)
	GetMonthPriceByMerchant(ctx context.Context, year, merchantID int) ([]CategoryMonthPrice, error)
	GetYearPriceByMerchant(ctx context.Context, year, merchantID int) ([]CategoryYearPrice, error)
	GetMonthlyTotalPricesById(ctx context.Context, year, month, categoryID int) ([]CategoryMonthTotalPrice, error)
	GetYearlyTotalPricesById(ctx context.Context, year, categoryID int) ([]CategoryYearTotalPrice, error)
	GetMonthPriceById(ctx context.Context, year, categoryID int) ([]CategoryMonthPrice, error)
	GetYearPriceById(ctx context.Context, year, categoryID int) ([]CategoryYearPrice, error)
}
