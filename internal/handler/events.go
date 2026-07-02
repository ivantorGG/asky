package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type EventRequest struct {
	Title   string `json:"title"`
	OwnerID int64  `json:"owner_id"`
}

type Event struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Code      string    `json:"code"`
	OwnerID   int64     `json:"owner_id"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req EventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if req.Title == "" {
		http.Error(w, "Название не может быть пустым", http.StatusBadRequest)
		return
	}
	userID := h.GetCurrentUserID(r)
	_, err := h.DB.Exec(
		r.Context(),
		`INSERT INTO events(title, owner_id)
		 VALUES ($1,$2)`,
		req.Title,
		userID,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) ListEvents(w http.ResponseWriter, r *http.Request) {
	// Пока у нас нет авторизации, зашиваем owner_id = 1, как и в CreateEvent
	var ownerID int64 = h.GetCurrentUserID(r)

	rows, err := h.DB.Query(
		r.Context(),
		`SELECT id, title, code, owner_id, is_active, created_at 
         FROM events 
         WHERE owner_id = $1 AND is_active = TRUE
         ORDER BY created_at DESC`,
		ownerID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	events := []Event{} // инициализируем сразу как пустой слайс, чтобы в JSON не было null
	for rows.Next() {
		var e Event
		err := rows.Scan(&e.ID, &e.Title, &e.Code, &e.OwnerID, &e.IsActive, &e.CreatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		events = append(events, e)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}

func (h *Handler) GetEventByCode(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement this function to retrieve an event by its code
}
func (h *Handler) DeleteEventByCode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		http.Error(w, "Код комнаты не указан", http.StatusBadRequest)
		return
	}
	userID := h.GetCurrentUserID(r)

	cmdTag, err := h.DB.Exec(
		r.Context(),
		`UPDATE events 
     SET is_active = FALSE 
     WHERE code = $1::uuid AND owner_id = $2 AND is_active = TRUE`,
		code,
		userID,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if cmdTag.RowsAffected() == 0 {
		http.Error(w, "Комната не найдена или у вас нет прав на её удаление", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetCurrentUserID(r *http.Request) int64 {
	return 1
}
