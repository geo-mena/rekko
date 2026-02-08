package cockroachdb

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/cockroachdb"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type DB struct {
	conn     *sql.DB
	dbDriver string
}

func NewDB(databaseURL, dbDriver string) (*DB, error) {
	conn, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(5 * time.Minute)

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{conn: conn, dbDriver: dbDriver}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) Conn() *sql.DB {
	return db.conn
}

func (db *DB) RunMigrations(migrationsPath string) error {
	if db.dbDriver != "postgres" {
		if err := db.ensureDatabase(); err != nil {
			return fmt.Errorf("failed to ensure database exists: %w", err)
		}
	}

	driver, driverName, err := db.createMigrationDriver()
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		driverName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func (db *DB) createMigrationDriver() (database.Driver, string, error) {
	if db.dbDriver == "postgres" {
		driver, err := postgres.WithInstance(db.conn, &postgres.Config{})
		return driver, "postgres", err
	}

	driver, err := cockroachdb.WithInstance(db.conn, &cockroachdb.Config{})
	return driver, "cockroachdb", err
}

func (db *DB) ensureDatabase() error {
	_, err := db.conn.Exec(`CREATE DATABASE IF NOT EXISTS stockdb`)
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	_, err = db.conn.Exec(`SET DATABASE = stockdb`)
	if err != nil {
		return fmt.Errorf("failed to set database: %w", err)
	}

	return nil
}
