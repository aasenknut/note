package note

import "context"

type User struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type UserService interface {
	GetByName(ctx context.Context, id int) (*User, error)
	GetByID(ctx context.Context, id int) (*User, error)
}
