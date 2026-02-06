package feature_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/geomena/stock-recommendation-system/backend/internal/domain"
	"github.com/geomena/stock-recommendation-system/backend/internal/i18n/en"
)

func TestGetDashboardStats_Success(t *testing.T) {
	app := newTestApp()

	app.mockRepo.CountAllFn = func(ctx context.Context) (int64, error) {
		return 150, nil
	}
	app.mockRepo.GetActionDistributionFn = func(ctx context.Context) ([]domain.ActionDistribution, error) {
		return sampleActionDistributions(), nil
	}
	app.mockRepo.GetBrokerageDistributionFn = func(ctx context.Context, limit int) ([]domain.BrokerageDistribution, error) {
		return sampleBrokerageDistributions(), nil
	}
	app.mockRepo.GetRecentActivityFn = func(ctx context.Context, days int) ([]domain.DailyActivity, error) {
		return sampleDailyActivities(), nil
	}

	rec, resp := doRequest(t, app.router, http.MethodGet, "/api/v1/dashboard/stats")

	assertStatus(t, rec, http.StatusOK)
	assertContentType(t, rec)
	assertSuccess(t, resp)
	assertHasData(t, resp)

	if resp.Message != en.DashboardStatsRetrieved {
		t.Errorf("expected message %q, got %q", en.DashboardStatsRetrieved, resp.Message)
	}

	var stats domain.DashboardStats
	if err := json.Unmarshal(resp.Data, &stats); err != nil {
		t.Fatalf("failed to unmarshal data: %v", err)
	}
	if stats.TotalStocks != 150 {
		t.Errorf("expected totalStocks 150, got %d", stats.TotalStocks)
	}
	if len(stats.ActionDistribution) != 5 {
		t.Errorf("expected 5 action distributions, got %d", len(stats.ActionDistribution))
	}
	if len(stats.BrokerageDistribution) != 4 {
		t.Errorf("expected 4 brokerage distributions, got %d", len(stats.BrokerageDistribution))
	}
	if len(stats.RecentActivity) != 3 {
		t.Errorf("expected 3 daily activities, got %d", len(stats.RecentActivity))
	}
}

func TestGetDashboardStats_InternalError(t *testing.T) {
	app := newTestApp()

	app.mockRepo.CountAllFn = func(ctx context.Context) (int64, error) {
		return 0, errors.New("database unavailable")
	}

	rec, resp := doRequest(t, app.router, http.MethodGet, "/api/v1/dashboard/stats")

	assertStatus(t, rec, http.StatusInternalServerError)
	assertContentType(t, rec)
	assertError(t, resp)
	assertNoData(t, resp)
}
