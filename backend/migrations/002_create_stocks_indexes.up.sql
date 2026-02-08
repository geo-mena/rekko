-- 002_create_stocks_indexes.up.sql
-- Creates indexes for optimizing stock queries

CREATE INDEX IF NOT EXISTS idx_stocks_ticker ON stocks(ticker);
CREATE INDEX IF NOT EXISTS idx_stocks_company ON stocks(company);
CREATE INDEX IF NOT EXISTS idx_stocks_action ON stocks(action);
CREATE INDEX IF NOT EXISTS idx_stocks_created_at ON stocks(created_at DESC);
