package note

import "context"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserService interface {
	Create(ctx context.Context, username, password string) (int, error)
	GetByUsername(ctx context.Context, id int) (*User, error)
	GetByID(ctx context.Context, id int) (*User, error)
}
