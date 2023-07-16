package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"
)

var parms = struct {
	serverAddr    string
	responseDelay time.Duration
	readTimeout   time.Duration
	writeTimeout  time.Duration
	idleTimeout   time.Duration
}{
	serverAddr:    ":8080",
	responseDelay: 30 * time.Second,
	readTimeout:   10 * time.Second,
	writeTimeout:  5 * time.Second,
	idleTimeout:   10 * time.Second,
}

/**
Experimental program used to play with and understand better the "timeout"
mechanisms available for a server within the Go standard library "http" APIs.
*/

func main() {

	var server = &http.Server{
		Addr:         parms.serverAddr,
		ReadTimeout:  parms.readTimeout,
		WriteTimeout: parms.writeTimeout,
		IdleTimeout:  parms.idleTimeout,
		Handler:      http.HandlerFunc(nonCancelingSlowHandler), // TRY THIS ONE
		//Handler: http.HandlerFunc(cancelingSlowHandler), // OR THIS ONE
		ConnState: func(conn net.Conn, state http.ConnState) {
			log.Printf("connection state change; conn(%v), state(%v)\n", conn.RemoteAddr(), state)
		},
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			cancelingContext, _ := context.WithTimeout(ctx, parms.writeTimeout)
			return cancelingContext
		},
	}

	log.Printf("running: %+v\n", parms)
	log.Fatal(server.ListenAndServe())
}

func cancelingSlowHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("entering canceling response handler sleep for %d\n", parms.responseDelay)

	ctx, cancel := context.WithTimeout(r.Context(), parms.writeTimeout*time.Second)
	defer cancel()

	go func() {
		generateSlowResponse(w, r)
		select {
		case <-ctx.Done():
			log.Println("request canceled")
		default:
			log.Println("response was not canceled")
		}
		log.Println("leaving goroutine")
	}()

	select {
	case <-ctx.Done():
		log.Println("connection timed out")
	case <-time.After(parms.responseDelay):
		log.Println("response was sent too late")
	}
	log.Println("leaving handler")
}

func nonCancelingSlowHandler(w http.ResponseWriter, r *http.Request) {
	// Below is to play with the newer (since go 1.20) "ResponseController" API
	// (hint: doesn't seem to work as we expect it until now)
	//rc := http.NewResponseController(w)
	//// Set a write deadline in 5 seconds time.
	//err := rc.SetWriteDeadline(time.Now().Add(parms.writeTimeout * time.Second))
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	_, _ = w.Write([]byte(err.Error()))
	//	return
	//}
	generateSlowResponse(w, r)
}

func generateSlowResponse(w http.ResponseWriter, r *http.Request) {
	log.Printf("entering response handler sleep for %d\n", parms.responseDelay)
	writeResponse(w, "response1 before sleep\n")
	select {
	case <-r.Context().Done():
		log.Printf("request context timed out\n")
		return
	case <-time.After(parms.responseDelay):
		log.Printf("simulated long-running task completed\n")
	}
	writeResponse(w, "response2 after sleep\n")
	log.Printf("leaving response handler\n")
}

func writeResponse(w http.ResponseWriter, message string) {
	log.Printf("sending a response(%s)\n", message)
	n, err := w.Write([]byte(message))
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	} else {
		log.Printf("NOTE: response writer could not be flushed\n")
	}
	log.Printf("return from send is n(%d) err(%v)\n", n, err)
}
