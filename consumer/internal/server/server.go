package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/noodnik2/incubator20/k8s/kafkloud/consumer/internal/handlers"
	"github.com/noodnik2/incubator20/k8s/kafkloud/consumer/internal/util"
)

type Server struct {
	Addr         string
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	IdleTimeout  time.Duration
	context.Context
	*handlers.Handlers
	util.ComponentErrorHandler
	*http.Server
}

func (s *Server) Launch(util.ComponentErrorHandler) error {

	h := s.Handlers
	r := mux.NewRouter()

	r.HandleFunc("/consume", h.ConsumerHandler).Methods(http.MethodGet)
	r.HandleFunc("/_health", h.HealthHandler).Methods(http.MethodGet)
	r.HandleFunc("/status", h.StatusHandler).Methods(http.MethodGet)
	http.Handle("/", r)

	s.Server = &http.Server{
		Addr:         s.Addr,
		WriteTimeout: s.WriteTimeout,
		ReadTimeout:  s.ReadTimeout,
		IdleTimeout:  s.IdleTimeout,
		Handler:      r,
		ErrorLog:     log.Default(),
		ConnState: func(conn net.Conn, state http.ConnState) {
			log.Printf("connection state change; conn(%v), state(%v)\n", conn.RemoteAddr(), state)
		},
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			log.Printf("returning updated context for connection\n")
			timeout, _ := context.WithTimeout(ctx, s.WriteTimeout)
			return timeout
		},
	}

	log.Printf("server timeouts: write(%v), read(%v), idle(%v)\n", s.Server.WriteTimeout, s.Server.ReadTimeout, s.Server.IdleTimeout)
	go func() {
		err := s.Server.ListenAndServe()
		log.Printf("Server stopped; err=%v...\n", err)
	}()

	log.Printf("launched Server on(%s)\n", s.Server.Addr)
	return nil
}

func (s *Server) Close() error {
	log.Printf("shutting down Server...\n")
	if errShutdown := s.Server.Shutdown(s.Context); errShutdown != nil {
		log.Printf("server shutdown error: %v\n", errShutdown)
	}
	return nil
}
