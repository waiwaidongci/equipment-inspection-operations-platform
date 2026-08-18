package application

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/example/inspection-platform/internal/domain"
	"github.com/example/inspection-platform/internal/platform/auth"
)

type Service struct {
	users     domain.UserRepository
	jwtSecret string
	tokenTTL  time.Duration
}

func NewService(users domain.UserRepository, jwtSecret string, tokenTTL time.Duration) *Service {
	return &Service{users: users, jwtSecret: jwtSecret, tokenTTL: tokenTTL}
}

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResult struct {
	Token string      `json:"token"`
	User  domain.User `json:"user"`
}

func (s *Service) Login(ctx context.Context, input LoginInput) (LoginResult, error) {
	username := strings.TrimSpace(input.Username)
	user, err := s.users.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return LoginResult{}, fmt.Errorf("%w: username or password is incorrect", domain.ErrInvalid)
		}
		return LoginResult{}, err
	}
	if !auth.CheckPassword(user.PasswordHash, input.Password) {
		return LoginResult{}, fmt.Errorf("%w: username or password is incorrect", domain.ErrInvalid)
	}
	claims := auth.Claims{
		Sub:      fmt.Sprintf("%d", user.ID),
		Username: user.Username,
		Exp:      time.Now().Add(s.tokenTTL).Unix(),
	}
	return LoginResult{Token: auth.Sign(s.jwtSecret, claims), User: *user}, nil
}

func (s *Service) EnsureDefault(ctx context.Context, username, password, displayName string) error {
	_, err := s.users.FindByUsername(ctx, username)
	if err == nil {
		return nil
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return err
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		return fmt.Errorf("hash default password: %w", err)
	}
	now := time.Now().UTC()
	user := &domain.User{
		Username:     username,
		PasswordHash: hash,
		DisplayName:  displayName,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.users.Create(ctx, user); err != nil {
		return fmt.Errorf("create default user: %w", err)
	}
	return nil
}

func (s *Service) List(ctx context.Context) ([]domain.User, error) {
	return s.users.List(ctx)
}

func (s *Service) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	return s.users.FindByID(ctx, id)
}
