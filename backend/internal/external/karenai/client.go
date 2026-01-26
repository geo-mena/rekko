package karenai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/geomena/stock-recommendation-system/backend/internal/domain"
)

type Client struct {
	baseURL    string
	authToken  string
	httpClient *http.Client
}

type APIResponse struct {
	Items    []StockItem `json:"items"`
	NextPage string      `json:"next_page"`
}

type StockItem struct {
	Ticker     string `json:"ticker"`
	Company    string `json:"company"`
	Brokerage  string `json:"brokerage"`
	Action     string `json:"action"`
	RatingFrom string `json:"rating_from"`
	RatingTo   string `json:"rating_to"`
	TargetFrom string `json:"target_from"`
	TargetTo   string `json:"target_to"`
}

func parsePriceString(price string) float64 {
	if price == "" {
		return 0.0
	}

	cleaned := strings.TrimSpace(price)
	cleaned = strings.TrimPrefix(cleaned, "$")
	cleaned = strings.TrimSpace(cleaned)

	if cleaned == "" {
		return 0.0
	}

	value, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0.0
	}

	return value
}

func NewClient(baseURL, authToken string) *Client {
	return &Client{
		baseURL:   baseURL,
		authToken: authToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) FetchStocks(ctx context.Context, nextPage string) (*APIResponse, error) {
	url := c.baseURL + "/swechallenge/list"
	if nextPage != "" {
		url += "?next_page=" + nextPage
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &apiResp, nil
}

func (c *Client) FetchAllStocks(ctx context.Context) ([]domain.Stock, error) {
	var allStocks []domain.Stock
	nextPage := ""

	for {
		resp, err := c.FetchStocks(ctx, nextPage)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.Items {
			stock := domain.Stock{
				Ticker:     item.Ticker,
				Company:    item.Company,
				Brokerage:  item.Brokerage,
				Action:     item.Action,
				RatingFrom: item.RatingFrom,
				RatingTo:   item.RatingTo,
				TargetFrom: parsePriceString(item.TargetFrom),
				TargetTo:   parsePriceString(item.TargetTo),
			}
			allStocks = append(allStocks, stock)
		}

		if resp.NextPage == "" {
			break
		}
		nextPage = resp.NextPage
	}

	return allStocks, nil
}
