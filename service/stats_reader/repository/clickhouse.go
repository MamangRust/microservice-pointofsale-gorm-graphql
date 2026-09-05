package repository

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type ClickHouseReaderRepository struct {
	conn clickhouse.Conn
}

func NewClickHouseReaderRepository(conn clickhouse.Conn) *ClickHouseReaderRepository {
	return &ClickHouseReaderRepository{conn: conn}
}

// ─── Order stats ─────────────────────────────────────────────────────────

func (r *ClickHouseReaderRepository) GetMonthlyTotalRevenue(ctx context.Context, year, month int) ([]MonthlyRevenue, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       formatDateTime(event_time, '%b') AS month,
		       sum(total_price) AS total_revenue,
		       count() AS order_count
		FROM order_daily FINAL
		WHERE toYear(event_time) = ? AND toMonth(event_time) = ?
		GROUP BY year, month, toMonth(event_time)
		ORDER BY year, toMonth(event_time)
	`
	rows, err := r.conn.Query(ctx, query, year, month)
	if err != nil {
		return nil, fmt.Errorf("query monthly total revenue: %w", err)
	}
	defer rows.Close()

	var results []MonthlyRevenue
	for rows.Next() {
		var m MonthlyRevenue
		if err := rows.Scan(&m.Year, &m.Month, &m.TotalRevenue, &m.OrderCount); err != nil {
			return nil, fmt.Errorf("scan monthly total revenue: %w", err)
		}
		results = append(results, m)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetYearlyTotalRevenue(ctx context.Context, year int) ([]YearlyRevenue, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       sum(total_price) AS total_revenue,
		       count() AS order_count
		FROM order_daily FINAL
		WHERE toYear(event_time) IN (?, ?)
		GROUP BY year
		ORDER BY year DESC
	`
	rows, err := r.conn.Query(ctx, query, year, year-1)
	if err != nil {
		return nil, fmt.Errorf("query yearly total revenue: %w", err)
	}
	defer rows.Close()

	var results []YearlyRevenue
	for rows.Next() {
		var y YearlyRevenue
		if err := rows.Scan(&y.Year, &y.TotalRevenue, &y.OrderCount); err != nil {
			return nil, fmt.Errorf("scan yearly total revenue: %w", err)
		}
		results = append(results, y)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetCashierMonthlyRevenue(ctx context.Context, cashierID int) ([]CashierMonthlyRevenue, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       formatDateTime(event_time, '%b') AS month,
		       cashier_id,
		       sum(total_price) AS total_revenue,
		       count() AS order_count
		FROM order_daily FINAL
		WHERE cashier_id = ?
		GROUP BY year, month, toMonth(event_time), cashier_id
		ORDER BY toMonth(event_time)
	`
	rows, err := r.conn.Query(ctx, query, cashierID)
	if err != nil {
		return nil, fmt.Errorf("query cashier monthly revenue: %w", err)
	}
	defer rows.Close()

	var results []CashierMonthlyRevenue
	for rows.Next() {
		var c CashierMonthlyRevenue
		if err := rows.Scan(&c.Year, &c.Month, &c.CashierID, &c.TotalRevenue, &c.OrderCount); err != nil {
			return nil, fmt.Errorf("scan cashier monthly revenue: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

// ─── Order detailed stats ────────────────────────────────────────────────

func (r *ClickHouseReaderRepository) GetOrderMonthlyDetail(ctx context.Context, year int) ([]OrderMonthlyDetail, error) {
	query := `
		SELECT formatDateTime(o.event_time, '%b') AS month,
		       count(DISTINCT o.order_id) AS order_count,
		       sum(o.total_price) AS total_rev,
		       sum(oi.quantity) AS items_sold
		FROM (SELECT * FROM order_daily FINAL) o
		LEFT JOIN (SELECT * FROM order_item_daily FINAL) oi
		  ON oi.order_id = o.order_id AND toDate(oi.event_time) = toDate(o.event_time)
		WHERE toYear(o.event_time) = ?
		GROUP BY month, toMonth(o.event_time)
		ORDER BY toMonth(o.event_time)
	`
	rows, err := r.conn.Query(ctx, query, year)
	if err != nil {
		return nil, fmt.Errorf("query order monthly detail: %w", err)
	}
	defer rows.Close()

	var results []OrderMonthlyDetail
	for rows.Next() {
		var d OrderMonthlyDetail
		if err := rows.Scan(&d.Month, &d.OrderCount, &d.TotalRev, &d.ItemsSold); err != nil {
			return nil, fmt.Errorf("scan order monthly detail: %w", err)
		}
		results = append(results, d)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetOrderYearlyDetail(ctx context.Context, year int) ([]OrderYearlyDetail, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       count() AS order_count,
		       sum(total_price) AS total_rev,
		       uniqExact(cashier_id) AS active_cashiers
		FROM order_daily FINAL
		WHERE toYear(event_time) = ?
		GROUP BY year
	`
	rows, err := r.conn.Query(ctx, query, year)
	if err != nil {
		return nil, fmt.Errorf("query order yearly detail: %w", err)
	}
	defer rows.Close()

	var results []OrderYearlyDetail
	for rows.Next() {
		var d OrderYearlyDetail
		if err := rows.Scan(&d.Year, &d.OrderCount, &d.TotalRev, &d.ActiveCashiers); err != nil {
			return nil, fmt.Errorf("scan order yearly detail: %w", err)
		}
		results = append(results, d)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetOrderMonthlyDetailByMerchant(ctx context.Context, year, merchantID int) ([]OrderMonthlyDetail, error) {
	query := `
		SELECT formatDateTime(event_time, '%b') AS month,
		       count() AS order_count,
		       sum(total_price) AS total_rev
		FROM order_daily FINAL
		WHERE toYear(event_time) = ? AND merchant_id = ?
		GROUP BY month, toMonth(event_time)
		ORDER BY toMonth(event_time)
	`
	rows, err := r.conn.Query(ctx, query, year, merchantID)
	if err != nil {
		return nil, fmt.Errorf("query order monthly detail by merchant: %w", err)
	}
	defer rows.Close()

	var results []OrderMonthlyDetail
	for rows.Next() {
		var d OrderMonthlyDetail
		if err := rows.Scan(&d.Month, &d.OrderCount, &d.TotalRev); err != nil {
			return nil, fmt.Errorf("scan order monthly detail by merchant: %w", err)
		}
		results = append(results, d)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetOrderYearlyDetailByMerchant(ctx context.Context, year, merchantID int) ([]OrderYearlyDetail, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       count() AS order_count,
		       sum(total_price) AS total_rev
		FROM order_daily FINAL
		WHERE toYear(event_time) = ? AND merchant_id = ?
		GROUP BY year
	`
	rows, err := r.conn.Query(ctx, query, year, merchantID)
	if err != nil {
		return nil, fmt.Errorf("query order yearly detail by merchant: %w", err)
	}
	defer rows.Close()

	var results []OrderYearlyDetail
	for rows.Next() {
		var d OrderYearlyDetail
		if err := rows.Scan(&d.Year, &d.OrderCount, &d.TotalRev); err != nil {
			return nil, fmt.Errorf("scan order yearly detail by merchant: %w", err)
		}
		results = append(results, d)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetOrderMonthlyTotalRevenue(ctx context.Context, year, month int) ([]OrderTotalRevenue, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       formatDateTime(event_time, '%b') AS month,
		       count() AS order_count,
		       sum(total_price) AS total_rev
		FROM order_daily FINAL
		WHERE toYear(event_time) = ? AND toMonth(event_time) = ?
		GROUP BY year, month, toMonth(event_time)
	`
	rows, err := r.conn.Query(ctx, query, year, month)
	if err != nil {
		return nil, fmt.Errorf("query order monthly total revenue: %w", err)
	}
	defer rows.Close()

	var results []OrderTotalRevenue
	for rows.Next() {
		var d OrderTotalRevenue
		if err := rows.Scan(&d.Year, &d.Month, &d.OrderCount, &d.TotalRev); err != nil {
			return nil, fmt.Errorf("scan order monthly total revenue: %w", err)
		}
		results = append(results, d)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetOrderYearlyTotalRevenue(ctx context.Context, year int) ([]OrderYearTotalRevenue, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       count() AS order_count,
		       sum(total_price) AS total_rev,
		       uniqExact(cashier_id) AS active_cashiers
		FROM order_daily FINAL
		WHERE toYear(event_time) = ?
		GROUP BY year
	`
	rows, err := r.conn.Query(ctx, query, year)
	if err != nil {
		return nil, fmt.Errorf("query order yearly total revenue: %w", err)
	}
	defer rows.Close()

	var results []OrderYearTotalRevenue
	for rows.Next() {
		var d OrderYearTotalRevenue
		if err := rows.Scan(&d.Year, &d.OrderCount, &d.TotalRev, &d.ActiveCashiers); err != nil {
			return nil, fmt.Errorf("scan order yearly total revenue: %w", err)
		}
		results = append(results, d)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetOrderMonthlyTotalRevenueById(ctx context.Context, year, month, orderID int) ([]OrderTotalRevenue, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       formatDateTime(event_time, '%b') AS month,
		       count() AS order_count,
		       sum(total_price) AS total_rev
		FROM order_daily FINAL
		WHERE toYear(event_time) = ? AND toMonth(event_time) = ? AND order_id = ?
		GROUP BY year, month, toMonth(event_time)
	`
	rows, err := r.conn.Query(ctx, query, year, month, orderID)
	if err != nil {
		return nil, fmt.Errorf("query order monthly total revenue by id: %w", err)
	}
	defer rows.Close()

	var results []OrderTotalRevenue
	for rows.Next() {
		var d OrderTotalRevenue
		if err := rows.Scan(&d.Year, &d.Month, &d.OrderCount, &d.TotalRev); err != nil {
			return nil, fmt.Errorf("scan order monthly total revenue by id: %w", err)
		}
		results = append(results, d)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetOrderYearlyTotalRevenueById(ctx context.Context, year, orderID int) ([]OrderYearTotalRevenue, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       count() AS order_count,
		       sum(total_price) AS total_rev
		FROM order_daily FINAL
		WHERE toYear(event_time) = ? AND order_id = ?
		GROUP BY year
	`
	rows, err := r.conn.Query(ctx, query, year, orderID)
	if err != nil {
		return nil, fmt.Errorf("query order yearly total revenue by id: %w", err)
	}
	defer rows.Close()

	var results []OrderYearTotalRevenue
	for rows.Next() {
		var d OrderYearTotalRevenue
		if err := rows.Scan(&d.Year, &d.OrderCount, &d.TotalRev); err != nil {
			return nil, fmt.Errorf("scan order yearly total revenue by id: %w", err)
		}
		results = append(results, d)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetOrderMonthlyTotalRevenueByMerchant(ctx context.Context, year, month, merchantID int) ([]OrderTotalRevenue, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       formatDateTime(event_time, '%b') AS month,
		       count() AS order_count,
		       sum(total_price) AS total_rev
		FROM order_daily FINAL
		WHERE toYear(event_time) = ? AND toMonth(event_time) = ? AND merchant_id = ?
		GROUP BY year, month, toMonth(event_time)
	`
	rows, err := r.conn.Query(ctx, query, year, month, merchantID)
	if err != nil {
		return nil, fmt.Errorf("query order monthly total revenue by merchant: %w", err)
	}
	defer rows.Close()

	var results []OrderTotalRevenue
	for rows.Next() {
		var d OrderTotalRevenue
		if err := rows.Scan(&d.Year, &d.Month, &d.OrderCount, &d.TotalRev); err != nil {
			return nil, fmt.Errorf("scan order monthly total revenue by merchant: %w", err)
		}
		results = append(results, d)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetOrderYearlyTotalRevenueByMerchant(ctx context.Context, year, merchantID int) ([]OrderYearTotalRevenue, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       count() AS order_count,
		       sum(total_price) AS total_rev
		FROM order_daily FINAL
		WHERE toYear(event_time) = ? AND merchant_id = ?
		GROUP BY year
	`
	rows, err := r.conn.Query(ctx, query, year, merchantID)
	if err != nil {
		return nil, fmt.Errorf("query order yearly total revenue by merchant: %w", err)
	}
	defer rows.Close()

	var results []OrderYearTotalRevenue
	for rows.Next() {
		var d OrderYearTotalRevenue
		if err := rows.Scan(&d.Year, &d.OrderCount, &d.TotalRev); err != nil {
			return nil, fmt.Errorf("scan order yearly total revenue by merchant: %w", err)
		}
		results = append(results, d)
	}
	return results, rows.Err()
}

// ─── Product / category stats ────────────────────────────────────────────

func (r *ClickHouseReaderRepository) GetProductMonthlySold(ctx context.Context, year, month int) ([]ProductMonthlySold, error) {
	query := `
		SELECT formatDateTime(event_time, '%b') AS month,
		       product_id,
		       sum(quantity) AS quantity,
		       sum(subtotal) AS subtotal
		FROM order_item_daily FINAL
		WHERE toYear(event_time) = ? AND toMonth(event_time) = ?
		GROUP BY month, toMonth(event_time), product_id
		ORDER BY toMonth(event_time), subtotal DESC
	`
	rows, err := r.conn.Query(ctx, query, year, month)
	if err != nil {
		return nil, fmt.Errorf("query product monthly sold: %w", err)
	}
	defer rows.Close()

	var results []ProductMonthlySold
	for rows.Next() {
		var p ProductMonthlySold
		if err := rows.Scan(&p.Month, &p.ProductID, &p.Quantity, &p.Subtotal); err != nil {
			return nil, fmt.Errorf("scan product monthly sold: %w", err)
		}
		results = append(results, p)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetCategoryMonthlySold(ctx context.Context, year, month int) ([]CategoryMonthlySold, error) {
	query := `
		SELECT formatDateTime(event_time, '%b') AS month,
		       category_id,
		       sum(quantity) AS quantity,
		       sum(subtotal) AS subtotal
		FROM order_item_daily FINAL
		WHERE toYear(event_time) = ? AND toMonth(event_time) = ?
		GROUP BY month, toMonth(event_time), category_id
		ORDER BY toMonth(event_time), subtotal DESC
	`
	rows, err := r.conn.Query(ctx, query, year, month)
	if err != nil {
		return nil, fmt.Errorf("query category monthly sold: %w", err)
	}
	defer rows.Close()

	var results []CategoryMonthlySold
	for rows.Next() {
		var c CategoryMonthlySold
		if err := rows.Scan(&c.Month, &c.CategoryID, &c.Quantity, &c.Subtotal); err != nil {
			return nil, fmt.Errorf("scan category monthly sold: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

// ─── Category price stats ────────────────────────────────────────────────

func (r *ClickHouseReaderRepository) GetCategoryMonthPrice(ctx context.Context, year int) ([]CategoryMonthPrice, error) {
	query := `
		SELECT formatDateTime(oi.event_time, '%b') AS month,
		       oi.category_id,
		       count(DISTINCT oi.order_id) AS order_count,
		       sum(oi.quantity) AS items_sold,
		       sum(oi.subtotal) AS total_rev
		FROM (SELECT * FROM order_item_daily FINAL) oi
		WHERE toYear(oi.event_time) = ?
		GROUP BY month, toMonth(oi.event_time), oi.category_id
		ORDER BY toMonth(oi.event_time)
	`
	rows, err := r.conn.Query(ctx, query, year)
	if err != nil {
		return nil, fmt.Errorf("query category month price: %w", err)
	}
	defer rows.Close()

	var results []CategoryMonthPrice
	for rows.Next() {
		var c CategoryMonthPrice
		if err := rows.Scan(&c.Month, &c.CategoryID, &c.OrderCount, &c.ItemsSold, &c.TotalRev); err != nil {
			return nil, fmt.Errorf("scan category month price: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetCategoryYearPrice(ctx context.Context, year int) ([]CategoryYearPrice, error) {
	query := `
		SELECT toString(toYear(oi.event_time)) AS year,
		       oi.category_id,
		       count(DISTINCT oi.order_id) AS order_count,
		       sum(oi.quantity) AS items_sold,
		       sum(oi.subtotal) AS total_rev,
		       uniqExact(oi.product_id) AS unique_prod
		FROM (SELECT * FROM order_item_daily FINAL) oi
		WHERE toYear(oi.event_time) = ?
		GROUP BY year, oi.category_id
	`
	rows, err := r.conn.Query(ctx, query, year)
	if err != nil {
		return nil, fmt.Errorf("query category year price: %w", err)
	}
	defer rows.Close()

	var results []CategoryYearPrice
	for rows.Next() {
		var c CategoryYearPrice
		if err := rows.Scan(&c.Year, &c.CategoryID, &c.OrderCount, &c.ItemsSold, &c.TotalRev, &c.UniqueProd); err != nil {
			return nil, fmt.Errorf("scan category year price: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetCategoryMonthPriceByMerchant(ctx context.Context, year, merchantID int) ([]CategoryMonthPrice, error) {
	query := `
		SELECT formatDateTime(oi.event_time, '%b') AS month,
		       oi.category_id,
		       count(DISTINCT oi.order_id) AS order_count,
		       sum(oi.quantity) AS items_sold,
		       sum(oi.subtotal) AS total_rev
		FROM (SELECT * FROM order_item_daily FINAL) oi
		INNER JOIN order_daily o ON oi.order_id = o.order_id AND toDate(oi.event_time) = toDate(o.event_time)
		WHERE toYear(oi.event_time) = ? AND o.merchant_id = ?
		GROUP BY month, toMonth(oi.event_time), oi.category_id
		ORDER BY toMonth(oi.event_time)
	`
	rows, err := r.conn.Query(ctx, query, year, merchantID)
	if err != nil {
		return nil, fmt.Errorf("query category month price by merchant: %w", err)
	}
	defer rows.Close()

	var results []CategoryMonthPrice
	for rows.Next() {
		var c CategoryMonthPrice
		if err := rows.Scan(&c.Month, &c.CategoryID, &c.OrderCount, &c.ItemsSold, &c.TotalRev); err != nil {
			return nil, fmt.Errorf("scan category month price by merchant: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetCategoryYearPriceByMerchant(ctx context.Context, year, merchantID int) ([]CategoryYearPrice, error) {
	query := `
		SELECT toString(toYear(oi.event_time)) AS year,
		       oi.category_id,
		       count(DISTINCT oi.order_id) AS order_count,
		       sum(oi.quantity) AS items_sold,
		       sum(oi.subtotal) AS total_rev,
		       uniqExact(oi.product_id) AS unique_prod
		FROM (SELECT * FROM order_item_daily FINAL) oi
		INNER JOIN order_daily o ON oi.order_id = o.order_id AND toDate(oi.event_time) = toDate(o.event_time)
		WHERE toYear(oi.event_time) = ? AND o.merchant_id = ?
		GROUP BY year, oi.category_id
	`
	rows, err := r.conn.Query(ctx, query, year, merchantID)
	if err != nil {
		return nil, fmt.Errorf("query category year price by merchant: %w", err)
	}
	defer rows.Close()

	var results []CategoryYearPrice
	for rows.Next() {
		var c CategoryYearPrice
		if err := rows.Scan(&c.Year, &c.CategoryID, &c.OrderCount, &c.ItemsSold, &c.TotalRev, &c.UniqueProd); err != nil {
			return nil, fmt.Errorf("scan category year price by merchant: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetCategoryMonthPriceById(ctx context.Context, year, categoryID int) ([]CategoryMonthPrice, error) {
	query := `
		SELECT formatDateTime(oi.event_time, '%b') AS month,
		       oi.category_id,
		       count(DISTINCT oi.order_id) AS order_count,
		       sum(oi.quantity) AS items_sold,
		       sum(oi.subtotal) AS total_rev
		FROM (SELECT * FROM order_item_daily FINAL) oi
		WHERE toYear(oi.event_time) = ? AND oi.category_id = ?
		GROUP BY month, toMonth(oi.event_time), oi.category_id
		ORDER BY toMonth(oi.event_time)
	`
	rows, err := r.conn.Query(ctx, query, year, categoryID)
	if err != nil {
		return nil, fmt.Errorf("query category month price by id: %w", err)
	}
	defer rows.Close()

	var results []CategoryMonthPrice
	for rows.Next() {
		var c CategoryMonthPrice
		if err := rows.Scan(&c.Month, &c.CategoryID, &c.OrderCount, &c.ItemsSold, &c.TotalRev); err != nil {
			return nil, fmt.Errorf("scan category month price by id: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetCategoryYearPriceById(ctx context.Context, year, categoryID int) ([]CategoryYearPrice, error) {
	query := `
		SELECT toString(toYear(oi.event_time)) AS year,
		       oi.category_id,
		       count(DISTINCT oi.order_id) AS order_count,
		       sum(oi.quantity) AS items_sold,
		       sum(oi.subtotal) AS total_rev,
		       uniqExact(oi.product_id) AS unique_prod
		FROM (SELECT * FROM order_item_daily FINAL) oi
		WHERE toYear(oi.event_time) = ? AND oi.category_id = ?
		GROUP BY year, oi.category_id
	`
	rows, err := r.conn.Query(ctx, query, year, categoryID)
	if err != nil {
		return nil, fmt.Errorf("query category year price by id: %w", err)
	}
	defer rows.Close()

	var results []CategoryYearPrice
	for rows.Next() {
		var c CategoryYearPrice
		if err := rows.Scan(&c.Year, &c.CategoryID, &c.OrderCount, &c.ItemsSold, &c.TotalRev, &c.UniqueProd); err != nil {
			return nil, fmt.Errorf("scan category year price by id: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetCategoryMonthlyTotalPrices(ctx context.Context, year, month int) ([]CategoryTotalPrice, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       formatDateTime(event_time, '%b') AS month,
		       sum(subtotal) AS total_rev
		FROM order_item_daily FINAL
		WHERE toYear(event_time) = ? AND toMonth(event_time) = ?
		GROUP BY year, month, toMonth(event_time)
	`
	rows, err := r.conn.Query(ctx, query, year, month)
	if err != nil {
		return nil, fmt.Errorf("query category monthly total prices: %w", err)
	}
	defer rows.Close()

	var results []CategoryTotalPrice
	for rows.Next() {
		var c CategoryTotalPrice
		if err := rows.Scan(&c.Year, &c.Month, &c.TotalRev); err != nil {
			return nil, fmt.Errorf("scan category monthly total prices: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetCategoryYearlyTotalPrices(ctx context.Context, year int) ([]CategoryYearTotalPrice, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       sum(subtotal) AS total_rev
		FROM order_item_daily FINAL
		WHERE toYear(event_time) = ?
		GROUP BY year
	`
	rows, err := r.conn.Query(ctx, query, year)
	if err != nil {
		return nil, fmt.Errorf("query category yearly total prices: %w", err)
	}
	defer rows.Close()

	var results []CategoryYearTotalPrice
	for rows.Next() {
		var c CategoryYearTotalPrice
		if err := rows.Scan(&c.Year, &c.TotalRev); err != nil {
			return nil, fmt.Errorf("scan category yearly total prices: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetCategoryMonthlyTotalPricesById(ctx context.Context, year, month, categoryID int) ([]CategoryTotalPrice, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       formatDateTime(event_time, '%b') AS month,
		       sum(subtotal) AS total_rev
		FROM order_item_daily FINAL
		WHERE toYear(event_time) = ? AND toMonth(event_time) = ? AND category_id = ?
		GROUP BY year, month, toMonth(event_time)
	`
	rows, err := r.conn.Query(ctx, query, year, month, categoryID)
	if err != nil {
		return nil, fmt.Errorf("query category monthly total prices by id: %w", err)
	}
	defer rows.Close()

	var results []CategoryTotalPrice
	for rows.Next() {
		var c CategoryTotalPrice
		if err := rows.Scan(&c.Year, &c.Month, &c.TotalRev); err != nil {
			return nil, fmt.Errorf("scan category monthly total prices by id: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetCategoryYearlyTotalPricesById(ctx context.Context, year, categoryID int) ([]CategoryYearTotalPrice, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       sum(subtotal) AS total_rev
		FROM order_item_daily FINAL
		WHERE toYear(event_time) = ? AND category_id = ?
		GROUP BY year
	`
	rows, err := r.conn.Query(ctx, query, year, categoryID)
	if err != nil {
		return nil, fmt.Errorf("query category yearly total prices by id: %w", err)
	}
	defer rows.Close()

	var results []CategoryYearTotalPrice
	for rows.Next() {
		var c CategoryYearTotalPrice
		if err := rows.Scan(&c.Year, &c.TotalRev); err != nil {
			return nil, fmt.Errorf("scan category yearly total prices by id: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetCategoryMonthlyTotalPricesByMerchant(ctx context.Context, year, month, merchantID int) ([]CategoryTotalPrice, error) {
	query := `
		SELECT toString(toYear(oi.event_time)) AS year,
		       formatDateTime(oi.event_time, '%b') AS month,
		       sum(oi.subtotal) AS total_rev
		FROM (SELECT * FROM order_item_daily FINAL) oi
		INNER JOIN order_daily o ON oi.order_id = o.order_id AND toDate(oi.event_time) = toDate(o.event_time)
		WHERE toYear(oi.event_time) = ? AND toMonth(oi.event_time) = ? AND o.merchant_id = ?
		GROUP BY year, month, toMonth(oi.event_time)
	`
	rows, err := r.conn.Query(ctx, query, year, month, merchantID)
	if err != nil {
		return nil, fmt.Errorf("query category monthly total prices by merchant: %w", err)
	}
	defer rows.Close()

	var results []CategoryTotalPrice
	for rows.Next() {
		var c CategoryTotalPrice
		if err := rows.Scan(&c.Year, &c.Month, &c.TotalRev); err != nil {
			return nil, fmt.Errorf("scan category monthly total prices by merchant: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetCategoryYearlyTotalPricesByMerchant(ctx context.Context, year, merchantID int) ([]CategoryYearTotalPrice, error) {
	query := `
		SELECT toString(toYear(oi.event_time)) AS year,
		       sum(oi.subtotal) AS total_rev
		FROM (SELECT * FROM order_item_daily FINAL) oi
		INNER JOIN order_daily o ON oi.order_id = o.order_id AND toDate(oi.event_time) = toDate(o.event_time)
		WHERE toYear(oi.event_time) = ? AND o.merchant_id = ?
		GROUP BY year
	`
	rows, err := r.conn.Query(ctx, query, year, merchantID)
	if err != nil {
		return nil, fmt.Errorf("query category yearly total prices by merchant: %w", err)
	}
	defer rows.Close()

	var results []CategoryYearTotalPrice
	for rows.Next() {
		var c CategoryYearTotalPrice
		if err := rows.Scan(&c.Year, &c.TotalRev); err != nil {
			return nil, fmt.Errorf("scan category yearly total prices by merchant: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

// ─── Transaction stats ───────────────────────────────────────────────────

func (r *ClickHouseReaderRepository) GetTransactionMonthlySuccess(ctx context.Context, year, month int) ([]TransactionMonthlySuccess, error) {
	query := `
		SELECT formatDateTime(event_time, '%b') AS month,
		       countIf(lower(status) IN ('success', 'completed')) AS total_count,
		       sumIf(amount, lower(status) IN ('success', 'completed')) AS total_amount
		FROM transaction_daily FINAL
		WHERE toYear(event_time) = ? AND toMonth(event_time) = ?
		GROUP BY month, toMonth(event_time)
		ORDER BY toMonth(event_time)
	`
	rows, err := r.conn.Query(ctx, query, year, month)
	if err != nil {
		return nil, fmt.Errorf("query transaction monthly success: %w", err)
	}
	defer rows.Close()

	var results []TransactionMonthlySuccess
	for rows.Next() {
		var t TransactionMonthlySuccess
		if err := rows.Scan(&t.Month, &t.TotalCount, &t.TotalAmount); err != nil {
			return nil, fmt.Errorf("scan transaction monthly success: %w", err)
		}
		results = append(results, t)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetTransactionStatusMonthly(ctx context.Context, year, month int, status string) ([]TransactionStatusMonthly, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       formatDateTime(event_time, '%b') AS month,
		       countIf(lower(status) = lower(?)) AS cnt,
		       sumIf(amount, lower(status) = lower(?)) AS amt
		FROM transaction_daily FINAL
		WHERE toYear(event_time) = ? AND toMonth(event_time) = ?
		GROUP BY year, month, toMonth(event_time)
	`
	rows, err := r.conn.Query(ctx, query, status, status, year, month)
	if err != nil {
		return nil, fmt.Errorf("query transaction status monthly: %w", err)
	}
	defer rows.Close()

	var results []TransactionStatusMonthly
	for rows.Next() {
		var t TransactionStatusMonthly
		if err := rows.Scan(&t.Year, &t.Month, &t.Count, &t.Amount); err != nil {
			return nil, fmt.Errorf("scan transaction status monthly: %w", err)
		}
		results = append(results, t)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetTransactionStatusYearly(ctx context.Context, year int, status string) ([]TransactionStatusYearly, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       countIf(lower(status) = lower(?)) AS cnt,
		       sumIf(amount, lower(status) = lower(?)) AS amt
		FROM transaction_daily FINAL
		WHERE toYear(event_time) = ?
		GROUP BY year
	`
	rows, err := r.conn.Query(ctx, query, status, status, year)
	if err != nil {
		return nil, fmt.Errorf("query transaction status yearly: %w", err)
	}
	defer rows.Close()

	var results []TransactionStatusYearly
	for rows.Next() {
		var t TransactionStatusYearly
		if err := rows.Scan(&t.Year, &t.Count, &t.Amount); err != nil {
			return nil, fmt.Errorf("scan transaction status yearly: %w", err)
		}
		results = append(results, t)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetTransactionStatusMonthlyByMerchant(ctx context.Context, year, month, merchantID int, status string) ([]TransactionStatusMonthly, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       formatDateTime(event_time, '%b') AS month,
		       countIf(lower(status) = lower(?)) AS cnt,
		       sumIf(amount, lower(status) = lower(?)) AS amt
		FROM transaction_daily FINAL
		WHERE toYear(event_time) = ? AND toMonth(event_time) = ? AND merchant_id = ?
		GROUP BY year, month, toMonth(event_time)
	`
	rows, err := r.conn.Query(ctx, query, status, status, year, month, merchantID)
	if err != nil {
		return nil, fmt.Errorf("query transaction status monthly by merchant: %w", err)
	}
	defer rows.Close()

	var results []TransactionStatusMonthly
	for rows.Next() {
		var t TransactionStatusMonthly
		if err := rows.Scan(&t.Year, &t.Month, &t.Count, &t.Amount); err != nil {
			return nil, fmt.Errorf("scan transaction status monthly by merchant: %w", err)
		}
		results = append(results, t)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetTransactionStatusYearlyByMerchant(ctx context.Context, year, merchantID int, status string) ([]TransactionStatusYearly, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       countIf(lower(status) = lower(?)) AS cnt,
		       sumIf(amount, lower(status) = lower(?)) AS amt
		FROM transaction_daily FINAL
		WHERE toYear(event_time) = ? AND merchant_id = ?
		GROUP BY year
	`
	rows, err := r.conn.Query(ctx, query, status, status, year, merchantID)
	if err != nil {
		return nil, fmt.Errorf("query transaction status yearly by merchant: %w", err)
	}
	defer rows.Close()

	var results []TransactionStatusYearly
	for rows.Next() {
		var t TransactionStatusYearly
		if err := rows.Scan(&t.Year, &t.Count, &t.Amount); err != nil {
			return nil, fmt.Errorf("scan transaction status yearly by merchant: %w", err)
		}
		results = append(results, t)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetTransactionMethodMonthly(ctx context.Context, year, month int, status string) ([]TransactionMethodMonthly, error) {
	query := `
		SELECT formatDateTime(event_time, '%b') AS month,
		       payment_method,
		       countIf(lower(status) = lower(?)) AS transactions,
		       sumIf(amount, lower(status) = lower(?)) AS amt
		FROM transaction_daily FINAL
		WHERE toYear(event_time) = ? AND toMonth(event_time) = ?
		GROUP BY month, toMonth(event_time), payment_method
	`
	rows, err := r.conn.Query(ctx, query, status, status, year, month)
	if err != nil {
		return nil, fmt.Errorf("query transaction method monthly: %w", err)
	}
	defer rows.Close()

	var results []TransactionMethodMonthly
	for rows.Next() {
		var t TransactionMethodMonthly
		if err := rows.Scan(&t.Month, &t.PaymentMethod, &t.Transactions, &t.Amount); err != nil {
			return nil, fmt.Errorf("scan transaction method monthly: %w", err)
		}
		results = append(results, t)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetTransactionMethodYearly(ctx context.Context, year int, status string) ([]TransactionMethodYearly, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       payment_method,
		       countIf(lower(status) = lower(?)) AS transactions,
		       sumIf(amount, lower(status) = lower(?)) AS amt
		FROM transaction_daily FINAL
		WHERE toYear(event_time) = ?
		GROUP BY year, payment_method
	`
	rows, err := r.conn.Query(ctx, query, status, status, year)
	if err != nil {
		return nil, fmt.Errorf("query transaction method yearly: %w", err)
	}
	defer rows.Close()

	var results []TransactionMethodYearly
	for rows.Next() {
		var t TransactionMethodYearly
		if err := rows.Scan(&t.Year, &t.PaymentMethod, &t.Transactions, &t.Amount); err != nil {
			return nil, fmt.Errorf("scan transaction method yearly: %w", err)
		}
		results = append(results, t)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetTransactionMethodMonthlyByMerchant(ctx context.Context, year, month, merchantID int, status string) ([]TransactionMethodMonthly, error) {
	query := `
		SELECT formatDateTime(event_time, '%b') AS month,
		       payment_method,
		       countIf(lower(status) = lower(?)) AS transactions,
		       sumIf(amount, lower(status) = lower(?)) AS amt
		FROM transaction_daily FINAL
		WHERE toYear(event_time) = ? AND toMonth(event_time) = ? AND merchant_id = ?
		GROUP BY month, toMonth(event_time), payment_method
	`
	rows, err := r.conn.Query(ctx, query, status, status, year, month, merchantID)
	if err != nil {
		return nil, fmt.Errorf("query transaction method monthly by merchant: %w", err)
	}
	defer rows.Close()

	var results []TransactionMethodMonthly
	for rows.Next() {
		var t TransactionMethodMonthly
		if err := rows.Scan(&t.Month, &t.PaymentMethod, &t.Transactions, &t.Amount); err != nil {
			return nil, fmt.Errorf("scan transaction method monthly by merchant: %w", err)
		}
		results = append(results, t)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetTransactionMethodYearlyByMerchant(ctx context.Context, year, merchantID int, status string) ([]TransactionMethodYearly, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       payment_method,
		       countIf(lower(status) = lower(?)) AS transactions,
		       sumIf(amount, lower(status) = lower(?)) AS amt
		FROM transaction_daily FINAL
		WHERE toYear(event_time) = ? AND merchant_id = ?
		GROUP BY year, payment_method
	`
	rows, err := r.conn.Query(ctx, query, status, status, year, merchantID)
	if err != nil {
		return nil, fmt.Errorf("query transaction method yearly by merchant: %w", err)
	}
	defer rows.Close()

	var results []TransactionMethodYearly
	for rows.Next() {
		var t TransactionMethodYearly
		if err := rows.Scan(&t.Year, &t.PaymentMethod, &t.Transactions, &t.Amount); err != nil {
			return nil, fmt.Errorf("scan transaction method yearly by merchant: %w", err)
		}
		results = append(results, t)
	}
	return results, rows.Err()
}

// ─── Cashier stats ───────────────────────────────────────────────────────

func (r *ClickHouseReaderRepository) GetCashierMonthlyOrders(ctx context.Context, cashierID int) ([]CashierMonthlyOrders, error) {
	query := `
		SELECT formatDateTime(event_time, '%b') AS month,
		       cashier_id,
		       count() AS order_count,
		       sum(total_price) AS total_amount
		FROM order_daily FINAL
		WHERE cashier_id = ?
		GROUP BY month, toMonth(event_time), cashier_id
		ORDER BY toMonth(event_time)
	`
	rows, err := r.conn.Query(ctx, query, cashierID)
	if err != nil {
		return nil, fmt.Errorf("query cashier monthly orders: %w", err)
	}
	defer rows.Close()

	var results []CashierMonthlyOrders
	for rows.Next() {
		var c CashierMonthlyOrders
		if err := rows.Scan(&c.Month, &c.CashierID, &c.OrderCount, &c.TotalAmount); err != nil {
			return nil, fmt.Errorf("scan cashier monthly orders: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetCashierMonthlyTotalSales(ctx context.Context, year, month int) ([]CashierTotalSales, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       formatDateTime(event_time, '%b') AS month,
		       sum(total_price) AS total_sales
		FROM order_daily FINAL
		WHERE toYear(event_time) = ? AND toMonth(event_time) = ?
		GROUP BY year, month, toMonth(event_time)
	`
	rows, err := r.conn.Query(ctx, query, year, month)
	if err != nil {
		return nil, fmt.Errorf("query cashier monthly total sales: %w", err)
	}
	defer rows.Close()

	var results []CashierTotalSales
	for rows.Next() {
		var c CashierTotalSales
		if err := rows.Scan(&c.Year, &c.Month, &c.TotalSales); err != nil {
			return nil, fmt.Errorf("scan cashier monthly total sales: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetCashierYearlyTotalSales(ctx context.Context, year int) ([]CashierYearTotalSales, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       sum(total_price) AS total_sales
		FROM order_daily FINAL
		WHERE toYear(event_time) = ?
		GROUP BY year
	`
	rows, err := r.conn.Query(ctx, query, year)
	if err != nil {
		return nil, fmt.Errorf("query cashier yearly total sales: %w", err)
	}
	defer rows.Close()

	var results []CashierYearTotalSales
	for rows.Next() {
		var c CashierYearTotalSales
		if err := rows.Scan(&c.Year, &c.TotalSales); err != nil {
			return nil, fmt.Errorf("scan cashier yearly total sales: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetCashierMonthlyTotalSalesById(ctx context.Context, year, month, cashierID int) ([]CashierTotalSales, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       formatDateTime(event_time, '%b') AS month,
		       sum(total_price) AS total_sales
		FROM order_daily FINAL
		WHERE toYear(event_time) = ? AND toMonth(event_time) = ? AND cashier_id = ?
		GROUP BY year, month, toMonth(event_time)
	`
	rows, err := r.conn.Query(ctx, query, year, month, cashierID)
	if err != nil {
		return nil, fmt.Errorf("query cashier monthly total sales by id: %w", err)
	}
	defer rows.Close()

	var results []CashierTotalSales
	for rows.Next() {
		var c CashierTotalSales
		if err := rows.Scan(&c.Year, &c.Month, &c.TotalSales); err != nil {
			return nil, fmt.Errorf("scan cashier monthly total sales by id: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetCashierYearlyTotalSalesById(ctx context.Context, year, cashierID int) ([]CashierYearTotalSales, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       sum(total_price) AS total_sales
		FROM order_daily FINAL
		WHERE toYear(event_time) = ? AND cashier_id = ?
		GROUP BY year
	`
	rows, err := r.conn.Query(ctx, query, year, cashierID)
	if err != nil {
		return nil, fmt.Errorf("query cashier yearly total sales by id: %w", err)
	}
	defer rows.Close()

	var results []CashierYearTotalSales
	for rows.Next() {
		var c CashierYearTotalSales
		if err := rows.Scan(&c.Year, &c.TotalSales); err != nil {
			return nil, fmt.Errorf("scan cashier yearly total sales by id: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetCashierMonthlyTotalSalesByMerchant(ctx context.Context, year, month, merchantID int) ([]CashierTotalSales, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       formatDateTime(event_time, '%b') AS month,
		       sum(total_price) AS total_sales
		FROM order_daily FINAL
		WHERE toYear(event_time) = ? AND toMonth(event_time) = ? AND merchant_id = ?
		GROUP BY year, month, toMonth(event_time)
	`
	rows, err := r.conn.Query(ctx, query, year, month, merchantID)
	if err != nil {
		return nil, fmt.Errorf("query cashier monthly total sales by merchant: %w", err)
	}
	defer rows.Close()

	var results []CashierTotalSales
	for rows.Next() {
		var c CashierTotalSales
		if err := rows.Scan(&c.Year, &c.Month, &c.TotalSales); err != nil {
			return nil, fmt.Errorf("scan cashier monthly total sales by merchant: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetCashierYearlyTotalSalesByMerchant(ctx context.Context, year, merchantID int) ([]CashierYearTotalSales, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       sum(total_price) AS total_sales
		FROM order_daily FINAL
		WHERE toYear(event_time) = ? AND merchant_id = ?
		GROUP BY year
	`
	rows, err := r.conn.Query(ctx, query, year, merchantID)
	if err != nil {
		return nil, fmt.Errorf("query cashier yearly total sales by merchant: %w", err)
	}
	defer rows.Close()

	var results []CashierYearTotalSales
	for rows.Next() {
		var c CashierYearTotalSales
		if err := rows.Scan(&c.Year, &c.TotalSales); err != nil {
			return nil, fmt.Errorf("scan cashier yearly total sales by merchant: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetCashierSalesMonthly(ctx context.Context, year int) ([]CashierSalesMonthly, error) {
	query := `
		SELECT formatDateTime(event_time, '%b') AS month,
		       cashier_id,
		       count() AS order_count,
		       sum(total_price) AS total_sales
		FROM order_daily FINAL
		WHERE toYear(event_time) = ?
		GROUP BY month, toMonth(event_time), cashier_id
		ORDER BY toMonth(event_time)
	`
	rows, err := r.conn.Query(ctx, query, year)
	if err != nil {
		return nil, fmt.Errorf("query cashier sales monthly: %w", err)
	}
	defer rows.Close()

	var results []CashierSalesMonthly
	for rows.Next() {
		var c CashierSalesMonthly
		if err := rows.Scan(&c.Month, &c.CashierID, &c.OrderCount, &c.TotalSales); err != nil {
			return nil, fmt.Errorf("scan cashier sales monthly: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetCashierSalesYearly(ctx context.Context, year int) ([]CashierSalesYearly, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       cashier_id,
		       count() AS order_count,
		       sum(total_price) AS total_sales
		FROM order_daily FINAL
		WHERE toYear(event_time) = ?
		GROUP BY year, cashier_id
	`
	rows, err := r.conn.Query(ctx, query, year)
	if err != nil {
		return nil, fmt.Errorf("query cashier sales yearly: %w", err)
	}
	defer rows.Close()

	var results []CashierSalesYearly
	for rows.Next() {
		var c CashierSalesYearly
		if err := rows.Scan(&c.Year, &c.CashierID, &c.OrderCount, &c.TotalSales); err != nil {
			return nil, fmt.Errorf("scan cashier sales yearly: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetCashierSalesMonthlyByMerchant(ctx context.Context, year, merchantID int) ([]CashierSalesMonthly, error) {
	query := `
		SELECT formatDateTime(event_time, '%b') AS month,
		       cashier_id,
		       count() AS order_count,
		       sum(total_price) AS total_sales
		FROM order_daily FINAL
		WHERE toYear(event_time) = ? AND merchant_id = ?
		GROUP BY month, toMonth(event_time), cashier_id
		ORDER BY toMonth(event_time)
	`
	rows, err := r.conn.Query(ctx, query, year, merchantID)
	if err != nil {
		return nil, fmt.Errorf("query cashier sales monthly by merchant: %w", err)
	}
	defer rows.Close()

	var results []CashierSalesMonthly
	for rows.Next() {
		var c CashierSalesMonthly
		if err := rows.Scan(&c.Month, &c.CashierID, &c.OrderCount, &c.TotalSales); err != nil {
			return nil, fmt.Errorf("scan cashier sales monthly by merchant: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetCashierSalesYearlyByMerchant(ctx context.Context, year, merchantID int) ([]CashierSalesYearly, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       cashier_id,
		       count() AS order_count,
		       sum(total_price) AS total_sales
		FROM order_daily FINAL
		WHERE toYear(event_time) = ? AND merchant_id = ?
		GROUP BY year, cashier_id
	`
	rows, err := r.conn.Query(ctx, query, year, merchantID)
	if err != nil {
		return nil, fmt.Errorf("query cashier sales yearly by merchant: %w", err)
	}
	defer rows.Close()

	var results []CashierSalesYearly
	for rows.Next() {
		var c CashierSalesYearly
		if err := rows.Scan(&c.Year, &c.CashierID, &c.OrderCount, &c.TotalSales); err != nil {
			return nil, fmt.Errorf("scan cashier sales yearly by merchant: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetCashierSalesMonthlyById(ctx context.Context, year, cashierID int) ([]CashierSalesMonthly, error) {
	query := `
		SELECT formatDateTime(event_time, '%b') AS month,
		       cashier_id,
		       count() AS order_count,
		       sum(total_price) AS total_sales
		FROM order_daily FINAL
		WHERE toYear(event_time) = ? AND cashier_id = ?
		GROUP BY month, toMonth(event_time), cashier_id
		ORDER BY toMonth(event_time)
	`
	rows, err := r.conn.Query(ctx, query, year, cashierID)
	if err != nil {
		return nil, fmt.Errorf("query cashier sales monthly by id: %w", err)
	}
	defer rows.Close()

	var results []CashierSalesMonthly
	for rows.Next() {
		var c CashierSalesMonthly
		if err := rows.Scan(&c.Month, &c.CashierID, &c.OrderCount, &c.TotalSales); err != nil {
			return nil, fmt.Errorf("scan cashier sales monthly by id: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

func (r *ClickHouseReaderRepository) GetCashierSalesYearlyById(ctx context.Context, year, cashierID int) ([]CashierSalesYearly, error) {
	query := `
		SELECT toString(toYear(event_time)) AS year,
		       cashier_id,
		       count() AS order_count,
		       sum(total_price) AS total_sales
		FROM order_daily FINAL
		WHERE toYear(event_time) = ? AND cashier_id = ?
		GROUP BY year, cashier_id
	`
	rows, err := r.conn.Query(ctx, query, year, cashierID)
	if err != nil {
		return nil, fmt.Errorf("query cashier sales yearly by id: %w", err)
	}
	defer rows.Close()

	var results []CashierSalesYearly
	for rows.Next() {
		var c CashierSalesYearly
		if err := rows.Scan(&c.Year, &c.CashierID, &c.OrderCount, &c.TotalSales); err != nil {
			return nil, fmt.Errorf("scan cashier sales yearly by id: %w", err)
		}
		results = append(results, c)
	}
	return results, rows.Err()
}
