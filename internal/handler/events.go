package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Code      string    `json:"code"`
	OwnerID   int       `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *Handler) NewEvents(w http.ResponseWriter, r *http.Request) {
	filePath := "тут будет путь к файлу с формой для создания нового события"
	http.ServeFile(w, r, filePath)
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	title := r.PostFormValue("title")
	if title == "" {
		http.Error(w, "Название не может быть пустым", http.StatusBadRequest)
		return
	}
	eventCode := uuid.New().String()
	ownerID := h.GetCurrentUserID(r)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO events (title, code, owner_id, created_at) 
		VALUES ($1, $2, $3, NOW())
	`
	_, err := h.DB.ExecContext(ctx, query, title, eventCode, ownerID)
	if err != nil {
		h.Logger.Printf("Ошибка вставки ивента в БД: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/events/"+eventCode, http.StatusSeeOther)

}
func (h *Handler) GetCurrentUserID(r *http.Request) int {
	// TODO: Сюда мы позже вставим чтение сессии/кук:
	// session, _ := h.Sessions.Get(r, "session-name")
	// return session.Values["userID"].(int)

	return 1 // Пока просто возвращаем ID первого лектора, который у тебя есть в БД
}
