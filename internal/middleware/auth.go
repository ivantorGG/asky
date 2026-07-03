package middleware

import (
	"context"
	"net/http"
	"strconv"
)

type ContextKey string

const UserIDKey ContextKey = "userID"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userIDStr := r.Header.Get("X-User-ID")
		if userIDStr == "" {
			writeJSONError(w, http.StatusInternalServerError, "server_error")
			return
		}

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "server_error")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func writeJSONError(w http.ResponseWriter, i int, s string) {
	panic("unimplemented")
}
