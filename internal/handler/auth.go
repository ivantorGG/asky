package handler

import (
	"asky/internal/utils"
	"encoding/json"
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
		utils.WriteJSONError(w, http.StatusBadRequest, "bad_request")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "server_error")
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
		utils.WriteJSONError(w, http.StatusBadRequest, "db_error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "registration_success",
	})
}

func (h *Handler) RegistrationPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/templates/speakerReg.html")
}

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/templates/speakerLog.html")
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "bad_request")
		return
	}

	var passwordHash string

	err := h.DB.QueryRow(
		r.Context(),
		`SELECT password_hash FROM users WHERE email = $1`,
		req.Email,
	).Scan(&passwordHash)

	if err != nil {
		utils.WriteJSONError(w, http.StatusUnauthorized, "bad_credentials")
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(passwordHash),
		[]byte(req.Password),
	)

	if err != nil {
		utils.WriteJSONError(w, http.StatusUnauthorized, "bad_credentials")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "login_success",
	})
}
