package handler

import (
	"SalesTracker/internal/model"
	"context"

	"github.com/google/uuid"
)

type SalesService interface {
	Create(ctx context.Context, req *model.CreateItemRequest) (*model.Item, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Item, error)
	List(ctx context.Context, filter model.AnalyticsFilter) ([]*model.Item, error)
	Update(ctx context.Context, id uuid.UUID, req *model.UpdateItemRequest) (*model.Item, error)
	Delete(ctx context.Context, id uuid.UUID) error

	GetAnalytics(ctx context.Context, filter model.AnalyticsFilter) (*model.AnalyticsResponse, error)
	ExportCSV(ctx context.Context, filter model.AnalyticsFilter) ([]byte, error)
}
