// internal/middleware/auth_middleware.go

package middleware

import (
	"strconv"
	"net/http"
	"strings"

	"asky/internal/jwt"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		if auth == "" {
			http.Error(w, "no token", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		claims, err := jwt.Parse(tokenStr)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		r.Header.Set("user_id", strconv.FormatInt(claims.UserID, 10))

		next(w, r)
	}
}