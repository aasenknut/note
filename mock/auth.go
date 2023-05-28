package mock

import (
	"context"
	"time"
)

type AuthService struct {
	getUserIDFn func(ctx context.Context, token string) (string, error)
	delAuthFn   func(ctx context.Context, userID string) error
	setAuthFn   func(ctx context.Context, token string, userID int) (time.Duration, error)
}

func (as *AuthService) GetUserID(ctx context.Context, token string) (string, error) {
	return as.getUserIDFn(ctx, token)
}
func (as *AuthService) DelAuth(ctx context.Context, userID string) error {
	return as.delAuthFn(ctx, userID)
}

func (as *AuthService) SetAuth(ctx context.Context, token string, userID int) (time.Duration, error) {
	return as.setAuthFn(ctx, token, userID)
}
