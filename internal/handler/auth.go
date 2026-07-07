package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"asky/internal/utils"

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

	// Минимальная валидация на пустые поля
	if req.Email == "" || req.Password == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid_input")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "server_error")
		return
	}

	_, err = h.DB.Exec(
		r.Context(),
		`INSERT INTO users(email, password_hash) VALUES ($1, $2)`,
		req.Email,
		string(hash),
	)

	if err != nil {
		// Ошибка чаще всего означает, что такой email уже зарегистрирован (Unique Constraint)
		utils.WriteJSONError(w, http.StatusConflict, "email_already_exists")
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

	var userID int64
	var passwordHash string

	// КРИТИЧНО: Вытягиваем id вместе с хэшем пароля, чтобы связать с сессией
	err := h.DB.QueryRow(
		r.Context(),
		`SELECT id, password_hash FROM users WHERE email = $1`,
		req.Email,
	).Scan(&userID, &passwordHash)

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

	// 1. Генерируем криптографически безопасный случайный токен (32 байта -> 64 символа)
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "server_error")
		return
	}
	sessionToken := hex.EncodeToString(b)

	// 2. Сохраняем сессию в базу данных.
	// Предполагается, что в вашей таблице sessions есть колонки token и user_id.
	// Если у вас есть колонка expires_at, вы можете добавить время жизни сессии в запрос.
	_, err = h.DB.Exec(
		r.Context(),
		`INSERT INTO sessions (token, user_id) VALUES ($1, $2)`,
		sessionToken,
		userID,
	)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "server_error")
		return
	}

	// 3. Устанавливаем куку сессии для браузера. Название совпадает с r.Cookie в middleware.
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour * 30), // Время жизни куки на 30 дней
		HttpOnly: true,                                // Защита от кражи токена через JS (XSS)
		Secure:   false,                               // Поставьте true, если тестируете локально через HTTPS или в продакшене
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "login_success",
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// 1. Пытаемся достать куку сессии
	cookie, err := r.Cookie("session_token")
	if err == nil {
		// 2. Если кука есть, удаляем запись о сессии из базы данных
		_, _ = h.DB.Exec(
			r.Context(),
			`DELETE FROM sessions WHERE token = $1`,
			cookie.Value,
		)
	}

	// 3. Перезаписываем куку на клиенте (ставим MaxAge: -1 для мгновенного удаления)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Инвалидирует куку в браузере
		HttpOnly: true,
		Secure:   false, // Поставьте true в продакшене (HTTPS)
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "logout_success",
	})
}
