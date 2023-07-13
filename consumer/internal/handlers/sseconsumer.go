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

	rid := nextRequestId()

	acceptType := r.Header[textproto.CanonicalMIMEHeaderKey("accept")]
	if len(acceptType) < 1 || acceptType[0] != "text/event-stream" {
		log.Printf("(%d) NOTE: improper 'Accept' mime-type in request(%v) ignored; assuming 'text/event-stream'", rid, acceptType)
	}

	topicsParam := r.URL.Query().Get("topics")
	topics := strings.Split(topicsParam, ",")
	chp := &consumerHandlerProcessor{
		rid: rid,
		ctx: r.Context(),
		h:   h,
		w:   w,
	}
	kConsumer := h.KConsumerFactory(chp, topics)

	w.Header().Set("Content-Type", "text/event-stream")
	//w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	var err error
	if err = chp.emit("retry: 3000"); err != nil {
		log.Printf("(%d) SSE connection could not be completed: %s\n", rid, err)
		return
	}

	if err = kConsumer.Launch(); err != nil {
		log.Printf("(%d) HTTP consumer could not be launched: %s\n", rid, err)
		return
	}

	log.Printf("(%d) waiting for HTTP consumer to finish...\n", rid)
	select {
	case <-chp.ctx.Done():
		log.Printf("(%d) HTTP context timed out at\n", rid)
		if kcErr := kConsumer.Close(); kcErr != nil {
			log.Printf("(%d) error closing kConsumer: %s\n", rid, kcErr)
		}
	case <-kConsumer.DoneChan:
		log.Printf("(%d) Kafka consumer exited\n", rid)
	}
	log.Printf("(%d) HTTP consumer finished\n", rid)
}

type consumerHandlerProcessor struct {
	rid uint64
	ctx context.Context
	h   *Handlers
	w   http.ResponseWriter
}

func (chp *consumerHandlerProcessor) emit(event string) error {
	rc := http.NewResponseController(chp.w)
	err := rc.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		log.Printf("(%d) unable to set write deadline for event(%s): %s\n", chp.rid, event, err)
		return err
	}

	if err := chp.h.respondWithString(chp.w, event); err != nil {
		log.Printf("(%d) unable to deliver event(%s): %s\n", chp.rid, event, err)
		return err
	}

	if err := rc.Flush(); err != nil {
		log.Printf("(%d) unable to flush event(%s): %s\n", chp.rid, event, err)
		return err
	}
	log.Printf("(%d) delivered event(%s)\n", chp.rid, event)
	return nil
}

func (chp *consumerHandlerProcessor) ProcessMessage(topic, message string) error {
	timeStamp := time.Now().Format(time.RFC3339)
	return chp.emit(fmt.Sprintf("id: %s-%d-%s\nevent: streamer\ndata: %s\n\n", topic, chp.rid, timeStamp, message))
}
