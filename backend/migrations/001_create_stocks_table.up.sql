-- 001_create_stocks_table.up.sql
-- Creates the stocks table for storing analyst recommendations

CREATE TABLE IF NOT EXISTS stocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticker VARCHAR(10) NOT NULL,
    company VARCHAR(255) NOT NULL,
    brokerage VARCHAR(255) NOT NULL,
    action VARCHAR(50) NOT NULL,
    rating_from VARCHAR(50),
    rating_to VARCHAR(50),
    target_from DECIMAL(10, 2),
    target_to DECIMAL(10, 2),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (ticker, brokerage, action, rating_from, rating_to, target_from, target_to)
);
