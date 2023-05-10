package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

const migrationPath string = "./sqlite/migration"

type DB struct {
	DSN string
	db  *sql.DB
}

func NewDB() *DB {
	return &DB{}
}

func (db *DB) Open() error {
	var err error
	db.db, err = sql.Open("sqlite3", db.DSN)
	if err != nil {
		return err
	}
	db.db.SetMaxOpenConns(1)

	// Enable foreign key check
	if _, err := db.db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return fmt.Errorf("foreign keys pragma: %w", err)
	}
	return nil
}

type Tx struct {
	*sql.Tx
	db  *DB
	now time.Time
}

func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("sql transaction")
	}
	return &Tx{tx, db, time.Now()}, nil
}

func (db *DB) Close() error {
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}

func (db *DB) Migrate() error {
	driver, err := sqlite.WithInstance(db.db, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("migrate instance: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", migrationPath), "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("migrate setup: %v", err)
	}
	if err := m.Up(); err != nil {
		return fmt.Errorf("up migration: %v", err)
	}
	return nil
}
