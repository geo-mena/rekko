package repository

import (
	"context"

	"github.com/geomena/stock-recommendation-system/backend/internal/domain"
	"github.com/google/uuid"
)

type MockStockRepository struct {
	CreateFn                  func(ctx context.Context, stock *domain.Stock) error
	FindByIDFn                func(ctx context.Context, id uuid.UUID) (*domain.Stock, error)
	FindByTickerFn            func(ctx context.Context, ticker string) ([]domain.Stock, error)
	FindAllFn                 func(ctx context.Context, filter domain.StockFilter) ([]domain.Stock, int64, error)
	BulkUpsertFn              func(ctx context.Context, stocks []domain.Stock) (int, error)
	GetDistinctActionsFn      func(ctx context.Context) ([]string, error)
	CountAllFn                func(ctx context.Context) (int64, error)
	GetActionDistributionFn   func(ctx context.Context) ([]domain.ActionDistribution, error)
	GetBrokerageDistributionFn func(ctx context.Context, limit int) ([]domain.BrokerageDistribution, error)
	GetRecentActivityFn       func(ctx context.Context, days int) ([]domain.DailyActivity, error)
}

func (m *MockStockRepository) Create(ctx context.Context, stock *domain.Stock) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, stock)
	}
	return nil
}

func (m *MockStockRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Stock, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockStockRepository) FindByTicker(ctx context.Context, ticker string) ([]domain.Stock, error) {
	if m.FindByTickerFn != nil {
		return m.FindByTickerFn(ctx, ticker)
	}
	return nil, nil
}

func (m *MockStockRepository) FindAll(ctx context.Context, filter domain.StockFilter) ([]domain.Stock, int64, error) {
	if m.FindAllFn != nil {
		return m.FindAllFn(ctx, filter)
	}
	return nil, 0, nil
}

func (m *MockStockRepository) BulkUpsert(ctx context.Context, stocks []domain.Stock) (int, error) {
	if m.BulkUpsertFn != nil {
		return m.BulkUpsertFn(ctx, stocks)
	}
	return 0, nil
}

func (m *MockStockRepository) GetDistinctActions(ctx context.Context) ([]string, error) {
	if m.GetDistinctActionsFn != nil {
		return m.GetDistinctActionsFn(ctx)
	}
	return nil, nil
}

func (m *MockStockRepository) CountAll(ctx context.Context) (int64, error) {
	if m.CountAllFn != nil {
		return m.CountAllFn(ctx)
	}
	return 0, nil
}

func (m *MockStockRepository) GetActionDistribution(ctx context.Context) ([]domain.ActionDistribution, error) {
	if m.GetActionDistributionFn != nil {
		return m.GetActionDistributionFn(ctx)
	}
	return nil, nil
}

func (m *MockStockRepository) GetBrokerageDistribution(ctx context.Context, limit int) ([]domain.BrokerageDistribution, error) {
	if m.GetBrokerageDistributionFn != nil {
		return m.GetBrokerageDistributionFn(ctx, limit)
	}
	return nil, nil
}

func (m *MockStockRepository) GetRecentActivity(ctx context.Context, days int) ([]domain.DailyActivity, error) {
	if m.GetRecentActivityFn != nil {
		return m.GetRecentActivityFn(ctx, days)
	}
	return nil, nil
}
