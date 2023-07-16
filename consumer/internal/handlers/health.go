package handlers

import (
	"net/http"
)

func (h *Handlers) HealthHandler(w http.ResponseWriter, _ *http.Request) {
	setJsonResponseHeaders(w)
	_ = h.respondWithString(w, "OK!\n")
}
