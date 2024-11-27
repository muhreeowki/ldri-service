package main

import (
	"log"
	"net/http"
)

type Server struct {
	listenAddress string
	store         *Store
	queue         *Queue
}

func NewServer(listenAddress string, store *Store, queue *Queue) *Server {
	return &Server{
		listenAddress: listenAddress,
		store:         store,
		queue:         queue,
	}
}

func (s *Server) ListenAndServer() error {
	router := http.NewServeMux()
	router.HandleFunc("/health", s.makeHandlerFunc(s.healthCheckHandler))
	return http.ListenAndServe(s.listenAddress, router)
}

func (s *Server) healthCheckHandler(w http.ResponseWriter, r *http.Request) error {
	return WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type ServerHandler func(w http.ResponseWriter, r *http.Request) error

func (s *Server) makeHandlerFunc(h ServerHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			log.Printf("Error occured: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
