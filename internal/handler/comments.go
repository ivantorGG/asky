package handler

import (
	"asky/internal/utils"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type NewCommentRequest struct {
	Text string `json:"text"`
}

func (h *Handler) NewComment(w http.ResponseWriter, r *http.Request) {
	var req NewCommentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad_request")
		return
	}

	id := chi.URLParam(r, "id")

	_, err := h.DB.Exec(
		r.Context(),
		`INSERT INTO comments(question_id, text)
		 VALUES ($1, $2)`,
		id,
		req.Text,
	)

	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "db_error")
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
