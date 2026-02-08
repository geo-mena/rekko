package handler

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/geomena/stock-recommendation-system/backend/internal/delivery/http/response"
	"github.com/geomena/stock-recommendation-system/backend/internal/i18n/en"
	"github.com/geomena/stock-recommendation-system/backend/web"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	healthTemplate *template.Template
}

func NewHealthHandler() *HealthHandler {
	tmpl, err := template.ParseFS(web.TemplatesFS, "templates/health.html")
	if err != nil {
		panic("failed to parse health template: " + err.Error())
	}
	return &HealthHandler{healthTemplate: tmpl}
}

type healthPageData struct {
	Message   string
	Timestamp string
}

// Health godoc
//
//	@Summary	Health check
//	@Description	Returns the current health status of the service
//	@Tags			Health
//	@Produce		json
//	@Produce		html
//	@Success		200	{object}	APIResponse	"Service is running"
//	@Router			/health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	if wantsHTML(c.GetHeader("Accept")) {
		h.renderHealthPage(c.Writer)
		return
	}
	response.MessageOnly(c.Writer, http.StatusOK, en.ServiceRunning)
}

func wantsHTML(accept string) bool {
	return strings.Contains(accept, "text/html")
}

func (h *HealthHandler) renderHealthPage(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := healthPageData{
		Message:   en.ServiceRunning,
		Timestamp: time.Now().UTC().Format("2006-01-02 15:04:05"),
	}
	if err := h.healthTemplate.Execute(w, data); err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
	}
}
