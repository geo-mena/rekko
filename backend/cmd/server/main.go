package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/geomena/stock-recommendation-system/backend/internal/config"
	httpDelivery "github.com/geomena/stock-recommendation-system/backend/internal/delivery/http"
	"github.com/geomena/stock-recommendation-system/backend/internal/delivery/http/handler"
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

	db, err := cockroachdb.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := db.RunMigrations(ctx); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed")

	stockRepo := cockroachdb.NewStockRepository(db)
	karenaiClient := karenai.NewClient(cfg.KarenaiAPIURL, cfg.KarenaiAPIToken)

	stockUsecase := usecase.NewStockUsecase(stockRepo, karenaiClient)
	recommendationUsecase := usecase.NewRecommendationUsecase(stockRepo)
	dashboardUsecase := usecase.NewDashboardUsecase(stockRepo)

	stockHandler := handler.NewStockHandler(stockUsecase, recommendationUsecase)
	healthHandler := handler.NewHealthHandler()
	dashboardHandler := handler.NewDashboardHandler(dashboardUsecase)

	router := httpDelivery.NewRouter(stockHandler, healthHandler, dashboardHandler)

	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := router.Run(":" + cfg.ServerPort); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
