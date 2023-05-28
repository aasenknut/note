package mock

import (
	"context"

	"github.com/aasenknut/note"
)

type UserService struct {
	createFn      func(ctx context.Context, username, password string) (int, error)
	getByUsnameFn func(ctx context.Context, username string) (*note.User, error)
	getByIDFn     func(ctx context.Context, id int) (*note.User, error)
}

func (us *UserService) Create(ctx context.Context, username, password string) (int, error) {
	return us.createFn(ctx, username, password)
}

func (us *UserService) GetByUsername(ctx context.Context, username string) (*note.User, error) {
	return us.getByUsnameFn(ctx, username)
}

func (us *UserService) GetByID(ctx context.Context, id int) (*note.User, error) {
	return us.getByIDFn(ctx, id)
}
