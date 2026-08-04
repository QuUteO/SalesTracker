package handler

import (
	"SalesTracker/internal/model"
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/wb-go/wbf/ginext"
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

type SalesHandler struct {
	srv SalesService
}

func New(srv SalesService) *SalesHandler {
	return &SalesHandler{srv: srv}
}

// InitRoutes настраивает роутинг с использованием ginext.Engine
func (h *SalesHandler) InitRoutes(ginMode string) *ginext.Engine {
	engine := ginext.New(ginMode)

	// Стандартные мидлвары ginext
	engine.Use(ginext.Logger(), ginext.Recovery())

	// Статика для веб-интерфейса
	engine.StaticFile("/", "./web/index.html")

	v1 := engine.Group("/api/v1")
	{
		// 1. CRUD транзакций/продаж (/items)
		items := v1.Group("/items")
		{
			items.POST("", h.CreateItem)       // POST /api/v1/items
			items.GET("", h.ListItems)         // GET /api/v1/items?from=...&to=...&category=...&sort_by=amount&order=desc
			items.GET("/:id", h.GetItemByID)   // GET /api/v1/items/{id}
			items.PUT("/:id", h.UpdateItem)    // PUT /api/v1/items/{id}
			items.DELETE("/:id", h.DeleteItem) // DELETE /api/v1/items/{id}
		}

		// 2. Сводная аналитика (sum, avg, count, median, p90)
		analytics := v1.Group("/analytics")
		{
			analytics.GET("", h.GetAnalytics) // GET /api/v1/analytics?from=...&to=...
		}

		// 3. Экспорт отчетов
		export := v1.Group("/export")
		{
			export.GET("/csv", h.ExportCSV) // GET /api/v1/export/csv?from=...&to=...
		}
	}

	return engine
}

// renderError централизованно обрабатывает ошибки доменного слоя
func renderError(c *ginext.Context, err error) {
	status := http.StatusInternalServerError

	switch {
	case errors.Is(err, model.ErrItemNotFound):
		status = http.StatusNotFound

	case errors.Is(err, model.ErrNegativeAmount),
		errors.Is(err, model.ErrInvalidDateRange),
		errors.Is(err, model.ErrEmptyTitle),
		errors.Is(err, model.ErrInvalidType):
		status = http.StatusBadRequest
	}

	c.JSON(status, ginext.H{"error": err.Error()})
}
