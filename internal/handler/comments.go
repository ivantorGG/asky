package handler

import (
	"asky/internal/utils"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"html"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type NewCommentRequest struct {
	Text string `json:"text"`
}

func (h *Handler) NewComment(w http.ResponseWriter, r *http.Request) {
	var req NewCommentRequest

	// 1. Декодируем JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad_request")
		return
	}

	// 2. Проверка на пустой текст
	if req.Text == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "empty_text")
		return
	}

	// 2.1. Экранируем HTML в тексте, чтобы скрипты не могли быть вставлены
	req.Text = html.EscapeString(req.Text)

	// 3. Получаем ID вопроса из URL
	idStr := chi.URLParam(r, "id")

	// 4. Конвертируем ID в нужный тип (выберите подходящий вариант)

	// ВАРИАНТ A: Если question_id в БД — это INT
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid_question_id")
		return
	}

	// ВАРИАНТ B: Если question_id в БД — это UUID (раскомментируйте, если нужно)
	// id, err := uuid.Parse(idStr)
	// if err != nil {
	//     log.Printf("❌ Не удалось конвертировать ID в UUID: %v", err)
	//     utils.WriteJSONError(w, http.StatusBadRequest, "invalid_question_id")
	//     return
	// }

	// 5. Получаем или создаём visitor_id
	cookie, err := r.Cookie("visitor_id")
	var visitorID string

	if err != nil {
		b := make([]byte, 32)
		if _, err := rand.Read(b); err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "server_error")
			return
		}
		visitorID = hex.EncodeToString(b)

		http.SetCookie(w, &http.Cookie{
			Name:     "visitor_id",
			Value:    visitorID,
			Path:     "/",
			Expires:  time.Now().AddDate(1, 0, 0),
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})
	} else {
		visitorID = cookie.Value
	}

	var last time.Time

	err = h.DB.QueryRow(
		r.Context(),
		`
    SELECT created_at
    FROM comments
    WHERE visitor_id=$1
    ORDER BY created_at DESC
    LIMIT 1
    `,
		visitorID,
	).Scan(&last)

	if err == nil && time.Since(last) < 15*time.Second {
		utils.WriteJSONError(w, http.StatusTooManyRequests, "wait_15_seconds")
		return
	}

	// 6. Вставляем комментарий в БД

	_, err = h.DB.Exec(
		r.Context(),
		`INSERT INTO comments(question_id, text, visitor_id)
		 VALUES ($1, $2, $3)`,
		id,
		req.Text,
		visitorID,
	)

	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "db_error: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "comment_created",
	})
}

type Comment struct {
	ID        int64     `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *Handler) ListQuestionComments(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	rows, err := h.DB.Query(
		r.Context(),
		`SELECT id, text, created_at
		 FROM comments
		 WHERE question_id = $1
		 ORDER BY created_at ASC`,
		id,
	)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "db_error")
		return
	}
	defer rows.Close()

	comments := make([]Comment, 0)

	for rows.Next() {
		var c Comment

		if err := rows.Scan(
			&c.ID,
			&c.Text,
			&c.CreatedAt,
		); err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "db_error")
			return
		}

		comments = append(comments, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	cookie, err := r.Cookie("visitor_id")
	if err != nil {
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	_, err = h.DB.Exec(
		r.Context(),
		`DELETE FROM comments WHERE id = $1 AND visitor_id = $2`,
		id,
		cookie.Value,
	)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "db_error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "comment_deleted",
	})
}
func (h *Handler) EditComment(w http.ResponseWriter, r *http.Request) {
	var req NewCommentRequest
	cookie, err := r.Cookie("visitor_id")
	if err != nil {
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad_request")
		return
	}
	_, err = h.DB.Exec(
		r.Context(),
		`UPDATE comments SET text = $1 WHERE id = $2 AND visitor_id = $3`,
		req.Text,
		chi.URLParam(r, "id"),
		cookie.Value,
	)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "db_error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "comment_updated",
	})
}
