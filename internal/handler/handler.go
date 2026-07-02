// internal/handler/handler.go

package handler

import (
	"asky/internal/service"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	Auth *service.AuthService
	DB   *pgxpool.Pool
    Logger *log.Logger
}

func New(db *pgxpool.Pool) *Handler {
	return &Handler{
		Auth: service.NewAuthService(db),
		DB:     db, // <-- ВОТ ТЕПЕРЬ МЫ ЕГО ЗАПИСЫВАЕМ!
		Logger: log.Default(), // на всякий случай инициализируем дефолтный логгер
	}
}
