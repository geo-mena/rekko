package handler

import (
	"github.com/geomena/stock-recommendation-system/backend/internal/delivery/http/response"
	"github.com/geomena/stock-recommendation-system/backend/internal/domain"
)

// Type aliases so swag resolves short names instead of full module paths.

type APIResponse = response.Response
type PaginationMeta = response.Meta
type Pagination = response.Pagination
type Stock = domain.Stock
type StockRecommendation = domain.StockRecommendation
type DashboardStats = domain.DashboardStats
type ActionDistribution = domain.ActionDistribution
type BrokerageDistribution = domain.BrokerageDistribution
type DailyActivity = domain.DailyActivity
