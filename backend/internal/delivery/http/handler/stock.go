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

// ListStocks godoc
//
//	@Summary	List stocks
//	@Description	Returns a paginated list of stocks with optional filtering and sorting
//	@Tags			Stocks
//	@Produce		json
//	@Param			page		query		int		false	"Page number"			default(1)
//	@Param			limit		query		int		false	"Items per page"		default(20)
//	@Param			search		query		string	false	"Search in ticker and company name"
//	@Param			ticker		query		string	false	"Filter by ticker symbol"
//	@Param			action		query		string	false	"Filter by action (e.g. upgraded, downgraded)"
//	@Param			sortBy		query		string	false	"Sort field"			default(created_at)
//	@Param			sortOrder	query		string	false	"Sort direction"		default(desc)	Enums(asc, desc)
//	@Success		200			{object}	APIResponse{data=[]Stock,meta=PaginationMeta}	"Stocks retrieved successfully"
//	@Failure		500			{object}	APIResponse	"Internal server error"
//	@Router			/stocks [get]
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

// GetStock godoc
//
//	@Summary	Get stock by ID
//	@Description	Returns a single stock record by its UUID
//	@Tags			Stocks
//	@Produce		json
//	@Param			id	path		string	true	"Stock UUID"
//	@Success		200	{object}	APIResponse{data=Stock}	"Stock retrieved successfully"
//	@Failure		400	{object}	APIResponse					"Invalid stock ID"
//	@Failure		404	{object}	APIResponse					"Stock not found"
//	@Failure		500	{object}	APIResponse					"Internal server error"
//	@Router			/stocks/{id} [get]
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

// GetByTicker godoc
//
//	@Summary	Get stocks by ticker
//	@Description	Returns all stock records matching the given ticker symbol
//	@Tags			Stocks
//	@Produce		json
//	@Param			ticker	path		string	true	"Ticker symbol (e.g. AAPL)"
//	@Success		200		{object}	APIResponse{data=[]Stock}	"Stocks retrieved successfully"
//	@Failure		400		{object}	APIResponse					"Ticker is required"
//	@Failure		500		{object}	APIResponse					"Internal server error"
//	@Router			/stocks/ticker/{ticker} [get]
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

// GetActions godoc
//
//	@Summary	Get distinct actions
//	@Description	Returns all distinct stock action types available in the system
//	@Tags			Stocks
//	@Produce		json
//	@Success		200	{object}	APIResponse{data=[]string}	"Actions retrieved successfully"
//	@Failure		500	{object}	APIResponse				"Internal server error"
//	@Router			/stocks/actions [get]
func (h *StockHandler) GetActions(c *gin.Context) {
	actions, err := h.stockUsecase.GetDistinctActions(c.Request.Context())
	if err != nil {
		response.InternalServerError(c.Writer, err.Error())
		return
	}

	response.Success(c.Writer, http.StatusOK, en.ActionsRetrieved, actions)
}

// SyncStocks godoc
//
//	@Summary	Sync stocks from external API
//	@Description	Fetches the latest stock data from the external Karenai API and upserts into the database
//	@Tags			Sync
//	@Produce		json
//	@Success		200	{object}	APIResponse{data=object{count=int}}	"Sync completed successfully"
//	@Failure		500	{object}	APIResponse						"Internal server error"
//	@Router			/sync [post]
func (h *StockHandler) SyncStocks(c *gin.Context) {
	count, err := h.stockUsecase.SyncFromExternalAPI(c.Request.Context())
	if err != nil {
		response.InternalServerError(c.Writer, err.Error())
		return
	}

	response.Success(c.Writer, http.StatusOK, en.SyncCompleted, gin.H{"count": count})
}

// GetRecommendations godoc
//
//	@Summary	Get stock recommendations
//	@Description	Returns ranked stock recommendations based on analyst consensus, momentum, rating upgrades, and target price changes
//	@Tags			Recommendations
//	@Produce		json
//	@Param			limit	query		int		false	"Maximum number of recommendations"	default(50)
//	@Param			search	query		string	false	"Search filter for ticker or company"
//	@Success		200		{object}	APIResponse{data=[]StockRecommendation}	"Recommendations retrieved successfully"
//	@Failure		500		{object}	APIResponse									"Internal server error"
//	@Router			/recommendations [get]
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

// GetTopRecommendation godoc
//
//	@Summary	Get top recommendation
//	@Description	Returns the single best stock recommendation with the highest composite score
//	@Tags			Recommendations
//	@Produce		json
//	@Success		200	{object}	APIResponse{data=StockRecommendation}	"Top recommendation retrieved successfully"
//	@Failure		404	{object}	APIResponse								"No recommendations available"
//	@Failure		500	{object}	APIResponse								"Internal server error"
//	@Router			/recommendations/top [get]
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
