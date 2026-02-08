package usecase

import (
	"context"
	"math"
	"sort"

	"github.com/geomena/stock-recommendation-system/backend/internal/domain"
	"github.com/geomena/stock-recommendation-system/backend/internal/external/finnhub"
	"github.com/geomena/stock-recommendation-system/backend/internal/repository"
)

const (
	WeightUpgrade        = 0.15
	WeightTargetIncrease = 0.10
	WeightActionType     = 0.15
	WeightConsensus      = 0.20
	WeightMomentum       = 0.10
	WeightRealUpside     = 0.15
	WeightMarketCap      = 0.10
	WeightPriceTrend     = 0.05
)

const (
	FallbackWeightUpgrade        = 0.20
	FallbackWeightTargetIncrease = 0.20
	FallbackWeightActionType     = 0.20
	FallbackWeightConsensus      = 0.25
	FallbackWeightMomentum       = 0.15
)

type RecommendationUsecase struct {
	stockRepo     repository.StockRepository
	finnhubClient *finnhub.Client
}

func NewRecommendationUsecase(stockRepo repository.StockRepository, finnhubClient *finnhub.Client) *RecommendationUsecase {
	return &RecommendationUsecase{
		stockRepo:     stockRepo,
		finnhubClient: finnhubClient,
	}
}

func (u *RecommendationUsecase) GetTopRecommendations(ctx context.Context, limit int, search string) ([]domain.StockRecommendation, error) {
	if limit < 1 || limit > 100 {
		limit = 50
	}

	filter := domain.StockFilter{
		Page:      1,
		Limit:     500,
		SortBy:    "created_at",
		SortOrder: "desc",
		Search:    search,
	}

	stocks, _, err := u.stockRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	tickerMap := groupByTicker(stocks)
	marketDataMap := u.fetchMarketDataForTickers(ctx, tickerMap)
	recommendations := u.scoreAllTickers(tickerMap, marketDataMap)

	if recommendations == nil {
		recommendations = []domain.StockRecommendation{}
	}

	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	return recommendations, nil
}

func (u *RecommendationUsecase) GetBestStock(ctx context.Context) (*domain.StockRecommendation, error) {
	recommendations, err := u.GetTopRecommendations(ctx, 1, "")
	if err != nil {
		return nil, err
	}
	if len(recommendations) == 0 {
		return nil, nil
	}
	return &recommendations[0], nil
}

func (u *RecommendationUsecase) fetchMarketDataForTickers(ctx context.Context, tickerMap map[string][]domain.Stock) map[string]*domain.MarketData {
	if u.finnhubClient == nil {
		return nil
	}

	tickers := make([]string, 0, len(tickerMap))
	for ticker := range tickerMap {
		tickers = append(tickers, ticker)
	}

	return u.finnhubClient.FetchBatch(ctx, tickers)
}

func (u *RecommendationUsecase) scoreAllTickers(tickerMap map[string][]domain.Stock, marketDataMap map[string]*domain.MarketData) []domain.StockRecommendation {
	var recommendations []domain.StockRecommendation

	for ticker, tickerStocks := range tickerMap {
		if len(tickerStocks) == 0 {
			continue
		}

		var md *domain.MarketData
		if marketDataMap != nil {
			md = marketDataMap[ticker]
		}

		rec := u.scoreTickerGroup(tickerStocks, md)
		if rec.Score > 0 {
			recommendations = append(recommendations, rec)
		}
	}

	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	return recommendations
}

func (u *RecommendationUsecase) scoreTickerGroup(tickerStocks []domain.Stock, md *domain.MarketData) domain.StockRecommendation {
	bestStock := tickerStocks[0]
	bestIndividualScore := 0.0
	var bestReasons []string

	for _, stock := range tickerStocks {
		score, reasons := u.calculateIndividualScore(stock)
		if score > bestIndividualScore {
			bestIndividualScore = score
			bestStock = stock
			bestReasons = reasons
		}
	}

	upgradeScore, _ := u.calculateRatingUpgrade(bestStock)
	targetIncreaseScore, _ := u.calculateTargetIncrease(bestStock)
	actionScore, _ := u.calculateActionScore(bestStock)
	consensusScore, consensusReason := u.calculateConsensusScore(tickerStocks)
	momentumScore, momentumReason := u.calculateMomentumScore(tickerStocks)

	var reasons []string
	reasons = append(reasons, bestReasons...)
	if consensusReason != "" {
		reasons = append(reasons, consensusReason)
	}
	if momentumReason != "" {
		reasons = append(reasons, momentumReason)
	}

	var totalScore float64
	upsidePotential := calculateUpsidePotentialFromAnalysts(tickerStocks)

	if md != nil {
		realUpsideScore, realUpsideReason := u.calculateRealUpside(tickerStocks, md)
		marketCapScore, marketCapReason := calculateMarketCapScore(md)
		priceTrendScore, priceTrendReason := u.calculatePriceTrendScore(md, tickerStocks)

		if realUpsideReason != "" {
			reasons = append(reasons, realUpsideReason)
		}
		if marketCapReason != "" {
			reasons = append(reasons, marketCapReason)
		}
		if priceTrendReason != "" {
			reasons = append(reasons, priceTrendReason)
		}

		totalScore = (upgradeScore*WeightUpgrade +
			targetIncreaseScore*WeightTargetIncrease +
			actionScore*WeightActionType +
			consensusScore*WeightConsensus +
			momentumScore*WeightMomentum +
			realUpsideScore*WeightRealUpside +
			marketCapScore*WeightMarketCap +
			priceTrendScore*WeightPriceTrend) / 10.0

		if md.CurrentPrice > 0 {
			avgTarget := averageTargetTo(tickerStocks)
			if avgTarget > 0 {
				upsidePotential = ((avgTarget - md.CurrentPrice) / md.CurrentPrice) * 100
			}
		}
	} else {
		totalScore = (upgradeScore*FallbackWeightUpgrade +
			targetIncreaseScore*FallbackWeightTargetIncrease +
			actionScore*FallbackWeightActionType +
			consensusScore*FallbackWeightConsensus +
			momentumScore*FallbackWeightMomentum) / 10.0
	}

	totalScore = math.Round(totalScore*10) / 10

	return domain.StockRecommendation{
		Stock:           bestStock,
		Score:           totalScore,
		Reasons:         reasons,
		UpsidePotential: math.Round(upsidePotential*10) / 10,
		AnalystCount:    countDistinctBrokerages(tickerStocks),
		MarketData:      md,
	}
}

func (u *RecommendationUsecase) calculateIndividualScore(stock domain.Stock) (float64, []string) {
	score := 0.0
	var reasons []string

	upgradeScore, upgradeReason := u.calculateRatingUpgrade(stock)
	score += upgradeScore * FallbackWeightUpgrade
	if upgradeReason != "" {
		reasons = append(reasons, upgradeReason)
	}

	targetScore, targetReason := u.calculateTargetIncrease(stock)
	score += targetScore * FallbackWeightTargetIncrease
	if targetReason != "" {
		reasons = append(reasons, targetReason)
	}

	actionScore, actionReason := u.calculateActionScore(stock)
	score += actionScore * FallbackWeightActionType
	if actionReason != "" {
		reasons = append(reasons, actionReason)
	}

	return score, reasons
}
