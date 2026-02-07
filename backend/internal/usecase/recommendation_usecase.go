package usecase

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/geomena/stock-recommendation-system/backend/internal/domain"
	"github.com/geomena/stock-recommendation-system/backend/internal/i18n/en"
	"github.com/geomena/stock-recommendation-system/backend/internal/repository"
)

const (
	WeightConsensus      = 0.25
	WeightMomentum       = 0.15
	WeightUpgrade        = 0.20
	WeightTargetIncrease = 0.20
	WeightActionType     = 0.20

	momentumDecayDays   = 30.0
	momentumSaturationK = 2.0
)

var ratingValues = map[string]int{
	"strong sell":    1,
	"sell":           2,
	"underweight":    2,
	"hold":           3,
	"neutral":        3,
	"equal-weight":   3,
	"market perform": 3,
	"sector perform": 3,
	"buy":            4,
	"overweight":     4,
	"outperform":     4,
	"strong buy":     5,
	"top pick":       5,
}

var actionScores = map[string]float64{
	"upgraded":       100,
	"initiated":      80,
	"reiterated":     60,
	"target raised":  70,
	"maintained":     50,
	"downgraded":     20,
	"target lowered": 30,
}

var bullishActions = map[string]bool{
	"upgraded":      true,
	"initiated":     true,
	"target raised": true,
	"reiterated":    true,
	"maintained":    true,
}

type RecommendationUsecase struct {
	stockRepo repository.StockRepository
}

func NewRecommendationUsecase(stockRepo repository.StockRepository) *RecommendationUsecase {
	return &RecommendationUsecase{stockRepo: stockRepo}
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

	recommendations := u.scoreAndRankStocks(stocks)

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

func (u *RecommendationUsecase) scoreAndRankStocks(stocks []domain.Stock) []domain.StockRecommendation {
	tickerMap := make(map[string][]domain.Stock)
	for _, stock := range stocks {
		tickerMap[stock.Ticker] = append(tickerMap[stock.Ticker], stock)
	}

	var recommendations []domain.StockRecommendation
	for _, tickerStocks := range tickerMap {
		if len(tickerStocks) == 0 {
			continue
		}

		rec := u.scoreTickerGroup(tickerStocks)
		if rec.Score > 0 {
			recommendations = append(recommendations, rec)
		}
	}

	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	return recommendations
}

func (u *RecommendationUsecase) scoreTickerGroup(tickerStocks []domain.Stock) domain.StockRecommendation {
	bestStock := tickerStocks[0]
	bestIndividualScore := 0.0
	var bestReasons []string
	var maxUpsidePotential float64

	for _, stock := range tickerStocks {
		score, reasons := u.calculateIndividualScore(stock)
		if score > bestIndividualScore {
			bestIndividualScore = score
			bestStock = stock
			bestReasons = reasons
		}

		if stock.TargetFrom > 0 && stock.TargetTo > 0 {
			potential := ((stock.TargetTo - stock.TargetFrom) / stock.TargetFrom) * 100
			if potential > maxUpsidePotential {
				maxUpsidePotential = potential
			}
		}
	}

	consensusScore, consensusReason := u.calculateConsensusScore(tickerStocks)
	momentumScore, momentumReason := u.calculateMomentumScore(tickerStocks)

	totalScore := (bestIndividualScore +
		consensusScore*WeightConsensus +
		momentumScore*WeightMomentum) / 10.0
	totalScore = math.Round(totalScore*10) / 10

	var reasons []string
	reasons = append(reasons, bestReasons...)
	if consensusReason != "" {
		reasons = append(reasons, consensusReason)
	}
	if momentumReason != "" {
		reasons = append(reasons, momentumReason)
	}

	return domain.StockRecommendation{
		Stock:           bestStock,
		Score:           totalScore,
		Reasons:         reasons,
		UpsidePotential: maxUpsidePotential,
		AnalystCount:    countDistinctBrokerages(tickerStocks),
	}
}

func (u *RecommendationUsecase) calculateIndividualScore(stock domain.Stock) (float64, []string) {
	score := 0.0
	var reasons []string

	upgradeScore, upgradeReason := u.calculateRatingUpgrade(stock)
	score += upgradeScore * WeightUpgrade
	if upgradeReason != "" {
		reasons = append(reasons, upgradeReason)
	}

	targetScore, targetReason := u.calculateTargetIncrease(stock)
	score += targetScore * WeightTargetIncrease
	if targetReason != "" {
		reasons = append(reasons, targetReason)
	}

	actionScore, actionReason := u.calculateActionScore(stock)
	score += actionScore * WeightActionType
	if actionReason != "" {
		reasons = append(reasons, actionReason)
	}

	return score, reasons
}

func (u *RecommendationUsecase) calculateConsensusScore(tickerStocks []domain.Stock) (float64, string) {
	brokerages := make(map[string]bool)
	bullishBrokerages := make(map[string]bool)

	for _, stock := range tickerStocks {
		brokerage := strings.ToLower(stock.Brokerage)
		brokerages[brokerage] = true
		if isBullishAction(stock.Action) {
			bullishBrokerages[brokerage] = true
		}
	}

	total := len(brokerages)
	bullish := len(bullishBrokerages)
	if total == 0 {
		return 0, ""
	}

	score := (float64(bullish) / float64(total)) * 100

	if total < 3 {
		score *= float64(total) / 3.0
	}

	return score, fmt.Sprintf(en.ReasonAnalystsBullish, bullish, total)
}

func (u *RecommendationUsecase) calculateMomentumScore(tickerStocks []domain.Stock) (float64, string) {
	now := time.Now()
	weightedSignals := 0.0
	recentCount := 0

	for _, stock := range tickerStocks {
		daysSince := now.Sub(stock.CreatedAt).Hours() / 24.0
		decayFactor := math.Exp(-daysSince / momentumDecayDays)

		if isBullishAction(stock.Action) {
			weightedSignals += decayFactor
		} else {
			weightedSignals -= decayFactor * 0.5
		}

		if daysSince <= 7 {
			recentCount++
		}
	}

	if weightedSignals <= 0 {
		return 0, ""
	}

	score := (weightedSignals / (weightedSignals + momentumSaturationK)) * 100

	if recentCount > 0 {
		return score, fmt.Sprintf(en.ReasonRecentSignals, recentCount)
	}
	return score, ""
}

func (u *RecommendationUsecase) calculateRatingUpgrade(stock domain.Stock) (float64, string) {
	fromValue := getRatingValue(stock.RatingFrom)
	toValue := getRatingValue(stock.RatingTo)

	if toValue > fromValue && fromValue > 0 {
		upgradePoints := float64(toValue-fromValue) / 4.0 * 100
		return upgradePoints, fmt.Sprintf(en.ReasonRatingUpgraded, stock.RatingFrom, stock.RatingTo)
	}

	if toValue >= 4 {
		return 50, fmt.Sprintf(en.ReasonStrongRating, stock.RatingTo)
	}

	return 0, ""
}

func (u *RecommendationUsecase) calculateTargetIncrease(stock domain.Stock) (float64, string) {
	if stock.TargetFrom <= 0 || stock.TargetTo <= 0 {
		return 0, ""
	}

	percentChange := ((stock.TargetTo - stock.TargetFrom) / stock.TargetFrom) * 100
	if percentChange > 0 {
		score := math.Min(percentChange, 100)
		return score, fmt.Sprintf(en.ReasonTargetIncreased, percentChange, stock.TargetTo)
	}

	return 0, ""
}

func (u *RecommendationUsecase) calculateActionScore(stock domain.Stock) (float64, string) {
	actionLower := strings.ToLower(stock.Action)

	for keyword, score := range actionScores {
		if strings.Contains(actionLower, keyword) {
			return score, fmt.Sprintf(en.ReasonActionBy, stock.Action, stock.Brokerage)
		}
	}

	return 30, ""
}

func isBullishAction(action string) bool {
	actionLower := strings.ToLower(action)
	for keyword := range bullishActions {
		if strings.Contains(actionLower, keyword) {
			return true
		}
	}
	return false
}

func countDistinctBrokerages(stocks []domain.Stock) int {
	seen := make(map[string]bool)
	for _, stock := range stocks {
		seen[strings.ToLower(stock.Brokerage)] = true
	}
	return len(seen)
}

func getRatingValue(rating string) int {
	ratingLower := strings.ToLower(strings.TrimSpace(rating))
	if value, ok := ratingValues[ratingLower]; ok {
		return value
	}
	return 0
}
