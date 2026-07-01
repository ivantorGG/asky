// internal/service/auth_service.go

package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"asky/internal/repository"
	"asky/internal/jwt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthService struct {
	repo *repository.UserRepository
}

func NewAuthService(db *pgxpool.Pool) *AuthService {
	return &AuthService{
		repo: repository.NewUserRepository(db),
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	return s.repo.CreateUser(ctx, email, string(hash))
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := jwt.Generate(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
