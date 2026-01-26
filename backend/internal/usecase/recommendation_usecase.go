package usecase

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/geomena/stock-recommendation-system/backend/internal/domain"
	"github.com/geomena/stock-recommendation-system/backend/internal/repository"
)

const (
	WeightUpgrade        = 0.30
	WeightTargetIncrease = 0.25
	WeightActionType     = 0.25
	WeightTargetPrice    = 0.20
)

var ratingValues = map[string]int{
	"strong sell": 1,
	"sell":        2,
	"underweight": 2,
	"hold":        3,
	"neutral":     3,
	"equal-weight": 3,
	"market perform": 3,
	"sector perform": 3,
	"buy":         4,
	"overweight":  4,
	"outperform":  4,
	"strong buy":  5,
	"top pick":    5,
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

type RecommendationUsecase struct {
	stockRepo repository.StockRepository
}

func NewRecommendationUsecase(stockRepo repository.StockRepository) *RecommendationUsecase {
	return &RecommendationUsecase{stockRepo: stockRepo}
}

func (u *RecommendationUsecase) GetTopRecommendations(ctx context.Context, limit int) ([]domain.StockRecommendation, error) {
	if limit < 1 || limit > 50 {
		limit = 10
	}

	filter := domain.StockFilter{
		Page:      1,
		Limit:     500,
		SortBy:    "created_at",
		SortOrder: "desc",
	}

	stocks, _, err := u.stockRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	recommendations := u.scoreAndRankStocks(stocks)

	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	return recommendations, nil
}

func (u *RecommendationUsecase) GetBestStock(ctx context.Context) (*domain.StockRecommendation, error) {
	recommendations, err := u.GetTopRecommendations(ctx, 1)
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

		bestStock := tickerStocks[0]
		totalScore := 0.0
		var reasons []string
		var upsidePotential float64

		for _, stock := range tickerStocks {
			score, stockReasons := u.calculateScore(stock)
			if score > totalScore {
				totalScore = score
				bestStock = stock
				reasons = stockReasons
			}

			if stock.TargetFrom > 0 && stock.TargetTo > 0 {
				potential := ((stock.TargetTo - stock.TargetFrom) / stock.TargetFrom) * 100
				if potential > upsidePotential {
					upsidePotential = potential
				}
			}
		}

		if totalScore > 0 {
			recommendations = append(recommendations, domain.StockRecommendation{
				Stock:           bestStock,
				Score:           totalScore,
				Reasons:         reasons,
				UpsidePotential: upsidePotential,
			})
		}
	}

	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	return recommendations
}

func (u *RecommendationUsecase) calculateScore(stock domain.Stock) (float64, []string) {
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

	priceScore := u.calculateTargetPriceScore(stock)
	score += priceScore * WeightTargetPrice

	return score, reasons
}

func (u *RecommendationUsecase) calculateRatingUpgrade(stock domain.Stock) (float64, string) {
	fromValue := getRatingValue(stock.RatingFrom)
	toValue := getRatingValue(stock.RatingTo)

	if toValue > fromValue && fromValue > 0 {
		upgradePoints := float64(toValue-fromValue) / 4.0 * 100
		reason := fmt.Sprintf("Rating upgraded from %s to %s", stock.RatingFrom, stock.RatingTo)
		return upgradePoints, reason
	}

	if toValue >= 4 {
		return 50, fmt.Sprintf("Strong rating: %s", stock.RatingTo)
	}

	return 0, ""
}

func (u *RecommendationUsecase) calculateTargetIncrease(stock domain.Stock) (float64, string) {
	if stock.TargetFrom <= 0 || stock.TargetTo <= 0 {
		return 0, ""
	}

	percentChange := ((stock.TargetTo - stock.TargetFrom) / stock.TargetFrom) * 100
	if percentChange > 0 {
		score := min(percentChange, 100)
		reason := fmt.Sprintf("Target price increased %.1f%% to $%.2f", percentChange, stock.TargetTo)
		return score, reason
	}

	return 0, ""
}

func (u *RecommendationUsecase) calculateActionScore(stock domain.Stock) (float64, string) {
	actionLower := strings.ToLower(stock.Action)

	for keyword, score := range actionScores {
		if strings.Contains(actionLower, keyword) {
			reason := fmt.Sprintf("%s by %s", stock.Action, stock.Brokerage)
			return score, reason
		}
	}

	return 30, ""
}

func (u *RecommendationUsecase) calculateTargetPriceScore(stock domain.Stock) float64 {
	if stock.TargetTo <= 0 {
		return 0
	}

	if stock.TargetTo >= 500 {
		return 100
	} else if stock.TargetTo >= 200 {
		return 80
	} else if stock.TargetTo >= 100 {
		return 60
	} else if stock.TargetTo >= 50 {
		return 40
	}
	return 20
}

func getRatingValue(rating string) int {
	ratingLower := strings.ToLower(strings.TrimSpace(rating))
	if value, ok := ratingValues[ratingLower]; ok {
		return value
	}
	return 0
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
