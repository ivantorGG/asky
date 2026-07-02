package router

import (
	"net/http"

	"asky/internal/handler"

	"github.com/go-chi/chi/v5"
)

func New(h *handler.Handler) http.Handler {
	r := chi.NewRouter()

	r.Get("/ping", h.Ping)

	r.Post("/register", h.Register)
	r.Post("/login", h.Login)

	return r
}