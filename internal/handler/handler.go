package handler

import "github.com/jackc/pgx/v5/pgxpool"

type Handler struct {
	DB *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Handler {
	return &Handler{
		DB: db,
	}
}
