package handler

import (
	"asky/internal/middleware"
	"encoding/json"
	"fmt"
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

type QuestionResponse struct {
	ID        int64     `json:"id"`
	EventCode string    `json:"event_code"`
	Text      string    `json:"text"`
	Likes     int       `json:"likes"`
	Answered  bool      `json:"answered"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req EventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad_request")
		return
	}
	if req.Title == "" {
		writeJSONError(w, http.StatusBadRequest, "bad_request")
		return
	}
	_, err := h.DB.Exec(
		r.Context(),
		`INSERT INTO events(title, owner_id)
		 VALUES ($1,$2)`,
		req.Title,
		req.OwnerID,
	)

	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "db_error")
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "event_created",
	})
}

func (h *Handler) ListEvents(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)
	rows, err := h.DB.Query(
		r.Context(),
		`SELECT id, title, code, owner_id, is_active, created_at 
         FROM events 
         WHERE owner_id = $1 AND is_active = TRUE
         ORDER BY created_at DESC`,
		userID,
	)
	fmt.Println("Rows:", rows)      // Debugging line
	fmt.Println("OwnerID:", userID) // Debugging line
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "server_error")
		return
	}
	defer rows.Close()

	events := []Event{}
	for rows.Next() {
		var e Event
		err := rows.Scan(&e.ID, &e.Title, &e.Code, &e.OwnerID, &e.IsActive, &e.CreatedAt)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "server_error")
			return
		}
		events = append(events, e)
	}

	if err = rows.Err(); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "server_error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}

func (h *Handler) GetQuestionsByEventCode(w http.ResponseWriter, r *http.Request) {
	eventCode := chi.URLParam(r, "code")
	if eventCode == "" {
		writeJSONError(w, http.StatusBadRequest, "bad_request")
		return
	}
	rows, err := h.DB.Query(
		r.Context(),
		`SELECT id, event_code, text, likes, answered, created_at 
		 FROM questions 
		 WHERE event_code = $1::uuid 
		 ORDER BY answered ASC, likes DESC, created_at DESC`,
		eventCode,
	)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "server_error")
		return
	}
	defer rows.Close()
	questions := make([]QuestionResponse, 0)
	for rows.Next() {
		var q QuestionResponse
		err := rows.Scan(&q.ID, &q.EventCode, &q.Text, &q.Likes, &q.Answered, &q.CreatedAt)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "server_error")
			return
		}
		questions = append(questions, q)
	}
	if err = rows.Err(); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "server_error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(questions)
}

func (h *Handler) DeleteQuestionByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "bad_request")
		return
	}

	cmdTag, err := h.DB.Exec(
		r.Context(),
		`UPDATE questions SET answered = TRUE WHERE id = $1`,
		id,
	)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "db_error")
		return
	}

	if cmdTag.RowsAffected() == 0 {
		writeJSONError(w, http.StatusNotFound, "question_not_found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteEventByCode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		writeJSONError(w, http.StatusBadRequest, "bad_request")
		return
	}
	var req EventRequest

	cmdTag, err := h.DB.Exec(
		r.Context(),
		`UPDATE events 
     SET is_active = FALSE 
     WHERE code = $1::uuid AND owner_id = $2 AND is_active = TRUE`,
		code,
		req.OwnerID,
	)

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "server_error")
		return
	}

	if cmdTag.RowsAffected() == 0 {
		writeJSONError(w, http.StatusNotFound, "event_not_found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "event_deleted",
	})
}
func (h *Handler) EventsPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/templates/eventList.html")
}
