package usecase

import (
	"context"

	"github.com/geomena/stock-recommendation-system/backend/internal/domain"
	"github.com/geomena/stock-recommendation-system/backend/internal/repository"
)

type DashboardUsecase struct {
	stockRepo repository.StockRepository
}

func NewDashboardUsecase(stockRepo repository.StockRepository) *DashboardUsecase {
	return &DashboardUsecase{stockRepo: stockRepo}
}

func (u *DashboardUsecase) GetDashboardStats(ctx context.Context) (*domain.DashboardStats, error) {
	totalStocks, err := u.stockRepo.CountAll(ctx)
	if err != nil {
		return nil, err
	}

	actionDist, err := u.stockRepo.GetActionDistribution(ctx)
	if err != nil {
		return nil, err
	}
	if actionDist == nil {
		actionDist = []domain.ActionDistribution{}
	}

	brokerageDist, err := u.stockRepo.GetBrokerageDistribution(ctx, 10)
	if err != nil {
		return nil, err
	}
	if brokerageDist == nil {
		brokerageDist = []domain.BrokerageDistribution{}
	}

	recentActivity, err := u.stockRepo.GetRecentActivity(ctx, 30)
	if err != nil {
		return nil, err
	}
	if recentActivity == nil {
		recentActivity = []domain.DailyActivity{}
	}

	return &domain.DashboardStats{
		TotalStocks:           totalStocks,
		ActionDistribution:    actionDist,
		BrokerageDistribution: brokerageDist,
		RecentActivity:        recentActivity,
	}, nil
}
