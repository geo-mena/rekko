package finnhub

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/geomena/stock-recommendation-system/backend/internal/domain"
)

const (
	baseURL       = "https://finnhub.io/api/v1"
	cacheTTL      = 15 * time.Minute
	maxConcurrent = 10
)

type quoteResponse struct {
	Current       float64 `json:"c"`
	Change        float64 `json:"d"`
	ChangePercent float64 `json:"dp"`
	High          float64 `json:"h"`
	Low           float64 `json:"l"`
	Open          float64 `json:"o"`
	PreviousClose float64 `json:"pc"`
}

type profileResponse struct {
	MarketCap float64 `json:"marketCapitalization"`
	Name      string  `json:"name"`
	Industry  string  `json:"finnhubIndustry"`
	Exchange  string  `json:"exchange"`
}

type cacheEntry struct {
	data      *domain.MarketData
	expiresAt time.Time
}

type Client struct {
	apiKey     string
	httpClient *http.Client
	cache      map[string]cacheEntry
	mu         sync.RWMutex
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		cache: make(map[string]cacheEntry),
	}
}

func (c *Client) FetchMarketData(ctx context.Context, symbol string) (*domain.MarketData, error) {
	if cached := c.getFromCache(symbol); cached != nil {
		return cached, nil
	}

	quote, err := c.fetchQuote(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("quote for %s: %w", symbol, err)
	}

	if quote.Current == 0 {
		return nil, nil
	}

	profile, err := c.fetchProfile(ctx, symbol)
	if err != nil {
		profile = &profileResponse{}
	}

	data := &domain.MarketData{
		CurrentPrice:  quote.Current,
		DayChange:     quote.Change,
		DayChangePct:  quote.ChangePercent,
		DayHigh:       quote.High,
		DayLow:        quote.Low,
		PreviousClose: quote.PreviousClose,
		MarketCap:     profile.MarketCap,
		Industry:      profile.Industry,
	}

	c.setCache(symbol, data)
	return data, nil
}

func (c *Client) FetchBatch(ctx context.Context, tickers []string) map[string]*domain.MarketData {
	results := make(map[string]*domain.MarketData)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrent)

	for _, ticker := range tickers {
		wg.Add(1)
		sem <- struct{}{}
		go func(t string) {
			defer wg.Done()
			defer func() { <-sem }()

			data, err := c.FetchMarketData(ctx, t)
			if err != nil {
				return
			}
			if data == nil {
				return
			}

			mu.Lock()
			results[t] = data
			mu.Unlock()
		}(ticker)
	}

	wg.Wait()
	return results
}

func (c *Client) fetchQuote(ctx context.Context, symbol string) (*quoteResponse, error) {
	url := fmt.Sprintf("%s/quote?symbol=%s", baseURL, symbol)
	return doRequest[quoteResponse](ctx, c.httpClient, url, c.apiKey)
}

func (c *Client) fetchProfile(ctx context.Context, symbol string) (*profileResponse, error) {
	url := fmt.Sprintf("%s/stock/profile2?symbol=%s", baseURL, symbol)
	return doRequest[profileResponse](ctx, c.httpClient, url, c.apiKey)
}

func doRequest[T any](ctx context.Context, client *http.Client, url, apiKey string) (*T, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("X-Finnhub-Token", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return &result, nil
}

func (c *Client) getFromCache(symbol string) *domain.MarketData {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.cache[symbol]
	if !ok || time.Now().After(entry.expiresAt) {
		return nil
	}
	return entry.data
}

func (c *Client) setCache(symbol string, data *domain.MarketData) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[symbol] = cacheEntry{
		data:      data,
		expiresAt: time.Now().Add(cacheTTL),
	}
}
