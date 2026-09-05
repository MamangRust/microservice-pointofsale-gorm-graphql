package response

// ── Order stats (ClickHouse, F4) ─────────────────────────────────────────────

type OrderMonthlyRevenueResponse struct {
	Year         string `json:"year"`
	Month        string `json:"month"`
	TotalRevenue int64  `json:"total_revenue"`
	OrderCount   int    `json:"order_count"`
}

type OrderYearlyRevenueResponse struct {
	Year         string `json:"year"`
	TotalRevenue int64  `json:"total_revenue"`
	OrderCount   int    `json:"order_count"`
}

type CashierMonthlyRevenueResponse struct {
	Year         string `json:"year"`
	Month        string `json:"month"`
	CashierID    int    `json:"cashier_id"`
	TotalRevenue int64  `json:"total_revenue"`
	OrderCount   int    `json:"order_count"`
}

type ApiResponseOrderMonthlyRevenue struct {
	Status  string                        `json:"status"`
	Message string                        `json:"message"`
	Data    []*OrderMonthlyRevenueResponse `json:"data"`
}

type ApiResponseOrderYearlyRevenue struct {
	Status  string                       `json:"status"`
	Message string                       `json:"message"`
	Data    []*OrderYearlyRevenueResponse `json:"data"`
}

type ApiResponseCashierMonthlyRevenue struct {
	Status  string                          `json:"status"`
	Message string                          `json:"message"`
	Data    []*CashierMonthlyRevenueResponse `json:"data"`
}

// ── Product / category stats ────────────────────────────────────────────────

type ProductMonthlySoldResponse struct {
	Month     string `json:"month"`
	ProductID int    `json:"product_id"`
	Quantity  int64  `json:"quantity"`
	Subtotal  int64  `json:"subtotal"`
}

type CategoryMonthlySoldResponse struct {
	Month      string `json:"month"`
	CategoryID int    `json:"category_id"`
	Quantity   int64  `json:"quantity"`
	Subtotal   int64  `json:"subtotal"`
}

type ApiResponseProductMonthlySold struct {
	Status  string                       `json:"status"`
	Message string                       `json:"message"`
	Data    []*ProductMonthlySoldResponse `json:"data"`
}

type ApiResponseCategoryMonthlySold struct {
	Status  string                        `json:"status"`
	Message string                        `json:"message"`
	Data    []*CategoryMonthlySoldResponse `json:"data"`
}

// ── Transaction stats ───────────────────────────────────────────────────────

type TransactionMonthlySuccessResponse struct {
	Month       string `json:"month"`
	TotalCount  int64  `json:"total_count"`
	TotalAmount int64  `json:"total_amount"`
}

type ApiResponseTransactionMonthlySuccess struct {
	Status  string                             `json:"status"`
	Message string                             `json:"message"`
	Data    []*TransactionMonthlySuccessResponse `json:"data"`
}

// ── Cashier stats ───────────────────────────────────────────────────────────

type CashierMonthlyOrdersResponse struct {
	Month       string `json:"month"`
	CashierID   int    `json:"cashier_id"`
	OrderCount  int64  `json:"order_count"`
	TotalAmount int64  `json:"total_amount"`
}

type ApiResponseCashierMonthlyOrders struct {
	Status  string                         `json:"status"`
	Message string                         `json:"message"`
	Data    []*CashierMonthlyOrdersResponse `json:"data"`
}
