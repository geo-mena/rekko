package feature_test

import (
	"net/http"
	"testing"

	"github.com/geomena/stock-recommendation-system/backend/internal/i18n/en"
)

func TestHealth_Success(t *testing.T) {
	app := newTestApp()

	rec, resp := doRequest(t, app.router, http.MethodGet, "/api/v1/health")

	assertStatus(t, rec, http.StatusOK)
	assertContentType(t, rec)
	assertSuccess(t, resp)

	if resp.Message != en.ServiceRunning {
		t.Errorf("expected message %q, got %q", en.ServiceRunning, resp.Message)
	}
}
