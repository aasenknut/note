package redis

import (
	"context"
	"fmt"
	"time"
)

type AuthService struct {
	db *DB
}

func NewAuthService(db *DB) *AuthService {
	return &AuthService{db: db}
}

func (as *AuthService) GetUserID(ctx context.Context, userID string) (string, error) {
	userID, err := as.db.get(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("get user_id by token: %v", err)
	}
	return userID, nil
}

func (as *AuthService) DelAuth(ctx context.Context, userID string) error {
	return fmt.Errorf("--- NOT IMPLEMENTED ---")
}

func (as *AuthService) SetAuth(ctx context.Context, token string, userID int) (time.Duration, error) {
	ttl := time.Hour * 24
	err := as.db.set(ctx, token, userID, ttl)
	if err != nil {
		return time.Duration(0), fmt.Errorf("transaction auth: %v", err)
	}
	return ttl, nil
}
