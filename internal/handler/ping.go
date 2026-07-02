// internal/handler/ping.go

package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

type Handler struct {
	DB     *sql.DB
    Logger *log.Logger
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "pong")
}
