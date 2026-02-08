package handler

import (
	"net/http"

	"github.com/geomena/stock-recommendation-system/backend/internal/delivery/http/response"
	"github.com/geomena/stock-recommendation-system/backend/internal/i18n/en"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health godoc
//
//	@Summary	Health check
//	@Description	Returns the current health status of the service
//	@Tags			Health
//	@Produce		json
//	@Success		200	{object}	APIResponse	"Service is running"
//	@Router			/health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	response.MessageOnly(c.Writer, http.StatusOK, en.ServiceRunning)
}
