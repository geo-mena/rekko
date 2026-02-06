package feature_test

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	httpdelivery "github.com/geomena/stock-recommendation-system/backend/internal/delivery/http"
	"github.com/geomena/stock-recommendation-system/backend/internal/delivery/http/handler"
	"github.com/geomena/stock-recommendation-system/backend/internal/domain"
	"github.com/geomena/stock-recommendation-system/backend/internal/repository"
	"github.com/geomena/stock-recommendation-system/backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type testApp struct {
	router   *gin.Engine
	mockRepo *repository.MockStockRepository
}

func newTestApp() *testApp {
	mockRepo := &repository.MockStockRepository{}

	stockUsecase := usecase.NewStockUsecase(mockRepo, nil)
	recommendationUsecase := usecase.NewRecommendationUsecase(mockRepo)
	dashboardUsecase := usecase.NewDashboardUsecase(mockRepo)

	stockHandler := handler.NewStockHandler(stockUsecase, recommendationUsecase)
	healthHandler := handler.NewHealthHandler()
	dashboardHandler := handler.NewDashboardHandler(dashboardUsecase)

	router := httpdelivery.NewRouter(stockHandler, healthHandler, dashboardHandler)

	return &testApp{
		router:   router,
		mockRepo: mockRepo,
	}
}

type jsonResponse struct {
	Status  bool            `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
	Meta    *jsonMeta       `json:"meta,omitempty"`
}

type jsonMeta struct {
	Pagination *jsonPagination `json:"pagination,omitempty"`
}

type jsonPagination struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
}

func doRequest(t *testing.T, router *gin.Engine, method, path string) (*httptest.ResponseRecorder, jsonResponse) {
	t.Helper()

	req := httptest.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	var resp jsonResponse
	body, err := io.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("failed to unmarshal response body: %v\nbody: %s", err, string(body))
	}

	return rec, resp
}

func assertStatus(t *testing.T, rec *httptest.ResponseRecorder, expected int) {
	t.Helper()
	if rec.Code != expected {
		t.Errorf("expected status %d, got %d", expected, rec.Code)
	}
}

func assertContentType(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()
	ct := rec.Header().Get("Content-Type")
	expected := "application/json; charset=utf-8"
	if ct != expected {
		t.Errorf("expected Content-Type %q, got %q", expected, ct)
	}
}

func assertSuccess(t *testing.T, resp jsonResponse) {
	t.Helper()
	if !resp.Status {
		t.Errorf("expected status true, got false")
	}
	if resp.Message == "" {
		t.Errorf("expected non-empty message")
	}
}

func assertError(t *testing.T, resp jsonResponse) {
	t.Helper()
	if resp.Status {
		t.Errorf("expected status false, got true")
	}
	if resp.Message == "" {
		t.Errorf("expected non-empty message")
	}
}

func assertHasData(t *testing.T, resp jsonResponse) {
	t.Helper()
	if resp.Data == nil || string(resp.Data) == "" || string(resp.Data) == "null" {
		t.Errorf("expected data to be present")
	}
}

func assertNoData(t *testing.T, resp jsonResponse) {
	t.Helper()
	if resp.Data != nil && string(resp.Data) != "" && string(resp.Data) != "null" {
		t.Errorf("expected data to be absent, got: %s", string(resp.Data))
	}
}

func assertPagination(t *testing.T, resp jsonResponse, page, perPage int, totalItems int64, totalPages int, hasNext bool) {
	t.Helper()
	if resp.Meta == nil {
		t.Fatal("expected meta to be present")
	}
	if resp.Meta.Pagination == nil {
		t.Fatal("expected meta.pagination to be present")
	}
	p := resp.Meta.Pagination
	if p.CurrentPage != page {
		t.Errorf("expected current_page %d, got %d", page, p.CurrentPage)
	}
	if p.PerPage != perPage {
		t.Errorf("expected per_page %d, got %d", perPage, p.PerPage)
	}
	if p.TotalItems != totalItems {
		t.Errorf("expected total_items %d, got %d", totalItems, p.TotalItems)
	}
	if p.TotalPages != totalPages {
		t.Errorf("expected total_pages %d, got %d", totalPages, p.TotalPages)
	}
	if p.HasNext != hasNext {
		t.Errorf("expected has_next %v, got %v", hasNext, p.HasNext)
	}
}

var (
	stockIDApple  = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	stockIDGoogle = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	stockIDMSFT   = uuid.MustParse("33333333-3333-3333-3333-333333333333")
	stockIDAMZN   = uuid.MustParse("44444444-4444-4444-4444-444444444444")
	stockIDTSLA   = uuid.MustParse("55555555-5555-5555-5555-555555555555")

	now = time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
)

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
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func sampleStocks() []domain.Stock {
	return []domain.Stock{
		makeStock(stockIDApple, "AAPL", "Apple Inc.", "Morgan Stanley", "upgraded", "Hold", "Buy", 180.0, 220.0),
		makeStock(stockIDGoogle, "GOOGL", "Alphabet Inc.", "Goldman Sachs", "initiated", "Neutral", "Buy", 140.0, 185.0),
		makeStock(stockIDMSFT, "MSFT", "Microsoft Corp.", "JP Morgan", "reiterated", "Buy", "Buy", 380.0, 420.0),
		makeStock(stockIDAMZN, "AMZN", "Amazon.com Inc.", "Barclays", "target raised", "Overweight", "Overweight", 175.0, 210.0),
		makeStock(stockIDTSLA, "TSLA", "Tesla Inc.", "Wedbush", "maintained", "Outperform", "Outperform", 250.0, 300.0),
	}
}

func sampleActionDistributions() []domain.ActionDistribution {
	return []domain.ActionDistribution{
		{Action: "upgraded", Count: 15},
		{Action: "initiated", Count: 12},
		{Action: "reiterated", Count: 25},
		{Action: "target raised", Count: 8},
		{Action: "maintained", Count: 20},
	}
}

func sampleBrokerageDistributions() []domain.BrokerageDistribution {
	return []domain.BrokerageDistribution{
		{Brokerage: "Morgan Stanley", Count: 18},
		{Brokerage: "Goldman Sachs", Count: 14},
		{Brokerage: "JP Morgan", Count: 22},
		{Brokerage: "Barclays", Count: 10},
	}
}

func sampleDailyActivities() []domain.DailyActivity {
	return []domain.DailyActivity{
		{Date: "2025-01-15", Count: 5},
		{Date: "2025-01-14", Count: 8},
		{Date: "2025-01-13", Count: 3},
	}
}
