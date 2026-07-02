package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) Vote(w http.ResponseWriter, r *http.Request) {
	questionID := chi.URLParam(r, "id")

	res, err := h.DB.Exec(
		r.Context(),
		`UPDATE questions
		 SET likes = likes + 1
		 WHERE id = $1`,
		questionID,
	)

	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	rows := res.RowsAffected()
	if rows == 0 {
		http.Error(w, "question not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "question "+questionID+" liked",
	})
}

func (h *Handler) UnVote(w http.ResponseWriter, r *http.Request) {
	questionID := chi.URLParam(r, "id")

	res, err := h.DB.Exec(
		r.Context(),
		`UPDATE questions
		 SET likes = likes - 1
		 WHERE id = $1`,
		questionID,
	)

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "server_error")
		return
	}

	rows := res.RowsAffected()
	if rows == 0 {
		writeJSONError(w, http.StatusNotFound, "question_not_found")
		return
	}
	
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "question "+questionID+" disliked",
	})
}


