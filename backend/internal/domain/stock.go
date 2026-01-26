package domain

import (
	"time"

	"github.com/google/uuid"
)

type Stock struct {
	ID         uuid.UUID `json:"id"`
	Ticker     string    `json:"ticker"`
	Company    string    `json:"company"`
	Brokerage  string    `json:"brokerage"`
	Action     string    `json:"action"`
	RatingFrom string    `json:"ratingFrom"`
	RatingTo   string    `json:"ratingTo"`
	TargetFrom float64   `json:"targetFrom"`
	TargetTo   float64   `json:"targetTo"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type StockFilter struct {
	Search    string
	Ticker    string
	Action    string
	SortBy    string
	SortOrder string
	Page      int
	Limit     int
}

func NewStockFilter() StockFilter {
	return StockFilter{
		SortBy:    "created_at",
		SortOrder: "desc",
		Page:      1,
		Limit:     20,
	}
}

type StockRecommendation struct {
	Stock           Stock    `json:"stock"`
	Score           float64  `json:"score"`
	Reasons         []string `json:"reasons"`
	UpsidePotential float64  `json:"upsidePotential"`
}

type PaginatedStocks struct {
	Data       []Stock `json:"data"`
	Page       int     `json:"page"`
	Limit      int     `json:"limit"`
	TotalCount int64   `json:"totalCount"`
	TotalPages int     `json:"totalPages"`
	HasNext    bool    `json:"hasNext"`
	HasPrev    bool    `json:"hasPrev"`
}
