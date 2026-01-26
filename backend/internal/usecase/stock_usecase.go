package usecase

import (
	"context"
	"math"

	"github.com/geomena/stock-recommendation-system/backend/internal/domain"
	"github.com/geomena/stock-recommendation-system/backend/internal/external/karenai"
	"github.com/geomena/stock-recommendation-system/backend/internal/repository"
	"github.com/google/uuid"
)

type StockUsecase struct {
	stockRepo     repository.StockRepository
	karenaiClient *karenai.Client
}

func NewStockUsecase(stockRepo repository.StockRepository, karenaiClient *karenai.Client) *StockUsecase {
	return &StockUsecase{
		stockRepo:     stockRepo,
		karenaiClient: karenaiClient,
	}
}

func (u *StockUsecase) ListStocks(ctx context.Context, filter domain.StockFilter) (*domain.PaginatedStocks, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 20
	}

	stocks, totalCount, err := u.stockRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(filter.Limit)))

	return &domain.PaginatedStocks{
		Data:       stocks,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
		HasNext:    filter.Page < totalPages,
		HasPrev:    filter.Page > 1,
	}, nil
}

func (u *StockUsecase) GetStockByID(ctx context.Context, id uuid.UUID) (*domain.Stock, error) {
	return u.stockRepo.FindByID(ctx, id)
}

func (u *StockUsecase) GetStocksByTicker(ctx context.Context, ticker string) ([]domain.Stock, error) {
	return u.stockRepo.FindByTicker(ctx, ticker)
}

func (u *StockUsecase) GetDistinctActions(ctx context.Context) ([]string, error) {
	return u.stockRepo.GetDistinctActions(ctx)
}

func (u *StockUsecase) SyncFromExternalAPI(ctx context.Context) (int, error) {
	stocks, err := u.karenaiClient.FetchAllStocks(ctx)
	if err != nil {
		return 0, err
	}

	inserted, err := u.stockRepo.BulkUpsert(ctx, stocks)
	if err != nil {
		return 0, err
	}

	return inserted, nil
}
