package server

import (
	"github.com/carinfinin/shortener-url/internal/app/router"
	"github.com/carinfinin/shortener-url/internal/app/storage"
	"github.com/carinfinin/shortener-url/internal/app/storage/store"
	"net/http"
)

type Server struct {
	Port   string
	Store  storage.Repositories
	Router *router.Router
}

func New(port string) *Server {
	s := store.New()

	return &Server{
		Port:   ":" + port,
		Store:  s,
		Router: router.ConfigureRouter(s),
	}
}

func (s *Server) Start() error {
	return http.ListenAndServe(s.Port, s.Router.Handle)
}
