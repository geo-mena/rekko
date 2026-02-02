package repository

import (
	"context"

	"github.com/geomena/stock-recommendation-system/backend/internal/domain"
	"github.com/google/uuid"
)

type StockRepository interface {
	Create(ctx context.Context, stock *domain.Stock) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Stock, error)
	FindByTicker(ctx context.Context, ticker string) ([]domain.Stock, error)
	FindAll(ctx context.Context, filter domain.StockFilter) ([]domain.Stock, int64, error)
	BulkUpsert(ctx context.Context, stocks []domain.Stock) (int, error)
	GetDistinctActions(ctx context.Context) ([]string, error)
	CountAll(ctx context.Context) (int64, error)
	GetActionDistribution(ctx context.Context) ([]domain.ActionDistribution, error)
	GetBrokerageDistribution(ctx context.Context, limit int) ([]domain.BrokerageDistribution, error)
	GetRecentActivity(ctx context.Context, days int) ([]domain.DailyActivity, error)
}
