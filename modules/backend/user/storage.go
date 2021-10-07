package user

import (
	"context"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("user not found")
	ErrRecordConflict = errors.New("user already exists")
)

type Storage interface {
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByUUID(ctx context.Context, uuid string) (User, error)
	Create(ctx context.Context, model User) error
}
