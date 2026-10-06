// Package category reads category price statistics from ClickHouse's
// order_item_daily table. It replaces the category slice of the old
// single stats-reader Repository god interface.
package category

import (
	"context"
	"fmt"
	"strings"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type clickhouseRepository struct {
	conn clickhouse.Conn
}

// NewRepository returns the ClickHouse-backed category stats repository.
func NewRepository(conn clickhouse.Conn) Repository {
	return &clickhouseRepository{conn: conn}
}

type itemFilter struct {
	year, month, merchantID, categoryID, orderID int
}

func itemWhere(f itemFilter) (string, []any) {
	var conds []string
	var args []any
	if f.year > 0 {
		conds = append(conds, "toYear(event_time) = ?")
		args = append(args, f.year)
	}
	if f.month > 0 {
		conds = append(conds, "toMonth(event_time) = ?")
		args = append(args, f.month)
	}
	if f.merchantID > 0 {
		conds = append(conds, "order_id IN (SELECT order_id FROM order_daily FINAL WHERE merchant_id = ?)")
		args = append(args, f.merchantID)
	}
	if f.categoryID > 0 {
		conds = append(conds, "category_id = ?")
		args = append(args, f.categoryID)
	}
	if f.orderID > 0 {
		conds = append(conds, "order_id = ?")
		args = append(args, f.orderID)
	}
	return strings.Join(conds, " AND "), args
}
func (r *clickhouseRepository) query(ctx context.Context, name, query string, args []any, scan func(rows driver.Rows) error) error {
	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("query %s: %w", name, err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := scan(rows); err != nil {
			return fmt.Errorf("scan %s: %w", name, err)
		}
	}
	return rows.Err()
}

// ── Category stats ──────────────────────────────────────────────────────────

type categoryAggRow struct {
	Year           string
	Month          string
	CategoryID     uint64
	OrderCount     uint64
	ItemsSold      uint64
	TotalRevenue   int64
	UniqueProducts uint64
}

// categoryAggregate groups order_item_daily by month and/or category. When
// byCategory is false the category_id/unique_products slots are constants.
func (r *clickhouseRepository) categoryAggregate(ctx context.Context, f itemFilter, byMonth, byCategory bool) ([]categoryAggRow, error) {
	where, args := itemWhere(f)
	keyExpr, orderExpr, groupExtra := "''", "year", ""
	if byMonth {
		keyExpr, orderExpr, groupExtra = "formatDateTime(event_time, '%b')", "toMonth(event_time)", ", toMonth(event_time)"
	}
	categoryExpr, productExpr := "toUInt64(0)", "toUInt64(0)"
	if byCategory {
		categoryExpr, productExpr = "category_id", "uniqExact(product_id)"
	}
	query := fmt.Sprintf(`
		SELECT toString(toYear(event_time)) AS year,
		       %s AS month_key,
		       %s AS category_id,
		       uniqExact(order_id) AS order_count,
		       sum(quantity) AS items_sold,
		       sum(subtotal) AS total_revenue,
		       %s AS unique_products
		FROM order_item_daily FINAL
		WHERE %s
		GROUP BY year, month_key, category_id%s
		ORDER BY %s, total_revenue DESC
	`, keyExpr, categoryExpr, productExpr, where, groupExtra, orderExpr)

	var results []categoryAggRow
	err := r.query(ctx, "category aggregate", query, args, func(rows driver.Rows) error {
		var row categoryAggRow
		if err := rows.Scan(&row.Year, &row.Month, &row.CategoryID, &row.OrderCount, &row.ItemsSold, &row.TotalRevenue, &row.UniqueProducts); err != nil {
			return err
		}
		results = append(results, row)
		return nil
	})
	return results, err
}

func (r *clickhouseRepository) GetMonthlyTotalPrices(ctx context.Context, year, month int) ([]CategoryMonthTotalPrice, error) {
	rows, err := r.categoryAggregate(ctx, itemFilter{year: year, month: month}, true, false)
	if err != nil {
		return nil, err
	}
	out := make([]CategoryMonthTotalPrice, 0, len(rows))
	for _, row := range rows {
		out = append(out, CategoryMonthTotalPrice{Year: row.Year, Month: row.Month, TotalRevenue: row.TotalRevenue})
	}
	return out, nil
}

func (r *clickhouseRepository) GetYearlyTotalPrices(ctx context.Context, year int) ([]CategoryYearTotalPrice, error) {
	rows, err := r.categoryAggregate(ctx, itemFilter{year: year}, false, false)
	if err != nil {
		return nil, err
	}
	out := make([]CategoryYearTotalPrice, 0, len(rows))
	for _, row := range rows {
		out = append(out, CategoryYearTotalPrice{Year: row.Year, TotalRevenue: row.TotalRevenue})
	}
	return out, nil
}

// GetMonthPrice returns per-category monthly price aggregates. category_name is
// not stored in ClickHouse; the front end resolves it via the category domain
// service.
func (r *clickhouseRepository) GetMonthPrice(ctx context.Context, year int) ([]CategoryMonthPrice, error) {
	return categoryMonthPrice(r.categoryAggregate(ctx, itemFilter{year: year}, true, true))
}

func (r *clickhouseRepository) GetYearPrice(ctx context.Context, year int) ([]CategoryYearPrice, error) {
	return categoryYearPrice(r.categoryAggregate(ctx, itemFilter{year: year}, false, true))
}

func (r *clickhouseRepository) GetMonthlyTotalPricesByMerchant(ctx context.Context, year, month, merchantID int) ([]CategoryMonthTotalPrice, error) {
	rows, err := r.categoryAggregate(ctx, itemFilter{year: year, month: month, merchantID: merchantID}, true, false)
	if err != nil {
		return nil, err
	}
	out := make([]CategoryMonthTotalPrice, 0, len(rows))
	for _, row := range rows {
		out = append(out, CategoryMonthTotalPrice{Year: row.Year, Month: row.Month, TotalRevenue: row.TotalRevenue})
	}
	return out, nil
}

func (r *clickhouseRepository) GetYearlyTotalPricesByMerchant(ctx context.Context, year, merchantID int) ([]CategoryYearTotalPrice, error) {
	rows, err := r.categoryAggregate(ctx, itemFilter{year: year, merchantID: merchantID}, false, false)
	if err != nil {
		return nil, err
	}
	out := make([]CategoryYearTotalPrice, 0, len(rows))
	for _, row := range rows {
		out = append(out, CategoryYearTotalPrice{Year: row.Year, TotalRevenue: row.TotalRevenue})
	}
	return out, nil
}

func (r *clickhouseRepository) GetMonthPriceByMerchant(ctx context.Context, year, merchantID int) ([]CategoryMonthPrice, error) {
	return categoryMonthPrice(r.categoryAggregate(ctx, itemFilter{year: year, merchantID: merchantID}, true, true))
}

func (r *clickhouseRepository) GetYearPriceByMerchant(ctx context.Context, year, merchantID int) ([]CategoryYearPrice, error) {
	return categoryYearPrice(r.categoryAggregate(ctx, itemFilter{year: year, merchantID: merchantID}, false, true))
}

func (r *clickhouseRepository) GetMonthlyTotalPricesById(ctx context.Context, year, month, categoryID int) ([]CategoryMonthTotalPrice, error) {
	rows, err := r.categoryAggregate(ctx, itemFilter{year: year, month: month, categoryID: categoryID}, true, false)
	if err != nil {
		return nil, err
	}
	out := make([]CategoryMonthTotalPrice, 0, len(rows))
	for _, row := range rows {
		out = append(out, CategoryMonthTotalPrice{Year: row.Year, Month: row.Month, TotalRevenue: row.TotalRevenue})
	}
	return out, nil
}

func (r *clickhouseRepository) GetYearlyTotalPricesById(ctx context.Context, year, categoryID int) ([]CategoryYearTotalPrice, error) {
	rows, err := r.categoryAggregate(ctx, itemFilter{year: year, categoryID: categoryID}, false, false)
	if err != nil {
		return nil, err
	}
	out := make([]CategoryYearTotalPrice, 0, len(rows))
	for _, row := range rows {
		out = append(out, CategoryYearTotalPrice{Year: row.Year, TotalRevenue: row.TotalRevenue})
	}
	return out, nil
}

func (r *clickhouseRepository) GetMonthPriceById(ctx context.Context, year, categoryID int) ([]CategoryMonthPrice, error) {
	return categoryMonthPrice(r.categoryAggregate(ctx, itemFilter{year: year, categoryID: categoryID}, true, true))
}

func (r *clickhouseRepository) GetYearPriceById(ctx context.Context, year, categoryID int) ([]CategoryYearPrice, error) {
	return categoryYearPrice(r.categoryAggregate(ctx, itemFilter{year: year, categoryID: categoryID}, false, true))
}

func categoryMonthPrice(rows []categoryAggRow, err error) ([]CategoryMonthPrice, error) {
	if err != nil {
		return nil, err
	}
	out := make([]CategoryMonthPrice, 0, len(rows))
	for _, row := range rows {
		out = append(out, CategoryMonthPrice{Month: row.Month, CategoryID: row.CategoryID, OrderCount: row.OrderCount, ItemsSold: row.ItemsSold, TotalRevenue: row.TotalRevenue})
	}
	return out, nil
}

func categoryYearPrice(rows []categoryAggRow, err error) ([]CategoryYearPrice, error) {
	if err != nil {
		return nil, err
	}
	out := make([]CategoryYearPrice, 0, len(rows))
	for _, row := range rows {
		out = append(out, CategoryYearPrice{Year: row.Year, CategoryID: row.CategoryID, OrderCount: row.OrderCount, ItemsSold: row.ItemsSold, TotalRevenue: row.TotalRevenue, UniqueProductsSold: row.UniqueProducts})
	}
	return out, nil
}
