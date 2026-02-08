-- 002_create_stocks_indexes.down.sql
-- Drops all stock indexes

DROP INDEX IF EXISTS idx_stocks_ticker;
DROP INDEX IF EXISTS idx_stocks_company;
DROP INDEX IF EXISTS idx_stocks_action;
DROP INDEX IF EXISTS idx_stocks_created_at;
