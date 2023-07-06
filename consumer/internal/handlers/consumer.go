package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"
)

func (h *Handlers) ConsumerHandler(w http.ResponseWriter, r *http.Request) {
	deadline, ok := r.Context().Deadline()
	log.Printf("consumer handler called; deadline(%s), ok(%t)\n", deadline.Format(time.RFC3339), ok)
	topicsParam := r.URL.Query().Get("topics")
	topics := strings.Split(topicsParam, ",")
	chp := &consumerHandlerProcessor{
		ctx: r.Context(),
		h:   h,
		w:   w,
	}
	kConsumer := h.KConsumerFactory(chp, topics)

	setChunkedResponseHeaders(w)

	if err := kConsumer.Launch(); err != nil {
		log.Printf("HTTP consumer could not be launched: %s\n", err)
		h.errHandler(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("waiting for HTTP consumer to finish...\n")
	select {
	case <-chp.ctx.Done():
		log.Printf("request timed out at\n")
		if kcErr := kConsumer.Close(); kcErr != nil {
			log.Printf("error closing kConsumer: %s\n", kcErr)
		}
	case <-kConsumer.DoneChan:
		log.Printf("request completed\n")
	}
	log.Printf("HTTP consumer finished\n")
}

type consumerHandlerProcessor struct {
	ctx context.Context
	h   *Handlers
	w   http.ResponseWriter
}

func (chp *consumerHandlerProcessor) ProcessMessage(topic, message string) error {
	if err := chp.h.respondWithBytesAndFlush(chp.w, []byte(message)); err != nil {
		log.Printf("unable to deliver message from topic(%s): %s\n", topic, err)
		return err
	}
	log.Printf("delivered message from topic(%s)\n", topic)
	return nil
}
