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
	query := `INSERT INTO user (username, password) VALUES (?, ?)`

	res, err := tx.ExecContext(ctx, query, username, password)
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

func (us *UserService) GetByUsername(ctx context.Context, username string) (*note.User, error) {
	tx, err := us.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return nil, fmt.Errorf("get user by username: %v", err)
	}
	query := `SELECT id, username, password FROM user WHERE username = ?`
	row := tx.QueryRowContext(ctx, query, username)
	var user note.User
	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		return nil, fmt.Errorf("get user by username: %v", err)
	}
	return &note.User{}, nil
}
