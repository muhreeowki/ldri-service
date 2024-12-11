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
	router.HandleFunc("GET /health", makeHandlerFunc(s.healthCheckHandler))
	router.HandleFunc("POST /login", makeHandlerFunc(s.loginHandler))
	router.HandleFunc("POST /signup", makeHandlerFunc(s.signUpHandler))
	router.HandleFunc("GET /data", makeHandlerFunc(s.getDataHandler))
	log.Printf("Listening on %s...\n", s.listenAddress)
	return http.ListenAndServe(s.listenAddress, router)
}

// ----------------- HANDLERS -----------------

func (s *HTTPServer) healthCheckHandler(w http.ResponseWriter, r *http.Request) error {
	log.Printf("Health check request")
	return WriteJSON(w, http.StatusOK, map[string]string{"message": "ok"})
}

func (s *HTTPServer) loginHandler(w http.ResponseWriter, r *http.Request) error {
	log.Printf("Login request")
	return WriteJSON(w, http.StatusOK, map[string]string{"message": "logged in"})
}

func (s *HTTPServer) signUpHandler(w http.ResponseWriter, r *http.Request) error {
	log.Printf("Sign up request")
	return WriteJSON(w, http.StatusOK, map[string]string{"message": "signed up"})
}

func (s *HTTPServer) getDataHandler(w http.ResponseWriter, r *http.Request) error {
	log.Printf("Get data request")
	return WriteJSON(w, http.StatusOK, map[string]string{"data": "some odk data"})
}

// ----------------- UTILS -----------------

type ServerHandler func(w http.ResponseWriter, r *http.Request) error

func makeHandlerFunc(h ServerHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			log.Printf("Error occured: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
