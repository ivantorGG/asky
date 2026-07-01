// internal/handler/ping.go

package handler

import (
	"fmt"
	"net/http"
)

type Handler struct {
	// позже сюда добавим DB
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "pong")
}
