package handlers

import (
	"encoding/json"
	"net/http"
)

type ConsumerStatus struct {
	NConsumed int `json:"n_consumed"`
}

func (h *Handlers) StatusHandler(w http.ResponseWriter, _ *http.Request) {
	cs := ConsumerStatus{
		NConsumed: h.Counter.GetNConsumed(),
	}
	marshal, errJsonMarshal := json.Marshal(&cs)
	if errJsonMarshal != nil {
		h.errHandler(w, http.StatusInternalServerError, errJsonMarshal)
		return
	}
	setJsonResponseHeaders(w)
	_ = h.respondWithString(w, string(marshal))
}
