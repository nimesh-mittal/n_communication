package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// Server represents server interface
type Server interface {
	StartServer(address string)
	Register(path string, handle func(http.ResponseWriter, *http.Request), method string)
	Mount(path string, handler http.Handler)
}

type server struct {
	Router *chi.Mux
}

// New creates new server object
func New() Server {
	router := chi.NewRouter()

	router.Use(middleware.Timeout(3 * time.Second))
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RealIP)
	router.Use(middleware.RequestID)
	router.Use(middleware.Throttle(1000))
	router.Use(middleware.NoCache)
	router.Use(middleware.SetHeader("Content-Type", "application/json"))

	return &server{Router: router}
}

// StartServer starts HTTP server at given address
func (s *server) StartServer(address string) {
	zap.L().Info("started running server",
		zap.String("address", address))
	_ = http.ListenAndServe(address, s.Router)
}

// Register registers a path with the server
func (s *server) Register(path string, handle func(http.ResponseWriter, *http.Request), method string) {
	s.Router.MethodFunc(method, path, handle)
}

// Mount registers a path with the server
func (s *server) Mount(path string, handler http.Handler) {
	s.Router.Mount(path, handler)
}
