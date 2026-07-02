// internal/repository/user_repository.go

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID           int64
	Email        string
	PasswordHash string
}

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, email, hash string) error {
	_, err := r.db.Exec(ctx,
		"INSERT INTO users(email, password_hash) VALUES ($1, $2)",
		email, hash,
	)
	return err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	u := &User{}

	err := r.db.QueryRow(ctx,
		"SELECT id, email, password_hash FROM users WHERE email=$1",
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash)

	if err != nil {
		return nil, err
	}

	return u, nil
}
