package service

import (
	"SalesTracker/internal/model"
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, item *model.Item) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Item, error)
	List(ctx context.Context, filter model.AnalyticsFilter) ([]*model.Item, error)
	Update(ctx context.Context, id uuid.UUID, item *model.UpdateItemRequest) error
	Delete(ctx context.Context, id uuid.UUID) error

	GetAnalytics(ctx context.Context, filter model.AnalyticsFilter) (*model.AnalyticsResponse, error)
}

type SalesService struct {
	repo Repository
}

func New(repo Repository) *SalesService {
	return &SalesService{repo: repo}
}

func (s *SalesService) Create(ctx context.Context, req *model.CreateItemRequest) (*model.Item, error) {
	if req.Date.IsZero() {
		req.Date = time.Now()
	}

	if req.Amount < 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}

	item := &model.Item{
		ID:        uuid.New(),
		Title:     req.Title,
		Amount:    req.Amount,
		Type:      req.Type,
		Category:  req.Category,
		Date:      req.Date,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return nil, fmt.Errorf("create item: %w", err)
	}

	return item, nil
}

func (s *SalesService) GetByID(ctx context.Context, id uuid.UUID) (*model.Item, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get item: %w", err)
	}

	return item, nil
}

func (s *SalesService) List(ctx context.Context, filter model.AnalyticsFilter) ([]*model.Item, error) {
	// 1. Валидация диапазона дат: 'From' не должна быть позже 'To'
	if filter.From != nil && filter.To != nil {
		if filter.From.After(*filter.To) {
			return nil, model.ErrInvalidDateRange
		}
	}

	// 2. Установка значений по умолчанию для сортировки
	if strings.TrimSpace(filter.SortBy) == "" {
		filter.SortBy = "date"
	}

	if strings.TrimSpace(filter.Order) == "" {
		filter.Order = "desc"
	}

	items, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("service.List: %w", err)
	}

	return items, nil
}

func (s *SalesService) Update(ctx context.Context, id uuid.UUID, req *model.UpdateItemRequest) (*model.Item, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil || item == nil {
		return nil, fmt.Errorf("get item: %w", err)
	}

	if req.Date.IsZero() {
		req.Date = time.Now()
	}

	if err := s.repo.Update(ctx, id, req); err != nil {
		return nil, fmt.Errorf("update item: %w", err)
	}

	return item, nil
}

func (s *SalesService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete item: %w", err)
	}
	return nil
}

func (s *SalesService) GetAnalytics(ctx context.Context, filter model.AnalyticsFilter) (*model.AnalyticsResponse, error) {
	if filter.From != nil && filter.To != nil {
		if filter.From.After(*filter.To) {
			return nil, model.ErrInvalidDateRange
		}
	}

	analytics, err := s.repo.GetAnalytics(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("service.GetAnalytics: %w", err)
	}

	return analytics, nil
}

func (s *SalesService) ExportCSV(ctx context.Context, filter model.AnalyticsFilter) ([]byte, error) {
	items, err := s.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("service.ExportCSV: %w", err)
	}

	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	headers := []string{"ID", "Название", "Сумма", "Тип", "Категория", "Дата", "Дата создания"}
	if err := writer.Write(headers); err != nil {
		return nil, fmt.Errorf("failed to write csv headers: %w", err)
	}

	for _, item := range items {
		record := []string{
			item.ID.String(),
			item.Title,
			fmt.Sprintf("%.2f", item.Amount),
			string(item.Type),
			item.Category,
			item.Date.Format("2006-01-02 15:04:05"),
			item.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("failed to write csv record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("csv writer flush error: %w", err)
	}

	return buf.Bytes(), nil
}
