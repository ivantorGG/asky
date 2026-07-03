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

	r.Get("/events", h.EventsPage) //загрузка страницы с событиями(просто шаблон)
	r.With(middleware.Auth).Post("/events/new", h.CreateEvent) //создание нового события
	r.With(middleware.Auth).Get("/api/events/teacher", h.ListTeachersEvents)//получение списка событий учителя
	r.With(middleware.Auth).Delete("/events/{code}", h.DeleteEventByCode)//удаление события по коду
	r.Get("/events/{code}/teacher", h.TeacherEventPage)//загрузка страницы события по коду(шаблон) УУЧИТЕЛЬ
	r.Get("/events/{code}/student", h.StudentEventPage)//загрузка страницы события по коду(шаблон) УЧЕНИК
	r.Get("/api/events/{code}/link", h.GetEventLink)//получение ссылки на событие
	r.Get("/api/events/{code}/qrcode", h.GetEventQRcode)//получение QR-кода события
	r.Get("/api/events/student", h.ListUsersEvents)//получение списка событий студента

	r.Get("/events/{code}/questions", h.GetQuestionsByEventCode)//получение вопросов события по коду с ОТСОРТИРОВАННЫМИ вопросами по лайкам
	
	r.Get("/register", h.RegistrationPage)
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	r.Get("/login", h.LoginPage)
	
	r.Post("/events/{code}/question", h.NewQuestion)
	
	r.Put("/questions/{id}/vote", h.Vote)
	r.Delete("/questions/{id}/vote", h.UnVote)
	r.Put("/questions/{id}/answer", h.AnswerQuestion)

	return r
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	fs := http.StripPrefix(path, http.FileServer(root))
	r.Handle(path+"/*", fs)
}
