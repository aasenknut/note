package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type DB struct {
	db *redis.Client

	Addr     string
	Password string
	DBIndex  int
}

func NewDB() *DB {
	return &DB{}
}

func (db *DB) Open() {
	db.db = redis.NewClient(&redis.Options{
		Addr:     db.Addr,
		Password: db.Password,
		DB:       db.DBIndex,
	})
}

func (db *DB) Close() error {
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}

func (db *DB) set(ctx context.Context, key string, value int, ttl time.Duration) error {
	if err := db.db.Set(ctx, key, value, ttl).Err(); err != nil {
		return fmt.Errorf("inseting to storage (SET): %v", err)
	}
	return nil
}

func (db *DB) get(ctx context.Context, key string) (string, error) {
	cmd, err := db.db.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("get from storage (GET): %v", err)
	}
	return cmd, nil
}
