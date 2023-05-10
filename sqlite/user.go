package sqlite

import (
	"context"
	"fmt"

	"github.com/aasenknut/note"
)

type UserService struct {
	db *DB
}

func NewUserService(db *DB) *UserService {
	return &UserService{db: db}
}

func (us *UserService) Create(ctx context.Context, username, password string) (int, error) {
	tx, err := us.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %v", err)
	}
	userQuery := `INSERT INTO user (username, password) VALUES (?, ?)`

	res, err := tx.ExecContext(ctx, userQuery, username, password)
	if err != nil {
		return 0, fmt.Errorf("create user: %s - %v", username, err)
	}
	userID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("last inserted user_id: %s - %v", username, err)
	}
	return int(userID), tx.Commit()
}

func (us *UserService) GetByID(ctx context.Context, id int) (*note.User, error) {
	return &note.User{}, fmt.Errorf("_ NOT IMPLEMENTED _")
}

func (us *UserService) GetByUsername(ctx context.Context, id int) (*note.User, error) {
	return &note.User{}, fmt.Errorf("_ NOT IMPLEMENTED _")
}
