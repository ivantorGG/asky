package router

import (
	"net/http"

	"asky/internal/handler"
	"asky/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func New(h *handler.Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:8080", "http://localhost:8080", "http://localhost:5500"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/ping", h.Ping)
	FileServer(r, "/static", http.Dir("./web/static"))

	// =========================================================================
	// Pages (HTML)
	// =========================================================================

	r.Get("/login", h.LoginPage)           // Login page
	r.Get("/register", h.RegistrationPage) // Registration page

	r.Get("/events", h.EventsPage)                      // Events dashboard page
	r.Get("/events/{code}/teacher", h.TeacherEventPage) // Teacher view of an event
	r.Get("/events/{code}/student", h.StudentEventPage) // Student view of an event

	// =========================================================================
	// Authentication API
	// =========================================================================

	r.Post("/api/register", h.Register) // Register a new user
	r.Post("/api/login", h.Login)       // Authenticate user

	// =========================================================================
	// Events API
	// =========================================================================

	r.With(middleware.Auth).Post("/api/events", h.CreateEvent)                // Create a new event
	r.With(middleware.Auth).Get("/api/events/teacher", h.ListTeachersEvents)  // Get teacher's events
	r.Get("/api/events/student", h.ListUsersEvents)                           // Get recently visited events
	r.With(middleware.Auth).Delete("/api/events/{code}", h.DeleteEventByCode) // Delete (deactivate) an event

	r.Get("/api/events/{code}/link", h.GetEventLink)     // Get event invitation link
	r.Get("/api/events/{code}/qrcode", h.GetEventQRcode) // Get event QR code

	// =========================================================================
	// Questions API
	// =========================================================================

	r.Get("/api/events/{code}/questions", h.GetQuestionsByEventCode) // Get all event questions
	r.Post("/api/events/{code}/questions", h.NewQuestion)            // Create a new question

	r.With(middleware.Auth).Put("/api/questions/{id}/answer", h.AnswerQuestion) // Mark question as answered

	r.Put("/api/questions/{id}/vote", h.Vote)      // Upvote a question
	r.Delete("/api/questions/{id}/vote", h.UnVote) // Remove question vote

	// =========================================================================
	// Comments API
	// =========================================================================

	r.Get("/api/questions/{id}/comments", h.ListQuestionComments) // Get question comments
	r.Post("/api/questions/{id}/comments", h.NewComment)          // Add a comment to a question

	return r
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	fs := http.StripPrefix(path, http.FileServer(root))
	r.Handle(path+"/*", fs)
}
