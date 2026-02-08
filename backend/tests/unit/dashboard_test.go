package unit_test

import (
	"context"
	"errors"
	"testing"

	"github.com/geomena/stock-recommendation-system/backend/internal/domain"
	"github.com/geomena/stock-recommendation-system/backend/internal/repository"
)

func TestGetDashboardStats_Success(t *testing.T) {
	mock := newMockRepo()
	mock.CountAllFn = func(ctx context.Context) (int64, error) {
		return 42, nil
	}
	mock.GetActionDistributionFn = func(ctx context.Context) ([]domain.ActionDistribution, error) {
		return []domain.ActionDistribution{
			{Action: "upgraded", Count: 15},
			{Action: "downgraded", Count: 5},
		}, nil
	}
	mock.GetBrokerageDistributionFn = func(ctx context.Context, limit int) ([]domain.BrokerageDistribution, error) {
		if limit != 10 {
			t.Errorf("expected limit 10, got %d", limit)
		}
		return []domain.BrokerageDistribution{
			{Brokerage: "Morgan Stanley", Count: 18},
		}, nil
	}
	mock.GetRecentActivityFn = func(ctx context.Context, days int) ([]domain.DailyActivity, error) {
		if days != 30 {
			t.Errorf("expected days 30, got %d", days)
		}
		return []domain.DailyActivity{
			{Date: "2025-01-15", Count: 5},
		}, nil
	}

	uc := newDashboardUsecase(mock)
	stats, err := uc.GetDashboardStats(context.Background())
	assertNoError(t, err)

	if stats.TotalStocks != 42 {
		t.Errorf("expected TotalStocks 42, got %d", stats.TotalStocks)
	}
	if len(stats.ActionDistribution) != 2 {
		t.Errorf("expected 2 action distributions, got %d", len(stats.ActionDistribution))
	}
	if len(stats.BrokerageDistribution) != 1 {
		t.Errorf("expected 1 brokerage distribution, got %d", len(stats.BrokerageDistribution))
	}
	if len(stats.RecentActivity) != 1 {
		t.Errorf("expected 1 recent activity, got %d", len(stats.RecentActivity))
	}
}

func TestGetDashboardStats_NilNormalized(t *testing.T) {
	mock := newMockRepo()
	mock.CountAllFn = func(ctx context.Context) (int64, error) {
		return 0, nil
	}
	mock.GetActionDistributionFn = func(ctx context.Context) ([]domain.ActionDistribution, error) {
		return nil, nil
	}
	mock.GetBrokerageDistributionFn = func(ctx context.Context, limit int) ([]domain.BrokerageDistribution, error) {
		return nil, nil
	}
	mock.GetRecentActivityFn = func(ctx context.Context, days int) ([]domain.DailyActivity, error) {
		return nil, nil
	}

	uc := newDashboardUsecase(mock)
	stats, err := uc.GetDashboardStats(context.Background())
	assertNoError(t, err)

	if stats.ActionDistribution == nil {
		t.Error("expected ActionDistribution to be non-nil empty slice")
	}
	if stats.BrokerageDistribution == nil {
		t.Error("expected BrokerageDistribution to be non-nil empty slice")
	}
	if stats.RecentActivity == nil {
		t.Error("expected RecentActivity to be non-nil empty slice")
	}
}

func TestGetDashboardStats_ErrorPropagation(t *testing.T) {
	repoErr := errors.New("db error")

	tests := []struct {
		name    string
		setupFn func(mock *repository.MockStockRepository)
	}{
		{
			name: "CountAll fails",
			setupFn: func(mock *repository.MockStockRepository) {
				mock.CountAllFn = func(ctx context.Context) (int64, error) {
					return 0, repoErr
				}
			},
		},
		{
			name: "GetActionDistribution fails",
			setupFn: func(mock *repository.MockStockRepository) {
				mock.CountAllFn = func(ctx context.Context) (int64, error) {
					return 10, nil
				}
				mock.GetActionDistributionFn = func(ctx context.Context) ([]domain.ActionDistribution, error) {
					return nil, repoErr
				}
			},
		},
		{
			name: "GetBrokerageDistribution fails",
			setupFn: func(mock *repository.MockStockRepository) {
				mock.CountAllFn = func(ctx context.Context) (int64, error) {
					return 10, nil
				}
				mock.GetActionDistributionFn = func(ctx context.Context) ([]domain.ActionDistribution, error) {
					return []domain.ActionDistribution{}, nil
				}
				mock.GetBrokerageDistributionFn = func(ctx context.Context, limit int) ([]domain.BrokerageDistribution, error) {
					return nil, repoErr
				}
			},
		},
		{
			name: "GetRecentActivity fails",
			setupFn: func(mock *repository.MockStockRepository) {
				mock.CountAllFn = func(ctx context.Context) (int64, error) {
					return 10, nil
				}
				mock.GetActionDistributionFn = func(ctx context.Context) ([]domain.ActionDistribution, error) {
					return []domain.ActionDistribution{}, nil
				}
				mock.GetBrokerageDistributionFn = func(ctx context.Context, limit int) ([]domain.BrokerageDistribution, error) {
					return []domain.BrokerageDistribution{}, nil
				}
				mock.GetRecentActivityFn = func(ctx context.Context, days int) ([]domain.DailyActivity, error) {
					return nil, repoErr
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock := newMockRepo()
			tc.setupFn(mock)

			uc := newDashboardUsecase(mock)
			stats, err := uc.GetDashboardStats(context.Background())
			assertError(t, err)

			if stats != nil {
				t.Errorf("expected nil stats on error, got %+v", stats)
			}
		})
	}
}
