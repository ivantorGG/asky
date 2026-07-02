package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type NewQuestionRequest struct {
	Text string `json:"text"`
}

func (h *Handler) NewQuestion(w http.ResponseWriter, r *http.Request) {
	var req NewQuestionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	code := chi.URLParam(r, "code")

	_, err := h.DB.Exec(
		r.Context(),
		`INSERT INTO questions(event_code, text)
		 VALUES ($1, $2)`,
		code,
		req.Text,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "question created")
}

