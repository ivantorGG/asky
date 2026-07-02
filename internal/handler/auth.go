package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad_request")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "server_error")
		return
	}

	_, err = h.DB.Exec(
		r.Context(),
		`INSERT INTO users(email, password_hash)
		 VALUES ($1,$2)`,
		req.Email,
		string(hash),
	)

	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "db_error")
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "registration success")
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "bad_request")
		return
	}

	var passwordHash string

	err := h.DB.QueryRow(
		r.Context(),
		`SELECT password_hash FROM users WHERE email = $1`,
		req.Email,
	).Scan(&passwordHash)

	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "bad_creditants")
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(passwordHash),
		[]byte(req.Password),
	)

	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "bad_creditants")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "login success")
}
