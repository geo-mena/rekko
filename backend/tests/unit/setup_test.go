package unit_test

import (
	"testing"
	"time"

	"github.com/geomena/stock-recommendation-system/backend/internal/domain"
	"github.com/geomena/stock-recommendation-system/backend/internal/repository"
	"github.com/geomena/stock-recommendation-system/backend/internal/usecase"
	"github.com/google/uuid"
)

var (
	stockID1 = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	stockID2 = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	stockID3 = uuid.MustParse("33333333-3333-3333-3333-333333333333")
	stockID4 = uuid.MustParse("44444444-4444-4444-4444-444444444444")
	stockID5 = uuid.MustParse("55555555-5555-5555-5555-555555555555")
	stockID6 = uuid.MustParse("66666666-6666-6666-6666-666666666666")

	fixedNow = time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
)

func newMockRepo() *repository.MockStockRepository {
	return &repository.MockStockRepository{}
}

func newStockUsecase(mock *repository.MockStockRepository) *usecase.StockUsecase {
	return usecase.NewStockUsecase(mock, nil)
}

func newRecommendationUsecase(mock *repository.MockStockRepository) *usecase.RecommendationUsecase {
	return usecase.NewRecommendationUsecase(mock, nil)
}

func newDashboardUsecase(mock *repository.MockStockRepository) *usecase.DashboardUsecase {
	return usecase.NewDashboardUsecase(mock)
}

func makeStock(id uuid.UUID, ticker, company, brokerage, action, ratingFrom, ratingTo string, targetFrom, targetTo float64) domain.Stock {
	return domain.Stock{
		ID:         id,
		Ticker:     ticker,
		Company:    company,
		Brokerage:  brokerage,
		Action:     action,
		RatingFrom: ratingFrom,
		RatingTo:   ratingTo,
		TargetFrom: targetFrom,
		TargetTo:   targetTo,
		CreatedAt:  fixedNow,
		UpdatedAt:  fixedNow,
	}
}

func makeStockAt(id uuid.UUID, ticker, company, brokerage, action, ratingFrom, ratingTo string, targetFrom, targetTo float64, createdAt time.Time) domain.Stock {
	s := makeStock(id, ticker, company, brokerage, action, ratingFrom, ratingTo, targetFrom, targetTo)
	s.CreatedAt = createdAt
	return s
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
