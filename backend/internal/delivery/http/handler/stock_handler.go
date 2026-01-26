package handler

import (
	"net/http"
	"strconv"

	"github.com/geomena/stock-recommendation-system/backend/internal/domain"
	"github.com/geomena/stock-recommendation-system/backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StockHandler struct {
	stockUsecase          *usecase.StockUsecase
	recommendationUsecase *usecase.RecommendationUsecase
}

func NewStockHandler(su *usecase.StockUsecase, ru *usecase.RecommendationUsecase) *StockHandler {
	return &StockHandler{
		stockUsecase:          su,
		recommendationUsecase: ru,
	}
}

func (h *StockHandler) ListStocks(c *gin.Context) {
	filter := domain.NewStockFilter()

	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil {
		filter.Page = page
	}
	if limit, err := strconv.Atoi(c.DefaultQuery("limit", "20")); err == nil {
		filter.Limit = limit
	}

	filter.Search = c.Query("search")
	filter.Ticker = c.Query("ticker")
	filter.Action = c.Query("action")
	filter.SortBy = c.DefaultQuery("sortBy", "created_at")
	filter.SortOrder = c.DefaultQuery("sortOrder", "desc")

	result, err := h.stockUsecase.ListStocks(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *StockHandler) GetStock(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stock ID"})
		return
	}

	stock, err := h.stockUsecase.GetStockByID(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrStockNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "stock not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stock)
}

func (h *StockHandler) GetByTicker(c *gin.Context) {
	ticker := c.Param("ticker")
	if ticker == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticker is required"})
		return
	}

	stocks, err := h.stockUsecase.GetStocksByTicker(c.Request.Context(), ticker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stocks})
}

func (h *StockHandler) GetActions(c *gin.Context) {
	actions, err := h.stockUsecase.GetDistinctActions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": actions})
}

func (h *StockHandler) SyncStocks(c *gin.Context) {
	count, err := h.stockUsecase.SyncFromExternalAPI(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "sync completed",
		"count":   count,
	})
}

func (h *StockHandler) GetRecommendations(c *gin.Context) {
	limit := 10
	if l, err := strconv.Atoi(c.DefaultQuery("limit", "10")); err == nil {
		limit = l
	}

	recommendations, err := h.recommendationUsecase.GetTopRecommendations(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": recommendations})
}

func (h *StockHandler) GetTopRecommendation(c *gin.Context) {
	recommendation, err := h.recommendationUsecase.GetBestStock(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if recommendation == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no recommendations available"})
		return
	}

	c.JSON(http.StatusOK, recommendation)
}
