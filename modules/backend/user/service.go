package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/iagapie/go-spring/modules/sys/password"
)

type (
	Service interface {
		GetByEmailAndPassword(ctx context.Context, email, password string) (User, error)
		GetByUUID(ctx context.Context, uuid string) (User, error)
		Create(ctx context.Context, dto CreateUserDTO) (string, error)
	}

	service struct {
		storage Storage
		encoder password.Encoder
	}
)

func NewService(storage Storage, encoder password.Encoder) Service {
	return &service{
		storage: storage,
		encoder: encoder,
	}
}

func (s *service) GetByEmailAndPassword(ctx context.Context, email, password string) (User, error) {
	if user, err := s.storage.FindByEmail(ctx, email); err == nil && s.encoder.IsValid(user.Password, password) {
		return user, nil
	}

	return User{}, ErrRecordNotFound
}

func (s *service) GetByUUID(ctx context.Context, uuid string) (User, error) {
	return s.storage.FindByUUID(ctx, uuid)
}

func (s *service) Create(ctx context.Context, dto CreateUserDTO) (string, error) {
	encoded, err := s.encoder.Encode(dto.Password)
	if err != nil {
		return "", fmt.Errorf("failed to create user. error: %w", err)
	}
	model := NewUser(dto)
	model.Password = encoded
	if err = s.storage.Create(ctx, model); err != nil {
		if errors.Is(err, ErrRecordConflict) {
			return "", err
		}
		return "", fmt.Errorf("failed to create user. error: %w", err)
	}
	return model.UUID, nil
}
