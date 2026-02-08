package main

import (
	"github.com/geomena/stock-recommendation-system/backend/internal/config"
	httpDelivery "github.com/geomena/stock-recommendation-system/backend/internal/delivery/http"
	"github.com/geomena/stock-recommendation-system/backend/internal/delivery/http/handler"
	"github.com/geomena/stock-recommendation-system/backend/internal/external/finnhub"
	"github.com/geomena/stock-recommendation-system/backend/internal/external/karenai"
	"github.com/geomena/stock-recommendation-system/backend/internal/repository/cockroachdb"
	"github.com/geomena/stock-recommendation-system/backend/internal/usecase"

	_ "github.com/geomena/stock-recommendation-system/backend/docs"
)

//	@title						Rekko API
//	@version					1.0
//	@description				API for stock recommendations, analyst ratings, and portfolio insights.
//	@host						localhost:8080
//	@BasePath					/api/v1
//	@produce					json
//	@consumes					json
func main() {
	cfg := config.Load()

	db := initDatabase(cfg.DatabaseURL, cfg.MigrationsPath, cfg.DBDriver)
	defer db.Close()

	stockRepo := cockroachdb.NewStockRepository(db)
	karenaiClient := karenai.NewClient(cfg.KarenaiAPIURL, cfg.KarenaiAPIToken)

	var finnhubClient *finnhub.Client
	if cfg.FinnhubAPIKey != "" {
		finnhubClient = finnhub.NewClient(cfg.FinnhubAPIKey)
	}

	stockUsecase := usecase.NewStockUsecase(stockRepo, karenaiClient)
	recommendationUsecase := usecase.NewRecommendationUsecase(stockRepo, finnhubClient)
	dashboardUsecase := usecase.NewDashboardUsecase(stockRepo)

	stockHandler := handler.NewStockHandler(stockUsecase, recommendationUsecase)
	healthHandler := handler.NewHealthHandler()
	dashboardHandler := handler.NewDashboardHandler(dashboardUsecase)

	router := httpDelivery.NewRouter(stockHandler, healthHandler, dashboardHandler, cfg.StaticDir)

	startServer(router, cfg.ServerPort)
	waitForShutdown()
}
