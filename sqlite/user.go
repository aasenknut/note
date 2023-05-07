package sqlite

import (
	"context"
	"fmt"

	"github.com/aasenknut/note"
)

type UserService struct {
	DB *DB
}

func NewUserService(db *DB) *UserService {
	return &UserService{DB: db}
}

func (d *DB) GetUserByID(ctx context.Context, id int) (*note.User, error) {
	return &note.User{}, fmt.Errorf("_ NOT IMPLEMENTED _")
}

func (d *DB) GetUserByName(ctx context.Context, id int) (*note.User, error) {
	return &note.User{}, fmt.Errorf("_ NOT IMPLEMENTED _")
}
