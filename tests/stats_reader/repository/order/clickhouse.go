// Package order reads order statistics (revenue, order counts and
// items sold) from ClickHouse. It replaces the order slice of the old
// single stats-reader Repository god interface.
package order

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

// NewRepository returns the ClickHouse-backed order stats repository.
func NewRepository(conn clickhouse.Conn) Repository {
	return &clickhouseRepository{conn: conn}
}

type orderFilter struct {
	year, month, merchantID, cashierID, orderID int
}
type itemFilter struct {
	year, month, merchantID, categoryID, orderID int
}

func orderWhere(f orderFilter) (string, []any) {
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
		conds = append(conds, "merchant_id = ?")
		args = append(args, f.merchantID)
	}
	if f.cashierID > 0 {
		conds = append(conds, "cashier_id = ?")
		args = append(args, f.cashierID)
	}
	if f.orderID > 0 {
		conds = append(conds, "order_id = ?")
		args = append(args, f.orderID)
	}
	return strings.Join(conds, " AND "), args
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

// ── Order stats ─────────────────────────────────────────────────────────────
//
// order_daily holds revenue/order/cashier aggregates; total_items_sold and
// unique_products_sold only exist in order_item_daily, so they are aggregated
// separately and merged in Go on the month/year key.

type orderAggRow struct {
	Year           string
	Month          string
	OrderCount     uint64
	TotalRevenue   int64
	ActiveCashiers uint64
}

type itemAggRow struct {
	ItemsSold      uint64
	UniqueProducts uint64
}

func (r *clickhouseRepository) orderAggregate(ctx context.Context, f orderFilter, byMonth bool) ([]orderAggRow, error) {
	where, args := orderWhere(f)
	keyExpr, orderExpr, groupExtra := "''", "year", ""
	if byMonth {
		keyExpr, orderExpr, groupExtra = "formatDateTime(event_time, '%b')", "toMonth(event_time)", ", toMonth(event_time)"
	}
	query := fmt.Sprintf(`
		SELECT toString(toYear(event_time)) AS year,
		       %s AS month_key,
		       count() AS order_count,
		       sum(total_price) AS total_revenue,
		       uniqExact(cashier_id) AS active_cashiers
		FROM order_daily FINAL
		WHERE %s
		GROUP BY year, month_key%s
		ORDER BY %s
	`, keyExpr, where, groupExtra, orderExpr)

	var results []orderAggRow
	err := r.query(ctx, "order aggregate", query, args, func(rows driver.Rows) error {
		var row orderAggRow
		if err := rows.Scan(&row.Year, &row.Month, &row.OrderCount, &row.TotalRevenue, &row.ActiveCashiers); err != nil {
			return err
		}
		results = append(results, row)
		return nil
	})
	return results, err
}

func (r *clickhouseRepository) itemAggregate(ctx context.Context, f itemFilter, byMonth bool) (map[string]itemAggRow, error) {
	where, args := itemWhere(f)
	keyExpr, orderExpr, groupExtra := "toString(toYear(event_time))", "year", ""
	if byMonth {
		keyExpr, orderExpr, groupExtra = "formatDateTime(event_time, '%b')", "toMonth(event_time)", ", toMonth(event_time)"
	}
	query := fmt.Sprintf(`
		SELECT %s AS item_key,
		       sum(quantity) AS items_sold,
		       uniqExact(product_id) AS unique_products
		FROM order_item_daily FINAL
		WHERE %s
		GROUP BY item_key%s
		ORDER BY %s
	`, keyExpr, where, groupExtra, orderExpr)

	items := map[string]itemAggRow{}
	err := r.query(ctx, "item aggregate", query, args, func(rows driver.Rows) error {
		var key string
		var agg itemAggRow
		if err := rows.Scan(&key, &agg.ItemsSold, &agg.UniqueProducts); err != nil {
			return err
		}
		items[key] = agg
		return nil
	})
	return items, err
}

func mergeOrderMonthlyTotalRevenue(rows []orderAggRow, items map[string]itemAggRow) []OrderMonthlyTotalRevenue {
	out := make([]OrderMonthlyTotalRevenue, 0, len(rows))
	for _, row := range rows {
		out = append(out, OrderMonthlyTotalRevenue{
			Year:           row.Year,
			Month:          row.Month,
			OrderCount:     row.OrderCount,
			TotalRevenue:   row.TotalRevenue,
			TotalItemsSold: items[row.Month].ItemsSold,
		})
	}
	return out
}

func mergeOrderYearly(rows []orderAggRow, items map[string]itemAggRow) []OrderYearlyTotalRevenue {
	out := make([]OrderYearlyTotalRevenue, 0, len(rows))
	for _, row := range rows {
		out = append(out, OrderYearlyTotalRevenue{
			Year:               row.Year,
			OrderCount:         row.OrderCount,
			TotalRevenue:       row.TotalRevenue,
			TotalItemsSold:     items[row.Year].ItemsSold,
			ActiveCashiers:     row.ActiveCashiers,
			UniqueProductsSold: items[row.Year].UniqueProducts,
		})
	}
	return out
}

func (r *clickhouseRepository) GetMonthlyTotalRevenue(ctx context.Context, year, month int) ([]OrderMonthlyTotalRevenue, error) {
	f := orderFilter{year: year, month: month}
	rows, err := r.orderAggregate(ctx, f, true)
	if err != nil {
		return nil, err
	}
	items, err := r.itemAggregate(ctx, itemFilter{year: year, month: month}, true)
	if err != nil {
		return nil, err
	}
	return mergeOrderMonthlyTotalRevenue(rows, items), nil
}

func (r *clickhouseRepository) GetYearlyTotalRevenue(ctx context.Context, year int) ([]OrderYearlyTotalRevenue, error) {
	return r.orderYearlyTotalRevenue(ctx, orderFilter{year: year})
}

// orderYearlyTotalRevenue backs both FindYearlyTotalRevenue and
// FindYearlyRevenue — their pb responses carry the same aggregates.
func (r *clickhouseRepository) orderYearlyTotalRevenue(ctx context.Context, f orderFilter) ([]OrderYearlyTotalRevenue, error) {
	rows, err := r.orderAggregate(ctx, f, false)
	if err != nil {
		return nil, err
	}
	items, err := r.itemAggregate(ctx, itemFilter{year: f.year, merchantID: f.merchantID, orderID: f.orderID}, false)
	if err != nil {
		return nil, err
	}
	return mergeOrderYearly(rows, items), nil
}

func (r *clickhouseRepository) GetMonthlyRevenue(ctx context.Context, year int) ([]OrderMonthly, error) {
	f := orderFilter{year: year}
	rows, err := r.orderAggregate(ctx, f, true)
	if err != nil {
		return nil, err
	}
	items, err := r.itemAggregate(ctx, itemFilter{year: year}, true)
	if err != nil {
		return nil, err
	}
	out := make([]OrderMonthly, 0, len(rows))
	for _, row := range rows {
		out = append(out, OrderMonthly{
			Month:          row.Month,
			OrderCount:     row.OrderCount,
			TotalRevenue:   row.TotalRevenue,
			TotalItemsSold: items[row.Month].ItemsSold,
		})
	}
	return out, nil
}

func (r *clickhouseRepository) GetYearlyRevenue(ctx context.Context, year int) ([]OrderYearly, error) {
	rows, err := r.orderYearlyTotalRevenue(ctx, orderFilter{year: year})
	if err != nil {
		return nil, err
	}
	out := make([]OrderYearly, 0, len(rows))
	for _, row := range rows {
		out = append(out, OrderYearly{
			Year:               row.Year,
			OrderCount:         row.OrderCount,
			TotalRevenue:       row.TotalRevenue,
			TotalItemsSold:     row.TotalItemsSold,
			ActiveCashiers:     row.ActiveCashiers,
			UniqueProductsSold: row.UniqueProductsSold,
		})
	}
	return out, nil
}

func (r *clickhouseRepository) GetMonthlyTotalRevenueByMerchant(ctx context.Context, year, month, merchantID int) ([]OrderMonthlyTotalRevenue, error) {
	f := orderFilter{year: year, month: month, merchantID: merchantID}
	rows, err := r.orderAggregate(ctx, f, true)
	if err != nil {
		return nil, err
	}
	items, err := r.itemAggregate(ctx, itemFilter{year: year, month: month, merchantID: merchantID}, true)
	if err != nil {
		return nil, err
	}
	return mergeOrderMonthlyTotalRevenue(rows, items), nil
}

func (r *clickhouseRepository) GetYearlyTotalRevenueByMerchant(ctx context.Context, year, merchantID int) ([]OrderYearlyTotalRevenue, error) {
	return r.orderYearlyTotalRevenue(ctx, orderFilter{year: year, merchantID: merchantID})
}

func (r *clickhouseRepository) GetMonthlyRevenueByMerchant(ctx context.Context, year, merchantID int) ([]OrderMonthly, error) {
	f := orderFilter{year: year, merchantID: merchantID}
	rows, err := r.orderAggregate(ctx, f, true)
	if err != nil {
		return nil, err
	}
	items, err := r.itemAggregate(ctx, itemFilter{year: year, merchantID: merchantID}, true)
	if err != nil {
		return nil, err
	}
	out := make([]OrderMonthly, 0, len(rows))
	for _, row := range rows {
		out = append(out, OrderMonthly{
			Month:          row.Month,
			OrderCount:     row.OrderCount,
			TotalRevenue:   row.TotalRevenue,
			TotalItemsSold: items[row.Month].ItemsSold,
		})
	}
	return out, nil
}

func (r *clickhouseRepository) GetYearlyRevenueByMerchant(ctx context.Context, year, merchantID int) ([]OrderYearly, error) {
	rows, err := r.orderYearlyTotalRevenue(ctx, orderFilter{year: year, merchantID: merchantID})
	if err != nil {
		return nil, err
	}
	out := make([]OrderYearly, 0, len(rows))
	for _, row := range rows {
		out = append(out, OrderYearly{
			Year:               row.Year,
			OrderCount:         row.OrderCount,
			TotalRevenue:       row.TotalRevenue,
			TotalItemsSold:     row.TotalItemsSold,
			ActiveCashiers:     row.ActiveCashiers,
			UniqueProductsSold: row.UniqueProductsSold,
		})
	}
	return out, nil
}

func (r *clickhouseRepository) GetMonthlyTotalRevenueById(ctx context.Context, year, month, orderID int) ([]OrderMonthlyTotalRevenue, error) {
	f := orderFilter{year: year, month: month, orderID: orderID}
	rows, err := r.orderAggregate(ctx, f, true)
	if err != nil {
		return nil, err
	}
	items, err := r.itemAggregate(ctx, itemFilter{year: year, month: month, orderID: orderID}, true)
	if err != nil {
		return nil, err
	}
	return mergeOrderMonthlyTotalRevenue(rows, items), nil
}

func (r *clickhouseRepository) GetYearlyTotalRevenueById(ctx context.Context, year, orderID int) ([]OrderYearlyTotalRevenue, error) {
	return r.orderYearlyTotalRevenue(ctx, orderFilter{year: year, orderID: orderID})
}
