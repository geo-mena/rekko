package domain

import "errors"

var (
	ErrStockNotFound       = errors.New("stock not found")
	ErrInvalidStockData    = errors.New("invalid stock data")
	ErrDatabaseConnection  = errors.New("database connection error")
	ErrExternalAPIFailure  = errors.New("external API failure")
	ErrInvalidFilter       = errors.New("invalid filter parameters")
	ErrSyncInProgress      = errors.New("sync already in progress")
)
