package unit_test

import (
	"context"
	"errors"
	"testing"

	"github.com/geomena/stock-recommendation-system/backend/internal/domain"
)

func TestGetTopRecommendations_LimitNormalization(t *testing.T) {
	tests := []struct {
		name  string
		limit int
	}{
		{"zero limit uses default", 0},
		{"negative limit uses default", -1},
		{"over 100 uses default", 150},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock := newMockRepo()
			mock.FindAllFn = func(ctx context.Context, filter domain.StockFilter) ([]domain.Stock, int64, error) {
				return []domain.Stock{}, 0, nil
			}

			uc := newRecommendationUsecase(mock)
			result, err := uc.GetTopRecommendations(context.Background(), tc.limit, "")
			assertNoError(t, err)

			if result == nil {
				t.Fatal("expected non-nil result")
			}
		})
	}
}

func TestGetTopRecommendations_RanksByScore(t *testing.T) {
	mock := newMockRepo()
	mock.FindAllFn = func(ctx context.Context, filter domain.StockFilter) ([]domain.Stock, int64, error) {
		stocks := []domain.Stock{
			// AAPL: upgraded hold→buy con target increase = debería puntuar alto
			makeStock(stockID1, "AAPL", "Apple Inc.", "Morgan Stanley", "upgraded", "hold", "buy", 180.0, 220.0),
			// GOOGL: downgraded buy→sell = debería puntuar bajo
			makeStock(stockID2, "GOOGL", "Alphabet Inc.", "Goldman Sachs", "downgraded", "buy", "sell", 185.0, 140.0),
		}
		return stocks, int64(len(stocks)), nil
	}

	uc := newRecommendationUsecase(mock)
	result, err := uc.GetTopRecommendations(context.Background(), 10, "")
	assertNoError(t, err)

	if len(result) < 1 {
		t.Fatal("expected at least 1 recommendation")
	}

	// AAPL debería rankear primero (upgraded + target increase + bullish action)
	if result[0].Stock.Ticker != "AAPL" {
		t.Errorf("expected AAPL ranked first, got %s", result[0].Stock.Ticker)
	}
	if result[0].Score <= 0 {
		t.Errorf("expected positive score for AAPL, got %f", result[0].Score)
	}
}

func TestGetTopRecommendations_MultipleTickersGrouped(t *testing.T) {
	mock := newMockRepo()
	mock.FindAllFn = func(ctx context.Context, filter domain.StockFilter) ([]domain.Stock, int64, error) {
		stocks := []domain.Stock{
			// Dos analistas bullish en AAPL → consenso alto
			makeStock(stockID1, "AAPL", "Apple Inc.", "Morgan Stanley", "upgraded", "hold", "buy", 180.0, 220.0),
			makeStock(stockID2, "AAPL", "Apple Inc.", "Goldman Sachs", "target raised", "buy", "buy", 200.0, 230.0),
			makeStock(stockID3, "AAPL", "Apple Inc.", "JP Morgan", "initiated", "neutral", "buy", 190.0, 225.0),
			// Un analista bearish en GOOGL
			makeStock(stockID4, "GOOGL", "Alphabet Inc.", "Barclays", "downgraded", "buy", "hold", 185.0, 150.0),
		}
		return stocks, int64(len(stocks)), nil
	}

	uc := newRecommendationUsecase(mock)
	result, err := uc.GetTopRecommendations(context.Background(), 10, "")
	assertNoError(t, err)

	// Deberían haber 2 recomendaciones (una por ticker) o al menos AAPL
	aaplFound := false
	for _, rec := range result {
		if rec.Stock.Ticker == "AAPL" {
			aaplFound = true
			if rec.AnalystCount != 3 {
				t.Errorf("expected 3 analysts for AAPL, got %d", rec.AnalystCount)
			}
			if len(rec.Reasons) == 0 {
				t.Error("expected reasons for AAPL recommendation")
			}
		}
	}

	if !aaplFound {
		t.Error("expected AAPL in recommendations")
	}
}

func TestGetTopRecommendations_EmptyStocks(t *testing.T) {
	mock := newMockRepo()
	mock.FindAllFn = func(ctx context.Context, filter domain.StockFilter) ([]domain.Stock, int64, error) {
		return []domain.Stock{}, 0, nil
	}

	uc := newRecommendationUsecase(mock)
	result, err := uc.GetTopRecommendations(context.Background(), 10, "")
	assertNoError(t, err)

	if result == nil {
		t.Fatal("expected non-nil empty slice, got nil")
	}
	if len(result) != 0 {
		t.Errorf("expected 0 recommendations, got %d", len(result))
	}
}

func TestGetTopRecommendations_RepoError(t *testing.T) {
	mock := newMockRepo()
	mock.FindAllFn = func(ctx context.Context, filter domain.StockFilter) ([]domain.Stock, int64, error) {
		return nil, 0, errors.New("database unavailable")
	}

	uc := newRecommendationUsecase(mock)
	result, err := uc.GetTopRecommendations(context.Background(), 10, "")
	assertError(t, err)

	if result != nil {
		t.Errorf("expected nil result on error, got %+v", result)
	}
}

func TestGetTopRecommendations_LimitApplied(t *testing.T) {
	mock := newMockRepo()
	mock.FindAllFn = func(ctx context.Context, filter domain.StockFilter) ([]domain.Stock, int64, error) {
		stocks := []domain.Stock{
			makeStock(stockID1, "AAPL", "Apple Inc.", "Morgan Stanley", "upgraded", "hold", "buy", 180.0, 220.0),
			makeStock(stockID2, "GOOGL", "Alphabet Inc.", "Goldman Sachs", "upgraded", "hold", "buy", 140.0, 185.0),
			makeStock(stockID3, "MSFT", "Microsoft Corp.", "JP Morgan", "upgraded", "hold", "buy", 380.0, 420.0),
			makeStock(stockID4, "AMZN", "Amazon.com Inc.", "Barclays", "upgraded", "hold", "buy", 175.0, 210.0),
			makeStock(stockID5, "TSLA", "Tesla Inc.", "Wedbush", "upgraded", "hold", "buy", 250.0, 300.0),
		}
		return stocks, int64(len(stocks)), nil
	}

	uc := newRecommendationUsecase(mock)
	result, err := uc.GetTopRecommendations(context.Background(), 2, "")
	assertNoError(t, err)

	if len(result) > 2 {
		t.Errorf("expected at most 2 recommendations, got %d", len(result))
	}
}

func TestGetTopRecommendations_SearchPassedToRepo(t *testing.T) {
	mock := newMockRepo()
	var receivedSearch string
	mock.FindAllFn = func(ctx context.Context, filter domain.StockFilter) ([]domain.Stock, int64, error) {
		receivedSearch = filter.Search
		return []domain.Stock{}, 0, nil
	}

	uc := newRecommendationUsecase(mock)
	_, err := uc.GetTopRecommendations(context.Background(), 10, "AAPL")
	assertNoError(t, err)

	if receivedSearch != "AAPL" {
		t.Errorf("expected search 'AAPL' passed to repo, got %q", receivedSearch)
	}
}

func TestGetBestStock_ReturnsTopOne(t *testing.T) {
	mock := newMockRepo()
	mock.FindAllFn = func(ctx context.Context, filter domain.StockFilter) ([]domain.Stock, int64, error) {
		stocks := []domain.Stock{
			makeStock(stockID1, "AAPL", "Apple Inc.", "Morgan Stanley", "upgraded", "hold", "buy", 180.0, 220.0),
			makeStock(stockID2, "GOOGL", "Alphabet Inc.", "Goldman Sachs", "downgraded", "buy", "sell", 185.0, 140.0),
		}
		return stocks, int64(len(stocks)), nil
	}

	uc := newRecommendationUsecase(mock)
	result, err := uc.GetBestStock(context.Background())
	assertNoError(t, err)

	if result == nil {
		t.Fatal("expected non-nil recommendation")
	}
	if result.Score <= 0 {
		t.Errorf("expected positive score, got %f", result.Score)
	}
}

func TestGetBestStock_NoStocks(t *testing.T) {
	mock := newMockRepo()
	mock.FindAllFn = func(ctx context.Context, filter domain.StockFilter) ([]domain.Stock, int64, error) {
		return []domain.Stock{}, 0, nil
	}

	uc := newRecommendationUsecase(mock)
	result, err := uc.GetBestStock(context.Background())
	assertNoError(t, err)

	if result != nil {
		t.Errorf("expected nil for no stocks, got %+v", result)
	}
}

func TestGetTopRecommendations_ScoreIsPositive(t *testing.T) {
	mock := newMockRepo()
	mock.FindAllFn = func(ctx context.Context, filter domain.StockFilter) ([]domain.Stock, int64, error) {
		stocks := []domain.Stock{
			makeStock(stockID1, "AAPL", "Apple Inc.", "Morgan Stanley", "upgraded", "hold", "strong buy", 180.0, 250.0),
		}
		return stocks, int64(len(stocks)), nil
	}

	uc := newRecommendationUsecase(mock)
	result, err := uc.GetTopRecommendations(context.Background(), 10, "")
	assertNoError(t, err)

	if len(result) != 1 {
		t.Fatalf("expected 1 recommendation, got %d", len(result))
	}

	if result[0].Score <= 0 {
		t.Errorf("expected positive score for upgraded stock, got %f", result[0].Score)
	}
	if result[0].UpsidePotential <= 0 {
		t.Errorf("expected positive upside for stock with target increase, got %f", result[0].UpsidePotential)
	}
}
