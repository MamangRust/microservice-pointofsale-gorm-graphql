// Package cashier reads cashier sales statistics from ClickHouse's
// order_daily table. It replaces the cashier slice of the old single
// stats-reader Repository god interface.
package cashier

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

// NewRepository returns the ClickHouse-backed cashier stats repository.
func NewRepository(conn clickhouse.Conn) Repository {
	return &clickhouseRepository{conn: conn}
}

type orderFilter struct {
	year, month, merchantID, cashierID, orderID int
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

// ── Cashier stats ───────────────────────────────────────────────────────────

type cashierAggRow struct {
	Year       string
	Month      string
	CashierID  uint64
	OrderCount uint64
	TotalSales int64
}

// cashierAggregate groups order_daily by month and/or cashier. When byCashier
// is false the cashier_id slot is a constant so the queries share one shape.
func (r *clickhouseRepository) cashierAggregate(ctx context.Context, f orderFilter, byMonth, byCashier bool) ([]cashierAggRow, error) {
	where, args := orderWhere(f)
	keyExpr, orderExpr, groupExtra := "''", "year", ""
	if byMonth {
		keyExpr, orderExpr, groupExtra = "formatDateTime(event_time, '%b')", "toMonth(event_time)", ", toMonth(event_time)"
	}
	cashierExpr := "toUInt64(0)"
	if byCashier {
		cashierExpr = "cashier_id"
	}
	query := fmt.Sprintf(`
		SELECT toString(toYear(event_time)) AS year,
		       %s AS month_key,
		       %s AS cashier_id,
		       count() AS order_count,
		       sum(total_price) AS total_sales
		FROM order_daily FINAL
		WHERE %s
		GROUP BY year, month_key, cashier_id%s
		ORDER BY %s, total_sales DESC
	`, keyExpr, cashierExpr, where, groupExtra, orderExpr)

	var results []cashierAggRow
	err := r.query(ctx, "cashier aggregate", query, args, func(rows driver.Rows) error {
		var row cashierAggRow
		if err := rows.Scan(&row.Year, &row.Month, &row.CashierID, &row.OrderCount, &row.TotalSales); err != nil {
			return err
		}
		results = append(results, row)
		return nil
	})
	return results, err
}

func (r *clickhouseRepository) GetMonthlyTotalSales(ctx context.Context, year, month int) ([]CashierMonthTotalSales, error) {
	rows, err := r.cashierAggregate(ctx, orderFilter{year: year, month: month}, true, false)
	if err != nil {
		return nil, err
	}
	out := make([]CashierMonthTotalSales, 0, len(rows))
	for _, row := range rows {
		out = append(out, CashierMonthTotalSales{Year: row.Year, Month: row.Month, TotalSales: row.TotalSales})
	}
	return out, nil
}

func (r *clickhouseRepository) GetYearlyTotalSales(ctx context.Context, year int) ([]CashierYearTotalSales, error) {
	rows, err := r.cashierAggregate(ctx, orderFilter{year: year}, false, false)
	if err != nil {
		return nil, err
	}
	out := make([]CashierYearTotalSales, 0, len(rows))
	for _, row := range rows {
		out = append(out, CashierYearTotalSales{Year: row.Year, TotalSales: row.TotalSales})
	}
	return out, nil
}

// GetMonthSales returns per-cashier monthly sales. cashier_name is not stored
// in ClickHouse; the front end resolves it via the cashier domain service.
func (r *clickhouseRepository) GetMonthSales(ctx context.Context, year int) ([]CashierMonthSales, error) {
	return cashierMonthSales(r.cashierAggregate(ctx, orderFilter{year: year}, true, true))
}

func (r *clickhouseRepository) GetYearSales(ctx context.Context, year int) ([]CashierYearSales, error) {
	return cashierYearSales(r.cashierAggregate(ctx, orderFilter{year: year}, false, true))
}

func (r *clickhouseRepository) GetMonthlyTotalSalesByMerchant(ctx context.Context, year, month, merchantID int) ([]CashierMonthTotalSales, error) {
	rows, err := r.cashierAggregate(ctx, orderFilter{year: year, month: month, merchantID: merchantID}, true, false)
	if err != nil {
		return nil, err
	}
	out := make([]CashierMonthTotalSales, 0, len(rows))
	for _, row := range rows {
		out = append(out, CashierMonthTotalSales{Year: row.Year, Month: row.Month, TotalSales: row.TotalSales})
	}
	return out, nil
}

func (r *clickhouseRepository) GetYearlyTotalSalesByMerchant(ctx context.Context, year, merchantID int) ([]CashierYearTotalSales, error) {
	rows, err := r.cashierAggregate(ctx, orderFilter{year: year, merchantID: merchantID}, false, false)
	if err != nil {
		return nil, err
	}
	out := make([]CashierYearTotalSales, 0, len(rows))
	for _, row := range rows {
		out = append(out, CashierYearTotalSales{Year: row.Year, TotalSales: row.TotalSales})
	}
	return out, nil
}

func (r *clickhouseRepository) GetMonthSalesByMerchant(ctx context.Context, year, merchantID int) ([]CashierMonthSales, error) {
	return cashierMonthSales(r.cashierAggregate(ctx, orderFilter{year: year, merchantID: merchantID}, true, true))
}

func (r *clickhouseRepository) GetYearSalesByMerchant(ctx context.Context, year, merchantID int) ([]CashierYearSales, error) {
	return cashierYearSales(r.cashierAggregate(ctx, orderFilter{year: year, merchantID: merchantID}, false, true))
}

func (r *clickhouseRepository) GetMonthlyTotalSalesById(ctx context.Context, year, month, cashierID int) ([]CashierMonthTotalSales, error) {
	rows, err := r.cashierAggregate(ctx, orderFilter{year: year, month: month, cashierID: cashierID}, true, false)
	if err != nil {
		return nil, err
	}
	out := make([]CashierMonthTotalSales, 0, len(rows))
	for _, row := range rows {
		out = append(out, CashierMonthTotalSales{Year: row.Year, Month: row.Month, TotalSales: row.TotalSales})
	}
	return out, nil
}

func (r *clickhouseRepository) GetYearlyTotalSalesById(ctx context.Context, year, cashierID int) ([]CashierYearTotalSales, error) {
	rows, err := r.cashierAggregate(ctx, orderFilter{year: year, cashierID: cashierID}, false, false)
	if err != nil {
		return nil, err
	}
	out := make([]CashierYearTotalSales, 0, len(rows))
	for _, row := range rows {
		out = append(out, CashierYearTotalSales{Year: row.Year, TotalSales: row.TotalSales})
	}
	return out, nil
}

func (r *clickhouseRepository) GetMonthSalesById(ctx context.Context, year, cashierID int) ([]CashierMonthSales, error) {
	return cashierMonthSales(r.cashierAggregate(ctx, orderFilter{year: year, cashierID: cashierID}, true, true))
}

func (r *clickhouseRepository) GetYearSalesById(ctx context.Context, year, cashierID int) ([]CashierYearSales, error) {
	return cashierYearSales(r.cashierAggregate(ctx, orderFilter{year: year, cashierID: cashierID}, false, true))
}

func cashierMonthSales(rows []cashierAggRow, err error) ([]CashierMonthSales, error) {
	if err != nil {
		return nil, err
	}
	out := make([]CashierMonthSales, 0, len(rows))
	for _, row := range rows {
		out = append(out, CashierMonthSales{Month: row.Month, CashierID: row.CashierID, OrderCount: row.OrderCount, TotalSales: row.TotalSales})
	}
	return out, nil
}

func cashierYearSales(rows []cashierAggRow, err error) ([]CashierYearSales, error) {
	if err != nil {
		return nil, err
	}
	out := make([]CashierYearSales, 0, len(rows))
	for _, row := range rows {
		out = append(out, CashierYearSales{Year: row.Year, CashierID: row.CashierID, OrderCount: row.OrderCount, TotalSales: row.TotalSales})
	}
	return out, nil
}
