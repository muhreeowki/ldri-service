package main

import (
	"log"
	"net/http"
)

type HTTPServer struct {
	listenAddress string
	store         *Store
	queue         *RabbitQueue
}

func NewServer(listenAddress string, store *Store, queue *RabbitQueue) *HTTPServer {
	return &HTTPServer{
		listenAddress: listenAddress,
		store:         store,
		queue:         queue,
	}
}

func (s *HTTPServer) ListenAndServer() error {
	router := http.NewServeMux()
	router.HandleFunc("/health", s.makeHandlerFunc(s.healthCheckHandler))
	return http.ListenAndServe(s.listenAddress, router)
}

func (s *HTTPServer) healthCheckHandler(w http.ResponseWriter, r *http.Request) error {
	return WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type ServerHandler func(w http.ResponseWriter, r *http.Request) error

func (s *HTTPServer) makeHandlerFunc(h ServerHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			log.Printf("Error occured: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
