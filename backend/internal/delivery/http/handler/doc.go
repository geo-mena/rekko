package handler

import (
	"github.com/geomena/stock-recommendation-system/backend/internal/delivery/http/response"
	"github.com/geomena/stock-recommendation-system/backend/internal/domain"
)

type APIResponse = response.Response
type PaginationMeta = response.Meta
type Pagination = response.Pagination
type Stock = domain.Stock
type StockRecommendation = domain.StockRecommendation
type DashboardStats = domain.DashboardStats
type ActionDistribution = domain.ActionDistribution
type BrokerageDistribution = domain.BrokerageDistribution
type DailyActivity = domain.DailyActivity
