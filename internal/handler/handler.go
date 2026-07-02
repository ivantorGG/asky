// internal/handler/handler.go

package handler

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"asky/internal/service"
)

type Handler struct {
	Auth *service.AuthService
}

func New(db *pgxpool.Pool) *Handler {
	return &Handler{
		Auth: service.NewAuthService(db),
	}
}
