package cockroachdb

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/geomena/stock-recommendation-system/backend/internal/domain"
	"github.com/google/uuid"
)

type StockRepository struct {
	db *DB
}

func NewStockRepository(db *DB) *StockRepository {
	return &StockRepository{db: db}
}

func (r *StockRepository) Create(ctx context.Context, stock *domain.Stock) error {
	query := `
		INSERT INTO stocks (ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	return r.db.Conn().QueryRowContext(ctx, query,
		stock.Ticker,
		stock.Company,
		stock.Brokerage,
		stock.Action,
		stock.RatingFrom,
		stock.RatingTo,
		stock.TargetFrom,
		stock.TargetTo,
	).Scan(&stock.ID, &stock.CreatedAt, &stock.UpdatedAt)
}

func (r *StockRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Stock, error) {
	query := `
		SELECT id, ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, created_at, updated_at
		FROM stocks
		WHERE id = $1`

	stock := &domain.Stock{}
	err := r.db.Conn().QueryRowContext(ctx, query, id).Scan(
		&stock.ID,
		&stock.Ticker,
		&stock.Company,
		&stock.Brokerage,
		&stock.Action,
		&stock.RatingFrom,
		&stock.RatingTo,
		&stock.TargetFrom,
		&stock.TargetTo,
		&stock.CreatedAt,
		&stock.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrStockNotFound
	}
	if err != nil {
		return nil, err
	}
	return stock, nil
}

func (r *StockRepository) FindByTicker(ctx context.Context, ticker string) ([]domain.Stock, error) {
	query := `
		SELECT id, ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, created_at, updated_at
		FROM stocks
		WHERE ticker = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Conn().QueryContext(ctx, query, ticker)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanStocks(rows)
}

func (r *StockRepository) FindAll(ctx context.Context, filter domain.StockFilter) ([]domain.Stock, int64, error) {
	baseQuery := "FROM stocks WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if filter.Search != "" {
		baseQuery += fmt.Sprintf(" AND (ticker ILIKE $%d OR company ILIKE $%d)", argIndex, argIndex+1)
		searchPattern := "%" + filter.Search + "%"
		args = append(args, searchPattern, searchPattern)
		argIndex += 2
	}

	if filter.Ticker != "" {
		baseQuery += fmt.Sprintf(" AND ticker = $%d", argIndex)
		args = append(args, filter.Ticker)
		argIndex++
	}

	if filter.Action != "" {
		baseQuery += fmt.Sprintf(" AND action ILIKE $%d", argIndex)
		args = append(args, "%"+filter.Action+"%")
		argIndex++
	}

	countQuery := "SELECT COUNT(*) " + baseQuery
	var totalCount int64
	if err := r.db.Conn().QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		return nil, 0, err
	}

	sortColumn := sanitizeSortColumn(filter.SortBy)
	sortOrder := sanitizeSortOrder(filter.SortOrder)

	selectQuery := fmt.Sprintf(`
		SELECT id, ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, created_at, updated_at
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d`,
		baseQuery, sortColumn, sortOrder, argIndex, argIndex+1)

	offset := (filter.Page - 1) * filter.Limit
	args = append(args, filter.Limit, offset)

	rows, err := r.db.Conn().QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	stocks, err := scanStocks(rows)
	if err != nil {
		return nil, 0, err
	}

	return stocks, totalCount, nil
}

func (r *StockRepository) BulkUpsert(ctx context.Context, stocks []domain.Stock) (int, error) {
	if len(stocks) == 0 {
		return 0, nil
	}

	query := `
		INSERT INTO stocks (ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (ticker, brokerage, action, rating_from, rating_to, target_from, target_to)
		DO UPDATE SET updated_at = NOW()`

	inserted := 0
	for _, stock := range stocks {
		result, err := r.db.Conn().ExecContext(ctx, query,
			stock.Ticker,
			stock.Company,
			stock.Brokerage,
			stock.Action,
			stock.RatingFrom,
			stock.RatingTo,
			stock.TargetFrom,
			stock.TargetTo,
		)
		if err != nil {
			continue
		}
		if affected, _ := result.RowsAffected(); affected > 0 {
			inserted++
		}
	}

	return inserted, nil
}

func (r *StockRepository) GetDistinctActions(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT action FROM stocks ORDER BY action`

	rows, err := r.db.Conn().QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []string
	for rows.Next() {
		var action string
		if err := rows.Scan(&action); err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}

	return actions, rows.Err()
}

func scanStocks(rows *sql.Rows) ([]domain.Stock, error) {
	var stocks []domain.Stock
	for rows.Next() {
		var stock domain.Stock
		if err := rows.Scan(
			&stock.ID,
			&stock.Ticker,
			&stock.Company,
			&stock.Brokerage,
			&stock.Action,
			&stock.RatingFrom,
			&stock.RatingTo,
			&stock.TargetFrom,
			&stock.TargetTo,
			&stock.CreatedAt,
			&stock.UpdatedAt,
		); err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}
	return stocks, rows.Err()
}

func sanitizeSortColumn(column string) string {
	allowed := map[string]string{
		"ticker":     "ticker",
		"company":    "company",
		"action":     "action",
		"targetTo":   "target_to",
		"target_to":  "target_to",
		"createdAt":  "created_at",
		"created_at": "created_at",
	}
	if col, ok := allowed[column]; ok {
		return col
	}
	return "created_at"
}

func sanitizeSortOrder(order string) string {
	order = strings.ToUpper(order)
	if order == "ASC" {
		return "ASC"
	}
	return "DESC"
}

func (r *StockRepository) CountAll(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.Conn().QueryRowContext(ctx, "SELECT COUNT(*) FROM stocks").Scan(&count)
	return count, err
}

func (r *StockRepository) GetActionDistribution(ctx context.Context) ([]domain.ActionDistribution, error) {
	rows, err := r.db.Conn().QueryContext(ctx,
		"SELECT action, COUNT(*) as count FROM stocks GROUP BY action ORDER BY count DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.ActionDistribution
	for rows.Next() {
		var item domain.ActionDistribution
		if err := rows.Scan(&item.Action, &item.Count); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (r *StockRepository) GetBrokerageDistribution(ctx context.Context, limit int) ([]domain.BrokerageDistribution, error) {
	rows, err := r.db.Conn().QueryContext(ctx,
		"SELECT brokerage, COUNT(*) as count FROM stocks GROUP BY brokerage ORDER BY count DESC LIMIT $1", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.BrokerageDistribution
	for rows.Next() {
		var item domain.BrokerageDistribution
		if err := rows.Scan(&item.Brokerage, &item.Count); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (r *StockRepository) GetRecentActivity(ctx context.Context, days int) ([]domain.DailyActivity, error) {
	rows, err := r.db.Conn().QueryContext(ctx, `
		SELECT TO_CHAR(created_at::DATE, 'YYYY-MM-DD') AS date, COUNT(*) AS count
		FROM stocks
		WHERE created_at >= NOW() - CAST($1 || ' days' AS INTERVAL)
		GROUP BY created_at::DATE
		ORDER BY date ASC`, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.DailyActivity
	for rows.Next() {
		var item domain.DailyActivity
		if err := rows.Scan(&item.Date, &item.Count); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}
