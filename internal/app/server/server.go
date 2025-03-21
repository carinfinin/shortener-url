package server

import (
	"context"
	"github.com/carinfinin/shortener-url/internal/app/config"
	"github.com/carinfinin/shortener-url/internal/app/router"
	"github.com/carinfinin/shortener-url/internal/app/service"
	"github.com/carinfinin/shortener-url/internal/app/storage"
	"github.com/carinfinin/shortener-url/internal/app/storage/store"
	"github.com/carinfinin/shortener-url/internal/app/storage/storefile"
	"github.com/carinfinin/shortener-url/internal/app/storage/storepg"
	"net/http"
)

type Server struct {
	Addr   string
	Store  storage.Repository
	Router *router.Router
}

func New(config *config.Config) (*Server, error) {

	var server Server

	if config.DBPath != "" {
		s, err := storepg.New(config)
		if err != nil {
			return nil, err
		}
		s.CreateTableForDB(context.Background())
		server.Store = s

	} else if config.FilePath != "" {
		s, err := storefile.New(config)
		if err != nil {
			return nil, err
		}
		server.Store = s

	} else {
		s, err := store.New(config)
		if err != nil {
			return nil, err
		}
		server.Store = s

	}

	server.Addr = config.Addr
	service := service.New(server.Store, config)

	server.Router = router.ConfigureRouter(service, config)

	return &server, nil
}

func (s *Server) Start() error {
	return http.ListenAndServe(s.Addr, s.Router.Handle)
}
