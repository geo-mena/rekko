package handler

import (
	"net/http"

	"github.com/geomena/stock-recommendation-system/backend/internal/usecase"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardUsecase *usecase.DashboardUsecase
}

func NewDashboardHandler(du *usecase.DashboardUsecase) *DashboardHandler {
	return &DashboardHandler{dashboardUsecase: du}
}

func (h *DashboardHandler) GetStats(c *gin.Context) {
	stats, err := h.dashboardUsecase.GetDashboardStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
