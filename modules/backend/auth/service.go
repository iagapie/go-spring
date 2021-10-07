package auth

import (
	"context"
	"fmt"
	"github.com/go-redis/cache/v8"
	"github.com/google/uuid"
	"github.com/iagapie/go-spring/modules/backend/user"
	"github.com/iagapie/go-spring/modules/sys/config"
	"github.com/iagapie/go-spring/modules/sys/token"
	"github.com/sirupsen/logrus"
)

type (
	Service interface {
		Auth(ctx context.Context, dto SignInDTO) (TokensResponse, error)
		RefreshToken(ctx context.Context, dto RefreshTokenDTO) (TokensResponse, error)
	}

	service struct {
		duration     config.JWTDuration
		userService  user.Service
		redisCache   *cache.Cache
		tokenManager token.Token
		log          *logrus.Entry
	}
)

func NewService(duration config.JWTDuration, userService user.Service, redisCache *cache.Cache, tokenManager token.Token, log *logrus.Entry) Service {
	return &service{
		duration:     duration,
		userService:  userService,
		redisCache:   redisCache,
		tokenManager: tokenManager,
		log:          log,
	}
}

func (s *service) Auth(ctx context.Context, dto SignInDTO) (TokensResponse, error) {
	u, err := s.userService.GetByEmailAndPassword(ctx, dto.Email, dto.Password)
	if err != nil {
		return TokensResponse{}, fmt.Errorf("authentication: %w", err)
	}
	return s.session(ctx, u.UUID)
}

func (s *service) RefreshToken(ctx context.Context, dto RefreshTokenDTO) (TokensResponse, error) {
	defer s.redisCache.Delete(ctx, dto.Token)

	var id string
	if err := s.redisCache.Get(ctx, dto.Token, &id); err != nil {
		return TokensResponse{}, fmt.Errorf("refresh token: %w", err)
	}
	if id == "" {
		return TokensResponse{}, fmt.Errorf("refresh token: %s not found", dto.Token)
	}

	return s.session(ctx, id)
}

func (s *service) session(ctx context.Context, id string) (TokensResponse, error) {
	accessToken, err := s.tokenManager.Create(s.duration.Access, id)
	if err != nil {
		return TokensResponse{}, fmt.Errorf("refresh token: %w", err)
	}
	refreshToken := uuid.NewString()
	if err = s.redisCache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   refreshToken,
		Value: id,
		TTL:   s.duration.Refresh,
	}); err != nil {
		return TokensResponse{}, fmt.Errorf("refresh token: %w", err)
	}

	return TokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
