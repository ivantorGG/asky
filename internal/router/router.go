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
	//r.Get("/events/new", h.NewEvents)
	r.Post("/events/new", h.CreateEvent)
	r.Get("/events", h.ListEvents)

	r.Post("/register", h.Register)
	r.Post("/login", h.Login)

	// временно убираем /me (пока нет метода)
	_ = middleware.Auth
	// Включаем раздачу статических файлов из папки tests/frontend
	// Метод chi.NewFileServer заставит роутер отдавать index.html по корневому пути "/"
	fileServer := http.FileServer(http.Dir("./tests/frontend"))
	
	// Символ /* означает, что Chi будет отдавать любые файлы из этой папки (HTML, JS, CSS)
	r.Handle("/*", fileServer) 
	return r
}
