package http

import (
	"github.com/geomena/stock-recommendation-system/backend/internal/delivery/http/handler"
	"github.com/geomena/stock-recommendation-system/backend/internal/delivery/http/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(stockHandler *handler.StockHandler, healthHandler *handler.HealthHandler) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.Logging())
	router.Use(middleware.CORS())

	router.GET("/api/v1/health", healthHandler.Health)

	api := router.Group("/api/v1")
	{
		api.GET("/stocks", stockHandler.ListStocks)
		api.GET("/stocks/:id", stockHandler.GetStock)
		api.GET("/stocks/ticker/:ticker", stockHandler.GetByTicker)
		api.GET("/stocks/actions", stockHandler.GetActions)

		api.POST("/sync", stockHandler.SyncStocks)

		api.GET("/recommendations", stockHandler.GetRecommendations)
		api.GET("/recommendations/top", stockHandler.GetTopRecommendation)
	}

	return router
}
