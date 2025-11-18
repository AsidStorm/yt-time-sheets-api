package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Server struct {
	Router  *mux.Router
	Stopped bool
}

type health struct {
	Status bool `json:"status"`
}

func (s *Server) Run(address string) error {
	return http.ListenAndServe(address, s.Router)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	response := &health{
		Status: s.Stopped,
	}

	j, _ := json.Marshal(response)

	if s.Stopped {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Write(j)
}

func (s *Server) init() {
	s.Router.Use(loggingMiddleware)
	s.Router.HandleFunc("/health", s.health).Methods(http.MethodGet)
}

func NewServer(r *mux.Router) *Server {
	s := &Server{
		Router: r,
	}

	s.init()

	return s
}

func (s *Server) Stop() {
	s.Stopped = true
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := r.Header.Get("X-Request-Id")

		if requestId == "" {
			requestId = uuid.New().String()
			r.Header.Set("X-Request-Id", requestId)
		}

		next.ServeHTTP(w, r)
	})
}
