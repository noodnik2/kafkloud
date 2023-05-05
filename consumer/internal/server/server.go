package server

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/noodnik2/incubator20/k8s/kafkloud/consumer/internal/configs"
	"github.com/noodnik2/incubator20/k8s/kafkloud/consumer/internal/controller"
	"github.com/noodnik2/incubator20/k8s/kafkloud/consumer/internal/handlers"
	"github.com/noodnik2/incubator20/k8s/kafkloud/consumer/internal/util"
)

type Server struct {
	sconfig    *configs.ServerConfig
	controller controller.Controller
	srv        *http.Server
	ctx        context.Context
}

func NewServer(ctx context.Context, sconfig *configs.ServerConfig, controller controller.Controller) *Server {
	return &Server{
		sconfig:    sconfig,
		controller: controller,
		ctx:        ctx,
	}
}

func (s *Server) Launch(util.ComponentErrorHandler) error {

	h := handlers.NewHandlers(s.controller)
	r := mux.NewRouter()

	// set content type
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("serving %s\n", r.URL.String())
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			next.ServeHTTP(w, r)
		})
	})

	r.HandleFunc("/_health", h.HealthHandler).Methods(http.MethodGet)
	r.HandleFunc("/status", h.StatusHandler).Methods(http.MethodGet)
	http.Handle("/", r)

	s.srv = &http.Server{
		Addr:         s.sconfig.Addr,
		WriteTimeout: s.sconfig.WriteTimeout,
		ReadTimeout:  s.sconfig.ReadTimeout,
		IdleTimeout:  s.sconfig.IdleTimeout,
		Handler:      r,
	}

	go func() {
		err := s.srv.ListenAndServe()
		log.Printf("Server stopped; err=%v...\n", err)
	}()

	log.Printf("launched Server on(%s)\n", s.srv.Addr)
	return nil
}

func (s *Server) Close() error {
	log.Printf("shutting down Server...\n")
	if errShutdown := s.srv.Shutdown(s.ctx); errShutdown != nil {
		log.Printf("server shutdown error: %v\n", errShutdown)
	}
	return nil
}
