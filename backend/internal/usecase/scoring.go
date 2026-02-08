package usecase

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/geomena/stock-recommendation-system/backend/internal/domain"
	"github.com/geomena/stock-recommendation-system/backend/internal/i18n/en"
)

const (
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
			return score, formatActionReason(stock.Action, stock.Brokerage)
		}
	}

	return 30, ""
}

func formatActionReason(action, brokerage string) string {
	if brokerage == "" {
		return strings.TrimSuffix(strings.TrimSpace(action), "by")
	}

	if strings.HasSuffix(strings.ToLower(strings.TrimSpace(action)), "by") {
		return fmt.Sprintf("%s %s", strings.TrimSpace(action), brokerage)
	}

	return fmt.Sprintf(en.ReasonActionBy, action, brokerage)
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

func (u *RecommendationUsecase) calculateRealUpside(tickerStocks []domain.Stock, md *domain.MarketData) (float64, string) {
	if md.CurrentPrice <= 0 {
		return 0, ""
	}

	avgTarget := averageTargetTo(tickerStocks)
	if avgTarget <= 0 {
		return 0, ""
	}

	upsidePct := ((avgTarget - md.CurrentPrice) / md.CurrentPrice) * 100
	if upsidePct <= 0 {
		return 0, ""
	}

	score := math.Min(upsidePct*2, 100)

	return score, fmt.Sprintf(en.ReasonRealUpside, upsidePct, md.CurrentPrice, avgTarget)
}

func calculateMarketCapScore(md *domain.MarketData) (float64, string) {
	if md.MarketCap <= 0 {
		return 50, ""
	}

	capInBillions := md.MarketCap / 1000.0

	if capInBillions >= 10 {
		return 100, fmt.Sprintf(en.ReasonLargeCap, capInBillions)
	}
	if capInBillions >= 2 {
		return 75, fmt.Sprintf(en.ReasonMidCap, capInBillions)
	}
	if capInBillions >= 0.3 {
		return 50, ""
	}
	return 25, ""
}

func (u *RecommendationUsecase) calculatePriceTrendScore(md *domain.MarketData, tickerStocks []domain.Stock) (float64, string) {
	if md.DayChangePct == 0 {
		return 50, ""
	}

	analystsBullish := isMajorityBullish(tickerStocks)

	if analystsBullish && md.DayChangePct > 0 {
		score := math.Min(50+md.DayChangePct*10, 100)
		return score, fmt.Sprintf(en.ReasonPriceTrendUp, md.DayChangePct)
	}

	if analystsBullish && md.DayChangePct < 0 {
		return math.Max(50+md.DayChangePct*10, 0), ""
	}

	return 50, ""
}

func averageTargetTo(tickerStocks []domain.Stock) float64 {
	sum := 0.0
	count := 0
	for _, stock := range tickerStocks {
		if stock.TargetTo > 0 {
			sum += stock.TargetTo
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return sum / float64(count)
}

func isMajorityBullish(tickerStocks []domain.Stock) bool {
	bullish := 0
	for _, stock := range tickerStocks {
		if isBullishAction(stock.Action) {
			bullish++
		}
	}
	return bullish > len(tickerStocks)/2
}

func calculateUpsidePotentialFromAnalysts(tickerStocks []domain.Stock) float64 {
	var maxUpside float64
	for _, stock := range tickerStocks {
		if stock.TargetFrom > 0 && stock.TargetTo > 0 {
			potential := ((stock.TargetTo - stock.TargetFrom) / stock.TargetFrom) * 100
			if potential > maxUpside {
				maxUpside = potential
			}
		}
	}
	return maxUpside
}

func groupByTicker(stocks []domain.Stock) map[string][]domain.Stock {
	tickerMap := make(map[string][]domain.Stock)
	for _, stock := range stocks {
		tickerMap[stock.Ticker] = append(tickerMap[stock.Ticker], stock)
	}
	return tickerMap
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
