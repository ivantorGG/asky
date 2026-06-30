package router

import (
	"net/http"

	"asky/internal/handler"

	"github.com/go-chi/chi/v5"
)

func New() http.Handler {
	r := chi.NewRouter()

	r.Get("/ping", handler.Ping)

	return r
}