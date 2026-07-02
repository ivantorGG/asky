package router

import (
	"net/http"

	"asky/internal/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func New(h *handler.Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:8080", "http://localhost:5500"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/ping", h.Ping)
	FileServer(r, "/static", http.Dir("./web/static"))

	r.Get("/events", h.ListEvents)
	r.Post("/events/new", h.CreateEvent)
	r.Delete("/events/{code}", h.DeleteEventByCode)
	r.Get("/events/{code}", h.EventsPage)

	r.Get("/register", h.RegistrationPage)
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	r.Get("/login", h.LoginPage)

	r.Post("/events/{code}/question", h.NewQuestion)

	r.Put("/questions/{id}/vote", h.Vote)
	r.Delete("/questions/{id}/vote", h.UnVote)

	return r
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	fs := http.StripPrefix(path, http.FileServer(root))
	r.Handle(path+"/*", fs)
}
