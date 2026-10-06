// Package transaction reads transaction status/payment-method
// statistics from ClickHouse's transaction_daily table. It replaces the
// transaction slice of the old single stats-reader Repository god
// interface.
package transaction

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

// NewRepository returns the ClickHouse-backed transaction stats repository.
func NewRepository(conn clickhouse.Conn) Repository {
	return &clickhouseRepository{conn: conn}
}

type transactionFilter struct {
	year, month, merchantID int
}

func transactionWhere(f transactionFilter) (string, []any) {
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
	return strings.Join(conds, " AND "), args
}

// POS stores both "success" (API creates) and "Completed"/"completed" (legacy
// seeder data) for a successful payment; match case-insensitively so the reader
// cross-checks with the OLTP payment_status values. Anything else counts as
// failed.
const successStatus = "lower(status) IN ('success', 'completed')"

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

// ── Transaction stats ───────────────────────────────────────────────────────

type txAmountRow struct {
	Year        string
	Month       string
	TotalCount  uint64
	TotalAmount int64
}

func (r *clickhouseRepository) transactionStatusAggregate(ctx context.Context, f transactionFilter, byMonth, success bool) ([]txAmountRow, error) {
	where, args := transactionWhere(f)
	statusCond := successStatus
	if !success {
		statusCond = "NOT " + successStatus
	}
	keyExpr, orderExpr, groupExtra := "''", "year", ""
	if byMonth {
		keyExpr, orderExpr, groupExtra = "formatDateTime(event_time, '%b')", "toMonth(event_time)", ", toMonth(event_time)"
	}
	query := fmt.Sprintf(`
		SELECT toString(toYear(event_time)) AS year,
		       %s AS month_key,
		       countIf(%s) AS total_count,
		       sumIf(amount, %s) AS total_amount
		FROM transaction_daily FINAL
		WHERE %s
		GROUP BY year, month_key%s
		ORDER BY %s
	`, keyExpr, statusCond, statusCond, where, groupExtra, orderExpr)

	var results []txAmountRow
	err := r.query(ctx, "transaction status aggregate", query, args, func(rows driver.Rows) error {
		var row txAmountRow
		if err := rows.Scan(&row.Year, &row.Month, &row.TotalCount, &row.TotalAmount); err != nil {
			return err
		}
		results = append(results, row)
		return nil
	})
	return results, err
}

func (r *clickhouseRepository) GetMonthStatusSuccess(ctx context.Context, year, month int) ([]TransactionMonthAmount, error) {
	return transactionMonthAmount(r.transactionStatusAggregate(ctx, transactionFilter{year: year, month: month}, true, true))
}

func (r *clickhouseRepository) GetYearStatusSuccess(ctx context.Context, year int) ([]TransactionYearAmount, error) {
	return transactionYearAmount(r.transactionStatusAggregate(ctx, transactionFilter{year: year}, false, true))
}

func (r *clickhouseRepository) GetMonthStatusFailed(ctx context.Context, year, month int) ([]TransactionMonthAmount, error) {
	return transactionMonthAmount(r.transactionStatusAggregate(ctx, transactionFilter{year: year, month: month}, true, false))
}

func (r *clickhouseRepository) GetYearStatusFailed(ctx context.Context, year int) ([]TransactionYearAmount, error) {
	return transactionYearAmount(r.transactionStatusAggregate(ctx, transactionFilter{year: year}, false, false))
}

func (r *clickhouseRepository) GetMonthStatusSuccessByMerchant(ctx context.Context, year, month, merchantID int) ([]TransactionMonthAmount, error) {
	return transactionMonthAmount(r.transactionStatusAggregate(ctx, transactionFilter{year: year, month: month, merchantID: merchantID}, true, true))
}

func (r *clickhouseRepository) GetYearStatusSuccessByMerchant(ctx context.Context, year, merchantID int) ([]TransactionYearAmount, error) {
	return transactionYearAmount(r.transactionStatusAggregate(ctx, transactionFilter{year: year, merchantID: merchantID}, false, true))
}

func (r *clickhouseRepository) GetMonthStatusFailedByMerchant(ctx context.Context, year, month, merchantID int) ([]TransactionMonthAmount, error) {
	return transactionMonthAmount(r.transactionStatusAggregate(ctx, transactionFilter{year: year, month: month, merchantID: merchantID}, true, false))
}

func (r *clickhouseRepository) GetYearStatusFailedByMerchant(ctx context.Context, year, merchantID int) ([]TransactionYearAmount, error) {
	return transactionYearAmount(r.transactionStatusAggregate(ctx, transactionFilter{year: year, merchantID: merchantID}, false, false))
}

type txMethodRow struct {
	Year              string
	Month             string
	PaymentMethod     string
	TotalTransactions uint64
	TotalAmount       int64
}

func (r *clickhouseRepository) transactionMethodAggregate(ctx context.Context, f transactionFilter, byMonth, success bool) ([]txMethodRow, error) {
	where, args := transactionWhere(f)
	statusCond := successStatus
	if !success {
		statusCond = "NOT " + successStatus
	}
	keyExpr, orderExpr, groupExtra := "''", "year", ""
	if byMonth {
		keyExpr, orderExpr, groupExtra = "formatDateTime(event_time, '%b')", "toMonth(event_time)", ", toMonth(event_time)"
	}
	query := fmt.Sprintf(`
		SELECT toString(toYear(event_time)) AS year,
		       %s AS month_key,
		       payment_method,
		       count() AS total_transactions,
		       sum(amount) AS total_amount
		FROM transaction_daily FINAL
		WHERE %s AND %s
		GROUP BY year, month_key, payment_method%s
		ORDER BY %s, total_amount DESC
	`, keyExpr, where, statusCond, groupExtra, orderExpr)

	var results []txMethodRow
	err := r.query(ctx, "transaction method aggregate", query, args, func(rows driver.Rows) error {
		var row txMethodRow
		if err := rows.Scan(&row.Year, &row.Month, &row.PaymentMethod, &row.TotalTransactions, &row.TotalAmount); err != nil {
			return err
		}
		results = append(results, row)
		return nil
	})
	return results, err
}

func (r *clickhouseRepository) GetMonthMethodSuccess(ctx context.Context, year, month int) ([]TransactionMonthMethod, error) {
	return transactionMonthMethod(r.transactionMethodAggregate(ctx, transactionFilter{year: year, month: month}, true, true))
}

func (r *clickhouseRepository) GetYearMethodSuccess(ctx context.Context, year int) ([]TransactionYearMethod, error) {
	return transactionYearMethod(r.transactionMethodAggregate(ctx, transactionFilter{year: year}, false, true))
}

func (r *clickhouseRepository) GetMonthMethodFailed(ctx context.Context, year, month int) ([]TransactionMonthMethod, error) {
	return transactionMonthMethod(r.transactionMethodAggregate(ctx, transactionFilter{year: year, month: month}, true, false))
}

func (r *clickhouseRepository) GetYearMethodFailed(ctx context.Context, year int) ([]TransactionYearMethod, error) {
	return transactionYearMethod(r.transactionMethodAggregate(ctx, transactionFilter{year: year}, false, false))
}

func (r *clickhouseRepository) GetMonthMethodByMerchantSuccess(ctx context.Context, year, month, merchantID int) ([]TransactionMonthMethod, error) {
	return transactionMonthMethod(r.transactionMethodAggregate(ctx, transactionFilter{year: year, month: month, merchantID: merchantID}, true, true))
}

func (r *clickhouseRepository) GetYearMethodByMerchantSuccess(ctx context.Context, year, merchantID int) ([]TransactionYearMethod, error) {
	return transactionYearMethod(r.transactionMethodAggregate(ctx, transactionFilter{year: year, merchantID: merchantID}, false, true))
}

func (r *clickhouseRepository) GetMonthMethodByMerchantFailed(ctx context.Context, year, month, merchantID int) ([]TransactionMonthMethod, error) {
	return transactionMonthMethod(r.transactionMethodAggregate(ctx, transactionFilter{year: year, month: month, merchantID: merchantID}, true, false))
}

func (r *clickhouseRepository) GetYearMethodByMerchantFailed(ctx context.Context, year, merchantID int) ([]TransactionYearMethod, error) {
	return transactionYearMethod(r.transactionMethodAggregate(ctx, transactionFilter{year: year, merchantID: merchantID}, false, false))
}

func transactionMonthAmount(rows []txAmountRow, err error) ([]TransactionMonthAmount, error) {
	if err != nil {
		return nil, err
	}
	out := make([]TransactionMonthAmount, 0, len(rows))
	for _, row := range rows {
		out = append(out, TransactionMonthAmount{Year: row.Year, Month: row.Month, TotalCount: row.TotalCount, TotalAmount: row.TotalAmount})
	}
	return out, nil
}

func transactionYearAmount(rows []txAmountRow, err error) ([]TransactionYearAmount, error) {
	if err != nil {
		return nil, err
	}
	out := make([]TransactionYearAmount, 0, len(rows))
	for _, row := range rows {
		out = append(out, TransactionYearAmount{Year: row.Year, TotalCount: row.TotalCount, TotalAmount: row.TotalAmount})
	}
	return out, nil
}

func transactionMonthMethod(rows []txMethodRow, err error) ([]TransactionMonthMethod, error) {
	if err != nil {
		return nil, err
	}
	out := make([]TransactionMonthMethod, 0, len(rows))
	for _, row := range rows {
		out = append(out, TransactionMonthMethod{Month: row.Month, PaymentMethod: row.PaymentMethod, TotalTransactions: row.TotalTransactions, TotalAmount: row.TotalAmount})
	}
	return out, nil
}

func transactionYearMethod(rows []txMethodRow, err error) ([]TransactionYearMethod, error) {
	if err != nil {
		return nil, err
	}
	out := make([]TransactionYearMethod, 0, len(rows))
	for _, row := range rows {
		out = append(out, TransactionYearMethod{Year: row.Year, PaymentMethod: row.PaymentMethod, TotalTransactions: row.TotalTransactions, TotalAmount: row.TotalAmount})
	}
	return out, nil
}
