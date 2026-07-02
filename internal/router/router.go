package router

import (
	"net/http"

	"asky/internal/handler"

	"github.com/go-chi/chi/v5"
)

func New(h *handler.Handler) http.Handler {
	r := chi.NewRouter()

	r.Get("/ping", h.Ping)
	r.Post("/events/new", h.CreateEvent)
	r.Get("/events", h.ListEvents)
	//r.Get("/events/{code}", h.GetEventByCode)
	r.Delete("/events/{code}", h.DeleteEventByCode)
	r.Post("/register", h.Register)
	r.Get("/register", h.OpenRegistration)
	r.Post("/login", h.Login)

	r.Post("/events/{code}/question", h.NewQuestion)

	r.Put("/questions/{id}/vote", h.Vote)
	r.Delete("/questions/{id}/vote", h.UnVote)

	return r
}
