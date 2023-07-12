package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/textproto"
	"strings"
	"time"
)

func (h *Handlers) SseConsumerHandler(w http.ResponseWriter, r *http.Request) {

	acceptType := r.Header[textproto.CanonicalMIMEHeaderKey("accept")]
	if len(acceptType) < 1 || acceptType[0] != "text/event-stream" {
		log.Printf("NOTE: improper 'Accept' mime-type in request(%v) ignored; assuming 'text/event-stream'", acceptType)
	}

	topicsParam := r.URL.Query().Get("topics")
	topics := strings.Split(topicsParam, ",")
	chp := &consumerHandlerProcessor{
		ctx: r.Context(),
		h:   h,
		w:   w,
	}
	kConsumer := h.KConsumerFactory(chp, topics)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	var err error
	if err = chp.emit("retry: 3000"); err != nil {
		log.Printf("SSE connection could not be completed: %s\n", err)
		return
	}

	if err = kConsumer.Launch(); err != nil {
		log.Printf("HTTP consumer could not be launched: %s\n", err)
		return
	}

	log.Printf("waiting for HTTP consumer to finish...\n")
	select {
	case <-chp.ctx.Done():
		log.Printf("HTTP context timed out at\n")
		if kcErr := kConsumer.Close(); kcErr != nil {
			log.Printf("error closing kConsumer: %s\n", kcErr)
		}
	case <-kConsumer.DoneChan:
		log.Printf("Kafka consumer exited\n")
	}
	log.Printf("HTTP consumer finished\n")
}

type consumerHandlerProcessor struct {
	ctx context.Context
	h   *Handlers
	w   http.ResponseWriter
}

func (chp *consumerHandlerProcessor) emit(event string) error {
	rc := http.NewResponseController(chp.w)
	err := rc.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		log.Printf("unable to set write deadline for event(%s): %s\n", event, err)
		return err
	}

	if err := chp.h.respondWithString(chp.w, event); err != nil {
		log.Printf("unable to deliver event(%s): %s\n", event, err)
		return err
	}

	if err := rc.Flush(); err != nil {
		log.Printf("unable to flush event(%s): %s\n", event, err)
		return err
	}
	log.Printf("delivered event(%s)\n", event)
	return nil
}

func (chp *consumerHandlerProcessor) ProcessMessage(topic, message string) error {
	timeStamp := time.Now().Format(time.RFC3339)
	return chp.emit(fmt.Sprintf("id: %s-%s\nevent: streamer\ndata: %s\n\n", topic, timeStamp, message))
}
