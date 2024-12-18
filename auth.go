package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// AuthServer is the main struct for the Auth server
type AuthServer struct {
	listenAddress string
	store         Store
	logger        *log.Logger
}

// NewAuthServer creates a new AuthServer instance
func NewAuthServer(listenAddress string, store Store) *AuthServer {
	return &AuthServer{
		listenAddress: listenAddress,
		store:         store,
		logger:        log.New(os.Stdout, "[auth-server] ", log.LstdFlags),
	}
}

// ListenAndServer starts the AuthServer and listens for incoming requests
func (s *AuthServer) ListenAndServer() error {
	router := http.NewServeMux()
	router.HandleFunc("GET /health", makeHandlerFunc(s.healthCheckHandler))
	router.HandleFunc("POST /login", makeHandlerFunc(s.loginHandler))
	router.HandleFunc("POST /signup", makeHandlerFunc(s.signUpHandler))
	router.HandleFunc("POST /validate", makeHandlerFunc(s.validateHandler))
	router.HandleFunc("GET /data", makeHandlerFunc(s.getDataHandler))
	s.logger.Printf("Listening on %s...\n", s.listenAddress)
	return http.ListenAndServe(s.listenAddress, router)
}

// ----------------- HANDLERS -----------------

// healthCheckHandler is a simple health check handler
func (s *AuthServer) healthCheckHandler(w http.ResponseWriter, r *http.Request) error {
	s.logger.Printf("Health check request")
	return WriteJSON(w, http.StatusOK, map[string]string{"message": "ok"})
}

// loginHandler handles the login request and returns a JWT token if the user is valid
func (s *AuthServer) loginHandler(w http.ResponseWriter, r *http.Request) error {
	// Get the user from the request body
	user := &User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		WriteJSON(w, http.StatusInternalServerError, "invalid request")
		return fmt.Errorf("unable to decode request body")
	}

	// Check if the user has provided the email and password
	if user.Email == "" || user.Password == "" {
		WriteJSON(w, http.StatusInternalServerError, fmt.Errorf("missing credentials"))
		return fmt.Errorf("User %s failed to log in", user.Email)
	}

	// Check if the user exists in the store
	dbUser, err := s.store.GetUser(user.Email, user.Password)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, "invalid credentials")
		return err
	}

	// Return a success message and a jwt token
	token, err := CreateJWT(dbUser)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, err.Error())
		return err
	}

	log.Printf("User %s logged in with token: %s\n", user.Email, token)
	return WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

// validateHandler validates the JWT token and returns a success message if the token is valid
func (s *AuthServer) validateHandler(w http.ResponseWriter, r *http.Request) error {
	// Get the token from the request
	token := r.Header.Get("Authorization")
	if token == "" {
		err := fmt.Errorf("missing token")
		WriteJSON(w, http.StatusUnauthorized, err)
		return err
	}

	// Check if the token is a bearer token
	if string(token[:7]) != "Bearer " {
		err := fmt.Errorf("invalid authorization header")
		WriteJSON(w, http.StatusUnauthorized, err)
		return err
	}
	token = token[7:]

	// Verify the token
	_, err := ValidateJWT(token)
	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, err.Error())
		return err
	}

	return WriteJSON(w, http.StatusOK, "valid token")
}

// signUpHandler handles the sign up request and creates a new user in the store
func (s *AuthServer) signUpHandler(w http.ResponseWriter, r *http.Request) error {
	log.Printf("Sign up request")

	reqUserData := new(User)
	if err := json.NewDecoder(r.Body).Decode(reqUserData); err != nil {
		WriteJSON(w, http.StatusOK, map[string]string{"error": "invalid request"})
		return fmt.Errorf("Error decoding request body")
	}

	usr, err := s.store.CreateUser(reqUserData.Email, reqUserData.Password)
	if err != nil {
		WriteJSON(w, http.StatusOK, map[string]string{"error": "error creating user"})
		return err
	}

	token, err := CreateJWT(usr)
	if err != nil {
		WriteJSON(w, http.StatusOK, map[string]string{"error": "error creating user"})
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

// getDataHandler handles the get data request and returns the data
func (s *AuthServer) getDataHandler(w http.ResponseWriter, r *http.Request) error {
	log.Printf("Get data request")
	return WriteJSON(w, http.StatusOK, map[string]string{"data": "some odk data"})
}

// ----------------- HANDLER UTILS -----------------

// ServerHandler is a custom handler type that returns an error
type ServerHandler func(w http.ResponseWriter, r *http.Request) error

// makeHandlerFunc is a utility function that converts a ServerHandler to an http.HandlerFunc
func makeHandlerFunc(h ServerHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			log.Printf("Error occured: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
