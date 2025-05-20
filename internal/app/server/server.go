package server

import (
	"context"
	"github.com/carinfinin/shortener-url/internal/app/config"
	"github.com/carinfinin/shortener-url/internal/app/router"
	"github.com/carinfinin/shortener-url/internal/app/service"
	"github.com/carinfinin/shortener-url/internal/app/storage/store"
	"github.com/carinfinin/shortener-url/internal/app/storage/storefile"
	"github.com/carinfinin/shortener-url/internal/app/storage/storepg"
	"net/http"
)

// Server заускает сервер и содержит ссылку на хранилище.
type Server struct {
	Addr   string
	Store  service.Repository
	Router *router.Router
}

// New конструктор для Server принимает кофиг.
func New(config *config.Config) (*Server, error) {

	var server Server

	switch {
	case config.DBPath != "":
		s, err := storepg.New(config)
		if err != nil {
			return nil, err
		}
		s.CreateTableForDB(context.Background())
		server.Store = s

	case config.FilePath != "":
		s, err := storefile.New(config)
		if err != nil {
			return nil, err
		}
		server.Store = s

	default:
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

// Start запускает сервер.
func (s *Server) Start() error {
	return http.ListenAndServe(s.Addr, s.Router.Handle)
}
