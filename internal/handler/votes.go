package handler

import (
	"fmt"
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
	fmt.Fprint(w, "question "+questionID+" liked")
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
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	rows := res.RowsAffected()
	if rows == 0 {
		http.Error(w, "question not found", http.StatusNotFound)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "question "+questionID+" disliked")
}


