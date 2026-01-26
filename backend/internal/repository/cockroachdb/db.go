package cockroachdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
}

func NewDB(databaseURL string) (*DB, error) {
	conn, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := conn.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{conn: conn}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) Conn() *sql.DB {
	return db.conn
}

func (db *DB) RunMigrations(ctx context.Context) error {
	migrations := []string{
		`CREATE DATABASE IF NOT EXISTS stockdb`,
		`SET DATABASE = stockdb`,
		`CREATE TABLE IF NOT EXISTS stocks (
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
		)`,
		`CREATE INDEX IF NOT EXISTS idx_stocks_ticker ON stocks(ticker)`,
		`CREATE INDEX IF NOT EXISTS idx_stocks_company ON stocks(company)`,
		`CREATE INDEX IF NOT EXISTS idx_stocks_action ON stocks(action)`,
		`CREATE INDEX IF NOT EXISTS idx_stocks_created_at ON stocks(created_at DESC)`,
	}

	for _, migration := range migrations {
		if _, err := db.conn.ExecContext(ctx, migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}
