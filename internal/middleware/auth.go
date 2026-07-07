package middleware

import (
	"context"
	"net/http"

	"asky/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ContextKey string

const UserIDKey ContextKey = "userID"

func Auth(db *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Получаем куку сессии
			cookie, err := r.Cookie("session_token")
			if err != nil {
				utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			// 2. Ищем userID по токену в базе данных
			var userID int64
			err = db.QueryRow(r.Context(),
				"SELECT user_id FROM sessions WHERE token = $1",
				cookie.Value).Scan(&userID)

			if err != nil {
				// Сессия не найдена или истекла
				utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			// 3. Передаем userID дальше по цепочке
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
