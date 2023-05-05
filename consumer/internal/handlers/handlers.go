package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/noodnik2/incubator20/k8s/kafkloud/consumer/internal/controller"
)

type Handlers struct {
	controller controller.Controller
}

type ConsumerStatus struct {
	NConsumed int `json:"n_consumed"`
}

func NewHandlers(controller controller.Controller) *Handlers {
	return &Handlers{
		controller: controller,
	}
}

func (h *Handlers) StatusHandler(writer http.ResponseWriter, _ *http.Request) {
	cs := ConsumerStatus{
		NConsumed: h.controller.GetNConsumed(),
	}
	marshal, errJsonMarshal := json.Marshal(&cs)
	if errJsonMarshal != nil {
		h.errHandler(writer, http.StatusInternalServerError, errJsonMarshal)
		return
	}
	h.respondWithString(writer, string(marshal))
}

func (h *Handlers) HealthHandler(writer http.ResponseWriter, _ *http.Request) {
	h.respondWithString(writer, "OK!\n")
}

func (h *Handlers) respondWithString(writer http.ResponseWriter, s string) {
	if err := writeResponse(writer, []byte(s)); err != nil {
		h.errHandler(writer, http.StatusInternalServerError, err)
	}
}

func (h *Handlers) errHandler(writer http.ResponseWriter, status int, err error) {
	writer.WriteHeader(status)
	_ = writeResponse(writer, []byte(err.Error()))
}

func writeResponse(writer http.ResponseWriter, response []byte) error {
	_, err := writer.Write(response)
	return err
}
