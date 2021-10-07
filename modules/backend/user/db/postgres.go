package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/iagapie/go-spring/modules/backend/user"
	"github.com/iagapie/go-spring/modules/sys/postgresdb"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
	"time"
)

var _ user.Storage = &storage{}

type storage struct {
	db  *postgresdb.Database
	log *logrus.Entry
}

func NewStorage(postgres *postgresdb.Database, log *logrus.Entry) user.Storage {
	return &storage{
		db:  postgres,
		log: log,
	}
}

func (s *storage) FindByEmail(ctx context.Context, email string) (user.User, error) {
	return s.findOne(ctx, "email", email)
}

func (s *storage) FindByUUID(ctx context.Context, uuid string) (user.User, error) {
	return s.findOne(ctx, "uuid", uuid)
}

func (s *storage) Create(ctx context.Context, model user.User) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.db.WithContext(ctx).Create(&model).Error; err != nil {
		if strings.Contains(err.Error(), "23505") {
			return user.ErrRecordConflict
		}
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	s.log.Tracef("Created user: %s.\n", model.UUID)

	return nil
}

func (s *storage) findOne(ctx context.Context, column, value string) (user.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var model user.User
	if err := s.db.WithContext(ctx).First(&model, fmt.Sprintf("%s = ?", column), value).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model, user.ErrRecordNotFound
		}
		return model, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return model, nil
}
