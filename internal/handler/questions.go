package handler

import (
	"asky/internal/utils"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"html"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

var validEventCode = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

type NewQuestionRequest struct {
	Text string `json:"text"`
}

func (h *Handler) NewQuestion(w http.ResponseWriter, r *http.Request) {
	var req NewQuestionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad_request")
		return
	}

	req.Text = strings.TrimSpace(req.Text)
	req.Text = html.EscapeString(req.Text)
	if req.Text == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad_request")
		return
	}

	code := strings.TrimSpace(chi.URLParam(r, "code"))
	if !validEventCode.MatchString(code) {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid_event_code")
		return
	}

	cookie, err := r.Cookie("visitor_id")
	var visitorID string
	validVisitorID := false

	if err == nil && cookie != nil {
		if _, decodeErr := hex.DecodeString(cookie.Value); decodeErr == nil && len(cookie.Value) == 64 {
			visitorID = cookie.Value
			validVisitorID = true
		}
	}

	if !validVisitorID {
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
	}

	var last time.Time

	err = h.DB.QueryRow(
		r.Context(),
		`
    SELECT created_at
    FROM questions
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

	_, err = h.DB.Exec(
		r.Context(),
		`INSERT INTO questions(event_code, text, visitor_id)
		 VALUES ($1, $2, $3)`,
		code,
		req.Text,
		visitorID,
	)

	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "db_error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "question_created",
	})
}
