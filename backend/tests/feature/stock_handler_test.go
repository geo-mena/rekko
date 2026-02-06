package feature_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/geomena/stock-recommendation-system/backend/internal/domain"
	"github.com/geomena/stock-recommendation-system/backend/internal/i18n/en"
	"github.com/google/uuid"
)

func TestListStocks_Success(t *testing.T) {
	app := newTestApp()
	stocks := sampleStocks()

	app.mockRepo.FindAllFn = func(ctx context.Context, filter domain.StockFilter) ([]domain.Stock, int64, error) {
		return stocks, 42, nil
	}

	rec, resp := doRequest(t, app.router, http.MethodGet, "/api/v1/stocks?page=1&limit=20")

	assertStatus(t, rec, http.StatusOK)
	assertContentType(t, rec)
	assertSuccess(t, resp)
	assertHasData(t, resp)
	assertPagination(t, resp, 1, 20, 42, 3, true)

	if resp.Message != en.StocksRetrieved {
		t.Errorf("expected message %q, got %q", en.StocksRetrieved, resp.Message)
	}

	var data []domain.Stock
	if err := json.Unmarshal(resp.Data, &data); err != nil {
		t.Fatalf("failed to unmarshal data: %v", err)
	}
	if len(data) != len(stocks) {
		t.Errorf("expected %d stocks, got %d", len(stocks), len(data))
	}
}

func TestListStocks_PaginationMeta(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		total      int64
		page       int
		perPage    int
		totalPages int
		hasNext    bool
	}{
		{
			name:       "first page with more pages",
			path:       "/api/v1/stocks?page=1&limit=2",
			total:      5,
			page:       1,
			perPage:    2,
			totalPages: 3,
			hasNext:    true,
		},
		{
			name:       "last page",
			path:       "/api/v1/stocks?page=3&limit=2",
			total:      5,
			page:       3,
			perPage:    2,
			totalPages: 3,
			hasNext:    false,
		},
		{
			name:       "single page",
			path:       "/api/v1/stocks?page=1&limit=20",
			total:      3,
			page:       1,
			perPage:    20,
			totalPages: 1,
			hasNext:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app := newTestApp()
			app.mockRepo.FindAllFn = func(ctx context.Context, filter domain.StockFilter) ([]domain.Stock, int64, error) {
				return sampleStocks()[:2], tc.total, nil
			}

			rec, resp := doRequest(t, app.router, http.MethodGet, tc.path)

			assertStatus(t, rec, http.StatusOK)
			assertContentType(t, rec)
			assertSuccess(t, resp)
			assertPagination(t, resp, tc.page, tc.perPage, tc.total, tc.totalPages, tc.hasNext)
		})
	}
}

func TestGetStock_Success(t *testing.T) {
	app := newTestApp()
	expected := sampleStocks()[0]

	app.mockRepo.FindByIDFn = func(ctx context.Context, id uuid.UUID) (*domain.Stock, error) {
		if id == stockIDApple {
			return &expected, nil
		}
		return nil, domain.ErrStockNotFound
	}

	rec, resp := doRequest(t, app.router, http.MethodGet, "/api/v1/stocks/"+stockIDApple.String())

	assertStatus(t, rec, http.StatusOK)
	assertContentType(t, rec)
	assertSuccess(t, resp)
	assertHasData(t, resp)

	if resp.Message != en.StockRetrieved {
		t.Errorf("expected message %q, got %q", en.StockRetrieved, resp.Message)
	}

	var stock domain.Stock
	if err := json.Unmarshal(resp.Data, &stock); err != nil {
		t.Fatalf("failed to unmarshal data: %v", err)
	}
	if stock.Ticker != "AAPL" {
		t.Errorf("expected ticker AAPL, got %s", stock.Ticker)
	}
	if stock.Company != "Apple Inc." {
		t.Errorf("expected company Apple Inc., got %s", stock.Company)
	}
}

func TestGetStock_InvalidUUID(t *testing.T) {
	app := newTestApp()

	rec, resp := doRequest(t, app.router, http.MethodGet, "/api/v1/stocks/not-a-uuid")

	assertStatus(t, rec, http.StatusBadRequest)
	assertContentType(t, rec)
	assertError(t, resp)
	assertNoData(t, resp)

	if resp.Message != en.StockInvalidID {
		t.Errorf("expected message %q, got %q", en.StockInvalidID, resp.Message)
	}
}

func TestGetStock_NotFound(t *testing.T) {
	app := newTestApp()

	app.mockRepo.FindByIDFn = func(ctx context.Context, id uuid.UUID) (*domain.Stock, error) {
		return nil, domain.ErrStockNotFound
	}

	missingID := uuid.MustParse("99999999-9999-9999-9999-999999999999")
	rec, resp := doRequest(t, app.router, http.MethodGet, "/api/v1/stocks/"+missingID.String())

	assertStatus(t, rec, http.StatusNotFound)
	assertContentType(t, rec)
	assertError(t, resp)
	assertNoData(t, resp)

	if resp.Message != en.StockNotFound {
		t.Errorf("expected message %q, got %q", en.StockNotFound, resp.Message)
	}
}

func TestGetStock_InternalError(t *testing.T) {
	app := newTestApp()

	app.mockRepo.FindByIDFn = func(ctx context.Context, id uuid.UUID) (*domain.Stock, error) {
		return nil, errors.New("database connection lost")
	}

	rec, resp := doRequest(t, app.router, http.MethodGet, "/api/v1/stocks/"+stockIDApple.String())

	assertStatus(t, rec, http.StatusInternalServerError)
	assertContentType(t, rec)
	assertError(t, resp)
	assertNoData(t, resp)
}

func TestGetByTicker_Success(t *testing.T) {
	app := newTestApp()
	appleStock := sampleStocks()[0]

	app.mockRepo.FindByTickerFn = func(ctx context.Context, ticker string) ([]domain.Stock, error) {
		if ticker == "AAPL" {
			return []domain.Stock{appleStock}, nil
		}
		return []domain.Stock{}, nil
	}

	rec, resp := doRequest(t, app.router, http.MethodGet, "/api/v1/stocks/ticker/AAPL")

	assertStatus(t, rec, http.StatusOK)
	assertContentType(t, rec)
	assertSuccess(t, resp)
	assertHasData(t, resp)

	if resp.Message != en.StocksRetrieved {
		t.Errorf("expected message %q, got %q", en.StocksRetrieved, resp.Message)
	}

	var stocks []domain.Stock
	if err := json.Unmarshal(resp.Data, &stocks); err != nil {
		t.Fatalf("failed to unmarshal data: %v", err)
	}
	if len(stocks) != 1 {
		t.Fatalf("expected 1 stock, got %d", len(stocks))
	}
	if stocks[0].Ticker != "AAPL" {
		t.Errorf("expected ticker AAPL, got %s", stocks[0].Ticker)
	}
}

func TestGetActions_Success(t *testing.T) {
	app := newTestApp()
	expectedActions := []string{"upgraded", "initiated", "reiterated", "target raised", "maintained"}

	app.mockRepo.GetDistinctActionsFn = func(ctx context.Context) ([]string, error) {
		return expectedActions, nil
	}

	rec, resp := doRequest(t, app.router, http.MethodGet, "/api/v1/stocks/actions")

	assertStatus(t, rec, http.StatusOK)
	assertContentType(t, rec)
	assertSuccess(t, resp)
	assertHasData(t, resp)

	if resp.Message != en.ActionsRetrieved {
		t.Errorf("expected message %q, got %q", en.ActionsRetrieved, resp.Message)
	}

	var actions []string
	if err := json.Unmarshal(resp.Data, &actions); err != nil {
		t.Fatalf("failed to unmarshal data: %v", err)
	}
	if len(actions) != len(expectedActions) {
		t.Errorf("expected %d actions, got %d", len(expectedActions), len(actions))
	}
}

func TestGetRecommendations_Success(t *testing.T) {
	app := newTestApp()
	stocks := sampleStocks()

	app.mockRepo.FindAllFn = func(ctx context.Context, filter domain.StockFilter) ([]domain.Stock, int64, error) {
		return stocks, int64(len(stocks)), nil
	}

	rec, resp := doRequest(t, app.router, http.MethodGet, "/api/v1/recommendations?limit=5")

	assertStatus(t, rec, http.StatusOK)
	assertContentType(t, rec)
	assertSuccess(t, resp)
	assertHasData(t, resp)

	if resp.Message != en.RecommendationsRetrieved {
		t.Errorf("expected message %q, got %q", en.RecommendationsRetrieved, resp.Message)
	}

	var recommendations []domain.StockRecommendation
	if err := json.Unmarshal(resp.Data, &recommendations); err != nil {
		t.Fatalf("failed to unmarshal data: %v", err)
	}
	if len(recommendations) == 0 {
		t.Error("expected at least one recommendation")
	}
	for i := 1; i < len(recommendations); i++ {
		if recommendations[i].Score > recommendations[i-1].Score {
			t.Errorf("recommendations not sorted by score: index %d (%.2f) > index %d (%.2f)",
				i, recommendations[i].Score, i-1, recommendations[i-1].Score)
		}
	}
}

func TestGetTopRecommendation_Success(t *testing.T) {
	app := newTestApp()
	stocks := sampleStocks()

	app.mockRepo.FindAllFn = func(ctx context.Context, filter domain.StockFilter) ([]domain.Stock, int64, error) {
		return stocks, int64(len(stocks)), nil
	}

	rec, resp := doRequest(t, app.router, http.MethodGet, "/api/v1/recommendations/top")

	assertStatus(t, rec, http.StatusOK)
	assertContentType(t, rec)
	assertSuccess(t, resp)
	assertHasData(t, resp)

	if resp.Message != en.TopRecommendationRetrieved {
		t.Errorf("expected message %q, got %q", en.TopRecommendationRetrieved, resp.Message)
	}

	var recommendation domain.StockRecommendation
	if err := json.Unmarshal(resp.Data, &recommendation); err != nil {
		t.Fatalf("failed to unmarshal data: %v", err)
	}
	if recommendation.Score <= 0 {
		t.Error("expected positive score for top recommendation")
	}
	if recommendation.Stock.Ticker == "" {
		t.Error("expected non-empty ticker in top recommendation")
	}
}

func TestGetTopRecommendation_NoRecommendations(t *testing.T) {
	app := newTestApp()

	app.mockRepo.FindAllFn = func(ctx context.Context, filter domain.StockFilter) ([]domain.Stock, int64, error) {
		return []domain.Stock{}, 0, nil
	}

	rec, resp := doRequest(t, app.router, http.MethodGet, "/api/v1/recommendations/top")

	assertStatus(t, rec, http.StatusNotFound)
	assertContentType(t, rec)
	assertError(t, resp)
	assertNoData(t, resp)

	if resp.Message != en.NoRecommendationsAvailable {
		t.Errorf("expected message %q, got %q", en.NoRecommendationsAvailable, resp.Message)
	}
}
