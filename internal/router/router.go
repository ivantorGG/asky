package router

import (
	"net/http"

	"asky/internal/handler"

	"github.com/go-chi/chi/v5"
)

func New(h *handler.Handler) http.Handler {
	r := chi.NewRouter()

	r.Get("/ping", h.Ping)
	//r.Get("/events/new", h.NewEvents)
	r.Post("/events/new", h.CreateEvent)
	r.Get("/events", h.ListEvents)

	r.Post("/register", h.Register)
	r.Post("/login", h.Login)

	fileServer := http.FileServer(http.Dir("./tests/frontend"))

	// /* ловит все остальные запросы (например, /, /app.js) и отдает их как файлы
	r.Handle("/*", fileServer)

	return r
}
