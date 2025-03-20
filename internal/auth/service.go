package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger"
	"github.com/x0k/skillrock-tasks-service/internal/shared"
	"golang.org/x/crypto/bcrypt"
)

type UsersRepo interface {
	SaveUser(ctx context.Context, user *User) error
	UserByLogin(ctx context.Context, login string) (*User, error)
}

type Service struct {
	log           *logger.Logger
	secret        []byte
	tokenLifetime time.Duration
	repo          UsersRepo
}

func NewService(log *logger.Logger, secret []byte, tokenLifetime time.Duration, repo UsersRepo) *Service {
	return &Service{log, secret, tokenLifetime, repo}
}

func (s *Service) Register(ctx context.Context, login string, password string) (string, *shared.ServiceError) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", shared.NewUnexpectedError(err, "failed to generate password hash")
	}
	user := NewUser(login, passwordHash)
	if err := s.repo.SaveUser(ctx, user); err != nil {
		if errors.Is(err, ErrLoginIsTaken) {
			return "", shared.NewServiceError(err, fmt.Sprintf("%q login is already taken", login))
		}
		return "", shared.NewUnexpectedError(err, "failed to create a new user")
	}
	return s.issueAccessToken(login)
}

func (s *Service) Login(ctx context.Context, login string, password string) (string, *shared.ServiceError) {
	user, err := s.repo.UserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return "", shared.NewServiceError(err, "failed to login")
		}
		return "", shared.NewUnexpectedError(err, "failed to login")
	}
	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", shared.NewServiceError(ErrPasswordsMismatch, "failed to login")
		}
		return "", shared.NewUnexpectedError(err, "failed to login")
	}
	return s.issueAccessToken(login)
}

func (s *Service) issueAccessToken(login string) (string, *shared.ServiceError) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": login,
		"exp": time.Now().Add(s.tokenLifetime).Unix(),
	})
	accessToken, err := t.SignedString(s.secret)
	if err != nil {
		return "", shared.NewUnexpectedError(err, "failed to sign access token")
	}
	return accessToken, nil
}
