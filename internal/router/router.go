package router

import (
	"net/http"

	"asky/internal/handler"
	"asky/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func New(h *handler.Handler) http.Handler {
	r := chi.NewRouter()

	r.Get("/ping", h.Ping)

	r.Post("/register", h.Register)
	r.Post("/login", h.Login)

	// временно убираем /me (пока нет метода)
	_ = middleware.Auth

	return r
}