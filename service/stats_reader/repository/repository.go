package repository

import "context"

// Result models mirror the ClickHouse aggregations in clickhouse.go.

type MonthlyRevenue struct {
	Year         string
	Month        string
	TotalRevenue int64
	OrderCount   uint64
}

type YearlyRevenue struct {
	Year         string
	TotalRevenue int64
	OrderCount   uint64
}

type CashierMonthlyRevenue struct {
	Year         string
	Month        string
	CashierID    uint64
	TotalRevenue int64
	OrderCount   uint64
}

type ProductMonthlySold struct {
	Month     string
	ProductID uint64
	Quantity  uint64
	Subtotal  int64
}

type CategoryMonthlySold struct {
	Month      string
	CategoryID uint64
	Quantity   uint64
	Subtotal   int64
}

type TransactionMonthlySuccess struct {
	Month       string
	TotalCount  uint64
	TotalAmount int64
}

type CashierMonthlyOrders struct {
	Month       string
	CashierID   uint64
	OrderCount  uint64
	TotalAmount int64
}

// ── Category price stats ──────────────────────────────────────────────────

type CategoryMonthPrice struct {
	Month      string
	CategoryID uint64
	OrderCount uint64
	ItemsSold  uint64
	TotalRev   int64
}

type CategoryYearPrice struct {
	Year         string
	CategoryID   uint64
	OrderCount   uint64
	ItemsSold    uint64
	TotalRev     int64
	UniqueProd   uint64
}

type CategoryTotalPrice struct {
	Year       string
	Month      string
	TotalRev   int64
}

type CategoryYearTotalPrice struct {
	Year     string
	TotalRev int64
}

// ── Order detailed stats ──────────────────────────────────────────────────

type OrderMonthlyDetail struct {
	Month        string
	OrderCount   uint64
	TotalRev     int64
	ItemsSold    uint64
}

type OrderYearlyDetail struct {
	Year           string
	OrderCount     uint64
	TotalRev       int64
	ItemsSold      uint64
	ActiveCashiers uint64
	UniqueProds    uint64
}

type OrderTotalRevenue struct {
	Year       string
	Month      string
	OrderCount uint64
	TotalRev   int64
	ItemsSold  uint64
}

type OrderYearTotalRevenue struct {
	Year           string
	OrderCount     uint64
	TotalRev       int64
	ItemsSold      uint64
	ActiveCashiers uint64
	UniqueProds    uint64
}

// ── Transaction detailed stats ────────────────────────────────────────────

type TransactionStatusMonthly struct {
	Year       string
	Month      string
	Count      uint64
	Amount     int64
}

type TransactionStatusYearly struct {
	Year   string
	Count  uint64
	Amount int64
}

type TransactionMethodMonthly struct {
	Month         string
	PaymentMethod string
	Transactions  uint64
	Amount        int64
}

type TransactionMethodYearly struct {
	Year          string
	PaymentMethod string
	Transactions  uint64
	Amount        int64
}

// ── Cashier detailed stats ────────────────────────────────────────────────

type CashierSalesMonthly struct {
	Month      string
	CashierID  uint64
	OrderCount uint64
	TotalSales int64
}

type CashierSalesYearly struct {
	Year       string
	CashierID  uint64
	OrderCount uint64
	TotalSales int64
}

type CashierTotalSales struct {
	Year       string
	Month      string
	TotalSales int64
}

type CashierYearTotalSales struct {
	Year       string
	TotalSales int64
}

// Repository serves the stats gRPC contracts straight from ClickHouse.
type Repository interface {
	// ── Order stats ─────────────────────────────────────────────────────
	GetMonthlyTotalRevenue(ctx context.Context, year, month int) ([]MonthlyRevenue, error)
	GetYearlyTotalRevenue(ctx context.Context, year int) ([]YearlyRevenue, error)
	GetCashierMonthlyRevenue(ctx context.Context, cashierID int) ([]CashierMonthlyRevenue, error)

	// ── Order detailed ──────────────────────────────────────────────────
	GetOrderMonthlyDetail(ctx context.Context, year int) ([]OrderMonthlyDetail, error)
	GetOrderYearlyDetail(ctx context.Context, year int) ([]OrderYearlyDetail, error)
	GetOrderMonthlyDetailByMerchant(ctx context.Context, year, merchantID int) ([]OrderMonthlyDetail, error)
	GetOrderYearlyDetailByMerchant(ctx context.Context, year, merchantID int) ([]OrderYearlyDetail, error)
	GetOrderMonthlyTotalRevenue(ctx context.Context, year, month int) ([]OrderTotalRevenue, error)
	GetOrderYearlyTotalRevenue(ctx context.Context, year int) ([]OrderYearTotalRevenue, error)
	GetOrderMonthlyTotalRevenueById(ctx context.Context, year, month, orderID int) ([]OrderTotalRevenue, error)
	GetOrderYearlyTotalRevenueById(ctx context.Context, year, orderID int) ([]OrderYearTotalRevenue, error)
	GetOrderMonthlyTotalRevenueByMerchant(ctx context.Context, year, month, merchantID int) ([]OrderTotalRevenue, error)
	GetOrderYearlyTotalRevenueByMerchant(ctx context.Context, year, merchantID int) ([]OrderYearTotalRevenue, error)

	// ── Product / category stats ────────────────────────────────────────
	GetProductMonthlySold(ctx context.Context, year, month int) ([]ProductMonthlySold, error)
	GetCategoryMonthlySold(ctx context.Context, year, month int) ([]CategoryMonthlySold, error)

	// ── Category price ──────────────────────────────────────────────────
	GetCategoryMonthPrice(ctx context.Context, year int) ([]CategoryMonthPrice, error)
	GetCategoryYearPrice(ctx context.Context, year int) ([]CategoryYearPrice, error)
	GetCategoryMonthPriceByMerchant(ctx context.Context, year, merchantID int) ([]CategoryMonthPrice, error)
	GetCategoryYearPriceByMerchant(ctx context.Context, year, merchantID int) ([]CategoryYearPrice, error)
	GetCategoryMonthPriceById(ctx context.Context, year, categoryID int) ([]CategoryMonthPrice, error)
	GetCategoryYearPriceById(ctx context.Context, year, categoryID int) ([]CategoryYearPrice, error)
	GetCategoryMonthlyTotalPrices(ctx context.Context, year, month int) ([]CategoryTotalPrice, error)
	GetCategoryYearlyTotalPrices(ctx context.Context, year int) ([]CategoryYearTotalPrice, error)
	GetCategoryMonthlyTotalPricesById(ctx context.Context, year, month, categoryID int) ([]CategoryTotalPrice, error)
	GetCategoryYearlyTotalPricesById(ctx context.Context, year, categoryID int) ([]CategoryYearTotalPrice, error)
	GetCategoryMonthlyTotalPricesByMerchant(ctx context.Context, year, month, merchantID int) ([]CategoryTotalPrice, error)
	GetCategoryYearlyTotalPricesByMerchant(ctx context.Context, year, merchantID int) ([]CategoryYearTotalPrice, error)

	// ── Transaction stats ───────────────────────────────────────────────
	GetTransactionMonthlySuccess(ctx context.Context, year, month int) ([]TransactionMonthlySuccess, error)
	GetTransactionStatusMonthly(ctx context.Context, year, month int, status string) ([]TransactionStatusMonthly, error)
	GetTransactionStatusYearly(ctx context.Context, year int, status string) ([]TransactionStatusYearly, error)
	GetTransactionStatusMonthlyByMerchant(ctx context.Context, year, month, merchantID int, status string) ([]TransactionStatusMonthly, error)
	GetTransactionStatusYearlyByMerchant(ctx context.Context, year, merchantID int, status string) ([]TransactionStatusYearly, error)
	GetTransactionMethodMonthly(ctx context.Context, year, month int, status string) ([]TransactionMethodMonthly, error)
	GetTransactionMethodYearly(ctx context.Context, year int, status string) ([]TransactionMethodYearly, error)
	GetTransactionMethodMonthlyByMerchant(ctx context.Context, year, month, merchantID int, status string) ([]TransactionMethodMonthly, error)
	GetTransactionMethodYearlyByMerchant(ctx context.Context, year, merchantID int, status string) ([]TransactionMethodYearly, error)

	// ── Cashier stats ───────────────────────────────────────────────────
	GetCashierMonthlyOrders(ctx context.Context, cashierID int) ([]CashierMonthlyOrders, error)
	GetCashierMonthlyTotalSales(ctx context.Context, year, month int) ([]CashierTotalSales, error)
	GetCashierYearlyTotalSales(ctx context.Context, year int) ([]CashierYearTotalSales, error)
	GetCashierMonthlyTotalSalesById(ctx context.Context, year, month, cashierID int) ([]CashierTotalSales, error)
	GetCashierYearlyTotalSalesById(ctx context.Context, year, cashierID int) ([]CashierYearTotalSales, error)
	GetCashierMonthlyTotalSalesByMerchant(ctx context.Context, year, month, merchantID int) ([]CashierTotalSales, error)
	GetCashierYearlyTotalSalesByMerchant(ctx context.Context, year, merchantID int) ([]CashierYearTotalSales, error)
	GetCashierSalesMonthly(ctx context.Context, year int) ([]CashierSalesMonthly, error)
	GetCashierSalesYearly(ctx context.Context, year int) ([]CashierSalesYearly, error)
	GetCashierSalesMonthlyByMerchant(ctx context.Context, year, merchantID int) ([]CashierSalesMonthly, error)
	GetCashierSalesYearlyByMerchant(ctx context.Context, year, merchantID int) ([]CashierSalesYearly, error)
	GetCashierSalesMonthlyById(ctx context.Context, year, cashierID int) ([]CashierSalesMonthly, error)
	GetCashierSalesYearlyById(ctx context.Context, year, cashierID int) ([]CashierSalesYearly, error)
}
