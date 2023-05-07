package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	DSN    string
	db     *sql.DB
	ctx    context.Context
	cancel func()
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

	db.ctx, db.cancel = context.WithCancel(context.Background())
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
	db.cancel()
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}
