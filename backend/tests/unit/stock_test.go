package unit_test

import (
	"context"
	"errors"
	"testing"

	"github.com/geomena/stock-recommendation-system/backend/internal/domain"
)

func TestListStocks_DefaultPagination(t *testing.T) {
	tests := []struct {
		name          string
		page          int
		limit         int
		expectedPage  int
		expectedLimit int
	}{
		{
			name:          "zero page defaults to 1",
			page:          0,
			limit:         10,
			expectedPage:  1,
			expectedLimit: 10,
		},
		{
			name:          "negative page defaults to 1",
			page:          -5,
			limit:         10,
			expectedPage:  1,
			expectedLimit: 10,
		},
		{
			name:          "zero limit defaults to 20",
			page:          1,
			limit:         0,
			expectedPage:  1,
			expectedLimit: 20,
		},
		{
			name:          "negative limit defaults to 20",
			page:          1,
			limit:         -1,
			expectedPage:  1,
			expectedLimit: 20,
		},
		{
			name:          "limit over 100 defaults to 20",
			page:          1,
			limit:         200,
			expectedPage:  1,
			expectedLimit: 20,
		},
		{
			name:          "valid values kept as-is",
			page:          3,
			limit:         50,
			expectedPage:  3,
			expectedLimit: 50,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock := newMockRepo()
			mock.FindAllFn = func(ctx context.Context, filter domain.StockFilter) ([]domain.Stock, int64, error) {
				if filter.Page != tc.expectedPage {
					t.Errorf("expected page %d, got %d", tc.expectedPage, filter.Page)
				}
				if filter.Limit != tc.expectedLimit {
					t.Errorf("expected limit %d, got %d", tc.expectedLimit, filter.Limit)
				}
				return []domain.Stock{}, 0, nil
			}

			uc := newStockUsecase(mock)
			_, err := uc.ListStocks(context.Background(), domain.StockFilter{Page: tc.page, Limit: tc.limit})
			assertNoError(t, err)
		})
	}
}

func TestListStocks_PaginationMetadata(t *testing.T) {
	tests := []struct {
		name           string
		page           int
		limit          int
		totalCount     int64
		expectedPages  int
		expectedNext   bool
		expectedPrev   bool
	}{
		{
			name:          "first page of many",
			page:          1,
			limit:         10,
			totalCount:    25,
			expectedPages: 3,
			expectedNext:  true,
			expectedPrev:  false,
		},
		{
			name:          "middle page",
			page:          2,
			limit:         10,
			totalCount:    25,
			expectedPages: 3,
			expectedNext:  true,
			expectedPrev:  true,
		},
		{
			name:          "last page",
			page:          3,
			limit:         10,
			totalCount:    25,
			expectedPages: 3,
			expectedNext:  false,
			expectedPrev:  true,
		},
		{
			name:          "single page",
			page:          1,
			limit:         50,
			totalCount:    10,
			expectedPages: 1,
			expectedNext:  false,
			expectedPrev:  false,
		},
		{
			name:          "empty results",
			page:          1,
			limit:         20,
			totalCount:    0,
			expectedPages: 0,
			expectedNext:  false,
			expectedPrev:  false,
		},
		{
			name:          "exact boundary",
			page:          1,
			limit:         10,
			totalCount:    10,
			expectedPages: 1,
			expectedNext:  false,
			expectedPrev:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock := newMockRepo()
			mock.FindAllFn = func(ctx context.Context, filter domain.StockFilter) ([]domain.Stock, int64, error) {
				return []domain.Stock{}, tc.totalCount, nil
			}

			uc := newStockUsecase(mock)
			result, err := uc.ListStocks(context.Background(), domain.StockFilter{Page: tc.page, Limit: tc.limit})
			assertNoError(t, err)

			if result.TotalPages != tc.expectedPages {
				t.Errorf("expected totalPages %d, got %d", tc.expectedPages, result.TotalPages)
			}
			if result.HasNext != tc.expectedNext {
				t.Errorf("expected hasNext %v, got %v", tc.expectedNext, result.HasNext)
			}
			if result.HasPrev != tc.expectedPrev {
				t.Errorf("expected hasPrev %v, got %v", tc.expectedPrev, result.HasPrev)
			}
			if result.Page != tc.page {
				t.Errorf("expected page %d, got %d", tc.page, result.Page)
			}
		})
	}
}

func TestListStocks_NilStocksNormalized(t *testing.T) {
	mock := newMockRepo()
	mock.FindAllFn = func(ctx context.Context, filter domain.StockFilter) ([]domain.Stock, int64, error) {
		return nil, 0, nil
	}

	uc := newStockUsecase(mock)
	result, err := uc.ListStocks(context.Background(), domain.StockFilter{Page: 1, Limit: 10})
	assertNoError(t, err)

	if result.Data == nil {
		t.Fatal("expected Data to be non-nil empty slice, got nil")
	}
	if len(result.Data) != 0 {
		t.Errorf("expected 0 stocks, got %d", len(result.Data))
	}
}

func TestListStocks_RepoError(t *testing.T) {
	mock := newMockRepo()
	mock.FindAllFn = func(ctx context.Context, filter domain.StockFilter) ([]domain.Stock, int64, error) {
		return nil, 0, errors.New("connection refused")
	}

	uc := newStockUsecase(mock)
	result, err := uc.ListStocks(context.Background(), domain.StockFilter{Page: 1, Limit: 10})
	assertError(t, err)

	if result != nil {
		t.Errorf("expected nil result on error, got %+v", result)
	}
}
