package handler

import (
	"net/http"
)

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/templates/index.html")
}
