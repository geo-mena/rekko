package handler

import (
	"net/http"

	"github.com/geomena/stock-recommendation-system/backend/internal/delivery/http/response"
	"github.com/geomena/stock-recommendation-system/backend/internal/i18n/en"
	"github.com/geomena/stock-recommendation-system/backend/internal/usecase"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardUsecase *usecase.DashboardUsecase
}

func NewDashboardHandler(du *usecase.DashboardUsecase) *DashboardHandler {
	return &DashboardHandler{dashboardUsecase: du}
}

// GetStats godoc
//
//	@Summary	Get dashboard statistics
//	@Description	Returns aggregated statistics including total stocks, action distribution, top brokerages, and recent daily activity
//	@Tags			Dashboard
//	@Produce		json
//	@Success		200	{object}	APIResponse{data=DashboardStats}	"Dashboard stats retrieved successfully"
//	@Failure		500	{object}	APIResponse								"Internal server error"
//	@Router			/dashboard/stats [get]
func (h *DashboardHandler) GetStats(c *gin.Context) {
	stats, err := h.dashboardUsecase.GetDashboardStats(c.Request.Context())
	if err != nil {
		response.InternalServerError(c.Writer, err.Error())
		return
	}
	response.Success(c.Writer, http.StatusOK, en.DashboardStatsRetrieved, stats)
}
