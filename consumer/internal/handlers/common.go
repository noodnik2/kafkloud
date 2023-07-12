package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/noodnik2/incubator20/k8s/kafkloud/consumer/internal/controller"
)

type Handlers struct {
	*controller.Counter
	controller.KConsumerFactory
}

func (h *Handlers) respondWithString(w http.ResponseWriter, s string) error {
	return h.respondWithBytes(w, []byte(s))
}

func (h *Handlers) respondWithBytes(w http.ResponseWriter, bytes []byte) error {
	if err := writeResponse(w, bytes); err != nil {
		h.errHandler(w, http.StatusInternalServerError, err)
		return err
	}
	return nil
}

func (h *Handlers) errHandler(w http.ResponseWriter, status int, err error) {
	log.Printf("errHandler(%d): %s\n", status, err)
	w.WriteHeader(status)
	if respErr := writeResponse(w, []byte(err.Error())); respErr != nil {
		log.Printf("could not deliver error to client: %s\n", respErr)
	}
}

func writeResponse(w http.ResponseWriter, r []byte) error {
	n, err := w.Write(r)
	if err != nil {
		return err
	}
	contentLength := len(r)
	if n != contentLength {
		return fmt.Errorf("couldn't not deliver all content; nWritten(%d) != contentLength(%d)\n", n, contentLength)
	}
	return nil
}

func setJsonResponseHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}
