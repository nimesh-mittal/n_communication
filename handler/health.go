package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// HealthHandler represents interface for health of a service
type HealthHandler interface {
	Home(w http.ResponseWriter, r *http.Request)
	Health(w http.ResponseWriter, r *http.Request)
	NewHealthRouter() http.Handler
}

type healthHandler struct{}

// NewHealthHandler creates new object of HealthHandler
func NewHealthHandler() HealthHandler {
	return &healthHandler{}
}

// NewHealthRouter constructs new router for health endpoints
func (hh *healthHandler) NewHealthRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", hh.Home)
	r.Get("/_health", hh.Health)
	return r
}

func (hh *healthHandler) Home(w http.ResponseWriter, r *http.Request) {
	type home struct {
		Greet string
	}

	h := home{Greet: "hello"}

	res, _ := json.Marshal(h)
	w.Write(res)
}

func (hh *healthHandler) Health(w http.ResponseWriter, r *http.Request) {
	type health struct {
		Status string
	}

	h := health{Status: "green"}

	res, _ := json.Marshal(h)
	w.Write(res)
}
