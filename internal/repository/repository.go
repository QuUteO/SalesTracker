package repository

import (
	"SalesTracker/internal/model"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	pgxdriver "github.com/wb-go/wbf/dbpg/pgx-driver"
)

type SalesRepository struct {
	conn *pgxdriver.Postgres
}

func New(conn *pgxdriver.Postgres) *SalesRepository {
	return &SalesRepository{conn: conn}
}

func (r *SalesRepository) Create(ctx context.Context, item *model.Item) error {
	query := `INSERT INTO items(id, title, amount, type, category, date, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.conn.Exec(ctx, query,
		item.ID,
		item.Title,
		item.Amount,
		item.Type,
		item.Category,
		item.Date,
		item.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("could not create item: %w", err)
	}

	return nil
}

func (r *SalesRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Item, error) {
	query := `SELECT id, title, amount, type, category, date, created_at FROM items WHERE id = $1`

	var item model.Item
	err := r.conn.QueryRow(ctx, query, id).Scan(
		&item.ID,
		&item.Title,
		&item.Amount,
		&item.Type,
		&item.Category,
		&item.Date,
		&item.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrItemNotFound
		}
		return nil, fmt.Errorf("could not get item: %w", err)
	}

	return &item, nil
}

func (r *SalesRepository) List(ctx context.Context, filter model.AnalyticsFilter) ([]*model.Item, error) {
	var builder strings.Builder
	builder.WriteString(`
		SELECT id, title, amount, type, category, date, created_at 
		FROM items 
		WHERE 1=1
	`)

	args := make([]any, 0, 4)
	argID := 1

	if filter.Type != "" {
		builder.WriteString(fmt.Sprintf(" AND type = $%d", argID))
		args = append(args, filter.Type)
		argID++
	}

	if filter.Category != "" {
		builder.WriteString(fmt.Sprintf(" AND category = $%d", argID))
		args = append(args, filter.Category)
		argID++
	}

	if filter.From != nil {
		builder.WriteString(fmt.Sprintf(" AND date >= $%d", argID))
		args = append(args, *filter.From)
		argID++
	}

	if filter.To != nil {
		builder.WriteString(fmt.Sprintf(" AND date <= $%d", argID))
		args = append(args, *filter.To)
		argID++
	}

	sortBy := "created_at"
	switch strings.ToLower(filter.SortBy) {
	case "amount":
		sortBy = "amount"
	case "date":
		sortBy = "date"
	case "title":
		sortBy = "title"
	}

	order := "DESC"
	if strings.ToLower(filter.Order) == "asc" {
		order = "ASC"
	}

	builder.WriteString(fmt.Sprintf(" ORDER BY %s %s", sortBy, order))

	rows, err := r.conn.Query(ctx, builder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query items: %w", err)
	}
	defer rows.Close()

	items := make([]*model.Item, 0)
	for rows.Next() {
		item := new(model.Item)
		err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Amount,
			&item.Type,
			&item.Category,
			&item.Date,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return items, nil
}

func (r *SalesRepository) Update(ctx context.Context, id uuid.UUID, item *model.UpdateItemRequest) error {
	query := `UPDATE items SET title=$2, amount=$3, type=$4, category=$5, date=$6 WHERE id = $1`

	_, err := r.conn.Exec(ctx, query,
		id,
		item.Title,
		item.Amount,
		item.Type,
		item.Category,
		item.Date,
	)
	if err != nil {
		return fmt.Errorf("could not update item: %w", err)
	}

	return nil
}

func (r *SalesRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM items WHERE id = $1`

	_, err := r.conn.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("could not delete item: %w", err)
	}

	return nil
}

func (r *SalesRepository) GetAnalytics(ctx context.Context, filter model.AnalyticsFilter) (*model.AnalyticsResponse, error) {
	var builder strings.Builder
	builder.WriteString(`
		SELECT 
			COUNT(*) AS total_count,
			COALESCE(SUM(amount), 0) AS total_sum,
			COALESCE(AVG(amount), 0) AS avg_amount,
			COALESCE(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY amount), 0) AS median,
			COALESCE(PERCENTILE_CONT(0.9) WITHIN GROUP (ORDER BY amount), 0) AS percentile_90
		FROM items 
		WHERE 1=1
	`)

	args := make([]any, 0, 4)
	argID := 1

	if filter.Type != "" {
		builder.WriteString(fmt.Sprintf(" AND type = $%d", argID))
		args = append(args, filter.Type)
		argID++
	}

	if filter.Category != "" {
		builder.WriteString(fmt.Sprintf(" AND category = $%d", argID))
		args = append(args, filter.Category)
		argID++
	}

	if filter.From != nil {
		builder.WriteString(fmt.Sprintf(" AND date >= $%d", argID))
		args = append(args, *filter.From)
		argID++
	}

	if filter.To != nil {
		builder.WriteString(fmt.Sprintf(" AND date <= $%d", argID))
		args = append(args, *filter.To)
		argID++
	}

	var analytics model.AnalyticsResponse
	err := r.conn.QueryRow(ctx, builder.String(), args...).Scan(
		&analytics.Count,
		&analytics.Sum,
		&analytics.Avg,
		&analytics.Median,
		&analytics.Percentile90,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate analytics: %w", err)
	}

	return &analytics, nil
}
