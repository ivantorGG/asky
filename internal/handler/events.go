package handler

import (
	"asky/internal/middleware"
	"asky/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"asky/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/skip2/go-qrcode"
)

var cfg *config.Config = config.Load()

type EventRequest struct {
	Title string `json:"title"`
}
type TeacherEvent struct {
	Title string `json:"title"`
	Code  string `json:"code"`
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
	userID := r.Context().Value(middleware.UserIDKey).(int64)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad_request")
		return
	}
	if req.Title == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad_request")
		return
	}
	_, err := h.DB.Exec(
		r.Context(),
		`INSERT INTO events(title, owner_id)
		VALUES ($1,$2)`,
		req.Title,
		userID,
	)

	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "db_error: "+err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "event_created",
	})
}

func (h *Handler) ListTeachersEvents(w http.ResponseWriter, r *http.Request) {
	// Безопасно достаем userID учителя из контекста
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	rows, err := h.DB.Query(
		r.Context(),
		`SELECT title, code
		 FROM events
		 WHERE owner_id = $1 AND is_active = TRUE
		 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "server_error")
		return
	}
	defer rows.Close()

	events := []TeacherEvent{}
	for rows.Next() {
		var e TeacherEvent
		if err := rows.Scan(&e.Title, &e.Code); err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "server_error")
			return
		}
		events = append(events, e)
	}

	w.Header().Set("Content-Type", "application/json")
	if len(events) == 0 {
		json.NewEncoder(w).Encode([]TeacherEvent{})
		return
	}
	json.NewEncoder(w).Encode(events)
}

func (h *Handler) ListUsersEvents(w http.ResponseWriter, r *http.Request) {
	// Студенты смотрят историю через КУКИ, так как у них нет аккаунтов
	cookie, err := r.Cookie("visited_events")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]TeacherEvent{})
		return
	}

	decoded, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad_cookie")
		return
	}

	var codes []string
	if err := json.Unmarshal([]byte(decoded), &codes); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad_cookie")
		return
	}

	// Выбираем из базы только те события, UUID которых сохранены в куках студента
	rows, err := h.DB.Query(
		r.Context(),
		`SELECT title, code
		 FROM events
		 WHERE code = ANY($1::uuid[]) AND is_active = TRUE
		 ORDER BY created_at DESC`,
		codes,
	)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "server_error")
		return
	}
	defer rows.Close()

	events := []TeacherEvent{}
	for rows.Next() {
		var e TeacherEvent
		if err := rows.Scan(&e.Title, &e.Code); err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "server_error")
			return
		}
		events = append(events, e)
	}

	w.Header().Set("Content-Type", "application/json")
	if len(events) == 0 {
		json.NewEncoder(w).Encode([]TeacherEvent{})
		return
	}
	json.NewEncoder(w).Encode(events)
}

func (h *Handler) GetQuestionsByEventCode(w http.ResponseWriter, r *http.Request) {
	eventCode := chi.URLParam(r, "code")
	if eventCode == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad_request")
		return
	}
	rows, err := h.DB.Query(
		r.Context(),
		`SELECT id, event_code, text, likes, answered, created_at 
		FROM questions 
		WHERE event_code = $1::uuid 
		ORDER BY answered ASC, likes DESC`,
		eventCode,
	)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "server_error: "+err.Error())
		return
	}
	defer rows.Close()
	questions := make([]QuestionResponse, 0)
	for rows.Next() {
		var q QuestionResponse
		err := rows.Scan(&q.ID, &q.EventCode, &q.Text, &q.Likes, &q.Answered, &q.CreatedAt)
		if err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "server_error: "+err.Error())
			return
		}
		questions = append(questions, q)
	}
	if err = rows.Err(); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "server_error: "+err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(questions)
}

func (h *Handler) AnswerQuestion(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad_request")
		return
	}

	cmdTag, err := h.DB.Exec(
		r.Context(),
		`UPDATE questions SET answered = TRUE WHERE id = $1`,
		id,
	)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "db_error")
		return
	}

	if cmdTag.RowsAffected() == 0 {
		utils.WriteJSONError(w, http.StatusNotFound, "question_not_found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteEventByCode(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)

	code := chi.URLParam(r, "code")
	if code == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad_request")
		return
	}

	cmdTag, err := h.DB.Exec(
		r.Context(),
		`UPDATE events 
	SET is_active = FALSE 
	WHERE code = $1::uuid AND owner_id = $2 AND is_active = TRUE`,
		code,
		userID,
	)

	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "server_error: "+err.Error())
		return
	}

	if cmdTag.RowsAffected() == 0 {
		utils.WriteJSONError(w, http.StatusNotFound, "event_not_found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "event_deleted",
	})
}

func (h *Handler) EventsPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/templates/eventList.html")
}

func (h *Handler) StudentEventPage(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	var events []string

	if cookie, err := r.Cookie("visited_events"); err == nil {
		decoded, err := url.QueryUnescape(cookie.Value)
		if err == nil {
			_ = json.Unmarshal([]byte(decoded), &events)
		}
	}

	exists := false
	for _, c := range events {
		if c == code {
			exists = true
			break
		}
	}

	if !exists {
		events = append(events, code)

		data, _ := json.Marshal(events)

		http.SetCookie(w, &http.Cookie{
			Name:   "visited_events",
			Value:  url.QueryEscape(string(data)),
			Path:   "/",
			MaxAge: 60 * 60 * 24 * 365,
		})
	}

	http.ServeFile(w, r, "./web/templates/studentWinds.html")
}
func (h *Handler) TeacherEventPage(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	var events []string

	if cookie, err := r.Cookie("visited_events"); err == nil {
		decoded, err := url.QueryUnescape(cookie.Value)
		if err == nil {
			_ = json.Unmarshal([]byte(decoded), &events)
		}
	}

	exists := false
	for _, c := range events {
		if c == code {
			exists = true
			break
		}
	}

	if !exists {
		events = append(events, code)

		data, _ := json.Marshal(events)

		http.SetCookie(w, &http.Cookie{
			Name:   "visited_events",
			Value:  url.QueryEscape(string(data)),
			Path:   "/",
			MaxAge: 60 * 60 * 24 * 365,
		})
	}

	http.ServeFile(w, r, "./web/templates/teacherWinds.html")
}
func (h *Handler) GetEventLink(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad_request")
		return
	}
	link := cfg.Domain + "/events/" + code

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"link": link,
	})
}
func (h *Handler) GetEventQRcode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad_request")
		return
	}
	link := fmt.Sprintf("%s/events/%s", cfg.Domain, code)

	png, err := qrcode.Encode(link, qrcode.Medium, 256)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "server_error: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	w.Write(png)
}
