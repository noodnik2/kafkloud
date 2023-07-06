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
		//Handler:      http.HandlerFunc(nonCancelingSlowHandler), // TRY THIS ONE
		Handler: http.HandlerFunc(cancelingSlowHandler), // OR THIS ONE
		ConnState: func(conn net.Conn, state http.ConnState) {
			log.Printf("connection state change; conn(%v), state(%v)\n", conn.RemoteAddr(), state)
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
		generateSlowResponse(w)
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

func nonCancelingSlowHandler(w http.ResponseWriter, _ *http.Request) {
	generateSlowResponse(w)
}

func generateSlowResponse(w http.ResponseWriter) {
	log.Printf("entering response handler sleep for %d\n", parms.responseDelay)
	writeResponse(w, "response1 before sleep\n")
	// Simulate a long-running task
	time.Sleep(parms.responseDelay)
	writeResponse(w, "response2 after sleep\n")
	log.Printf("leaving response handler\n")
}

func writeResponse(w http.ResponseWriter, message string) {
	log.Printf("sending a response(%s)\n", message)
	n, err := w.Write([]byte(message))
	w.(http.Flusher).Flush()
	log.Printf("return from send is n(%d) err(%v)\n", n, err)
}
