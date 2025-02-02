package server

import (
	"github.com/carinfinin/shortener-url/internal/app/config"
	"github.com/carinfinin/shortener-url/internal/app/router"
	"github.com/carinfinin/shortener-url/internal/app/storage"
	"github.com/carinfinin/shortener-url/internal/app/storage/storefile"
	"net/http"
)

type Server struct {
	Addr   string
	Store  storage.Repositories
	Router *router.Router
}

func New(config *config.Config) (*Server, error) {

	s, err := storefile.New(config.FilePath)

	if err != nil {
		return nil, err
	}
	return &Server{
		Addr:   config.Addr,
		Store:  s,
		Router: router.ConfigureRouter(s, config),
	}, nil
}

func (s *Server) Start() error {
	return http.ListenAndServe(s.Addr, s.Router.Handle)
}
