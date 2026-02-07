package handler

import (
	"net/http"
	"strconv"

	"github.com/geomena/stock-recommendation-system/backend/internal/delivery/http/response"
	"github.com/geomena/stock-recommendation-system/backend/internal/domain"
	"github.com/geomena/stock-recommendation-system/backend/internal/i18n/en"
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
		response.InternalServerError(c.Writer, err.Error())
		return
	}

	response.SuccessWithPagination(c.Writer, http.StatusOK, en.StocksRetrieved, result.Data, response.PaginationParams{
		Page:    result.Page,
		PerPage: result.Limit,
		Total:   result.TotalCount,
	})
}

func (h *StockHandler) GetStock(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c.Writer, en.StockInvalidID)
		return
	}

	stock, err := h.stockUsecase.GetStockByID(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrStockNotFound {
			response.NotFound(c.Writer, en.StockNotFound)
			return
		}
		response.InternalServerError(c.Writer, err.Error())
		return
	}

	response.Success(c.Writer, http.StatusOK, en.StockRetrieved, stock)
}

func (h *StockHandler) GetByTicker(c *gin.Context) {
	ticker := c.Param("ticker")
	if ticker == "" {
		response.BadRequest(c.Writer, en.StockTickerRequired)
		return
	}

	stocks, err := h.stockUsecase.GetStocksByTicker(c.Request.Context(), ticker)
	if err != nil {
		response.InternalServerError(c.Writer, err.Error())
		return
	}

	response.Success(c.Writer, http.StatusOK, en.StocksRetrieved, stocks)
}

func (h *StockHandler) GetActions(c *gin.Context) {
	actions, err := h.stockUsecase.GetDistinctActions(c.Request.Context())
	if err != nil {
		response.InternalServerError(c.Writer, err.Error())
		return
	}

	response.Success(c.Writer, http.StatusOK, en.ActionsRetrieved, actions)
}

func (h *StockHandler) SyncStocks(c *gin.Context) {
	count, err := h.stockUsecase.SyncFromExternalAPI(c.Request.Context())
	if err != nil {
		response.InternalServerError(c.Writer, err.Error())
		return
	}

	response.Success(c.Writer, http.StatusOK, en.SyncCompleted, gin.H{"count": count})
}

func (h *StockHandler) GetRecommendations(c *gin.Context) {
	limit := 50
	if l, err := strconv.Atoi(c.DefaultQuery("limit", "50")); err == nil {
		limit = l
	}

	search := c.Query("search")

	recommendations, err := h.recommendationUsecase.GetTopRecommendations(c.Request.Context(), limit, search)
	if err != nil {
		response.InternalServerError(c.Writer, err.Error())
		return
	}

	response.Success(c.Writer, http.StatusOK, en.RecommendationsRetrieved, recommendations)
}

func (h *StockHandler) GetTopRecommendation(c *gin.Context) {
	recommendation, err := h.recommendationUsecase.GetBestStock(c.Request.Context())
	if err != nil {
		response.InternalServerError(c.Writer, err.Error())
		return
	}

	if recommendation == nil {
		response.NotFound(c.Writer, en.NoRecommendationsAvailable)
		return
	}

	response.Success(c.Writer, http.StatusOK, en.TopRecommendationRetrieved, recommendation)
}
