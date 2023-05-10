package note

import (
	"context"
	"time"
)

type Auth struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Token  string `json:"token"`
}

type AuthService interface {
	GetUserID(ctx context.Context, token string) (string, error)
	DelAuth(ctx context.Context, userID string) error
	SetAuth(ctx context.Context, token string, userID int) (time.Duration, error)
}
