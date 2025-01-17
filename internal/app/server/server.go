package server

import (
	"github.com/carinfinin/shortener-url/internal/app/config"
	"github.com/carinfinin/shortener-url/internal/app/router"
	"github.com/carinfinin/shortener-url/internal/app/storage"
	"github.com/carinfinin/shortener-url/internal/app/storage/store"
	"net/http"
)

type Server struct {
	Addr   string
	Store  storage.Repositories
	Router *router.Router
}

func New(config *config.Config) *Server {
	s := store.New()

	return &Server{
		Addr:   config.Addr,
		Store:  s,
		Router: router.ConfigureRouter(s, config.URL),
	}
}

func (s *Server) Start() error {
	return http.ListenAndServe(s.Addr, s.Router.Handle)
}
