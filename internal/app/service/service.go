package service

import (
	"context"
	"github.com/carinfinin/shortener-url/internal/app/auth"
	"github.com/carinfinin/shortener-url/internal/app/config"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	"github.com/carinfinin/shortener-url/internal/app/models"
	"github.com/carinfinin/shortener-url/internal/app/storage"
	"github.com/carinfinin/shortener-url/internal/app/storage/storepg"
	"strings"
	"time"
)

type Service struct {
	Store  storage.Repository
	Config *config.Config
	ch     chan models.DeleteURLUser
}

func New(store storage.Repository, cfg *config.Config) *Service {
	s := &Service{
		Store:  store,
		Config: cfg,
		ch:     make(chan models.DeleteURLUser),
	}

	go s.Worker(context.TODO())
	return s
}

func (s *Service) CreateURL(ctx context.Context, url string) (string, error) {
	return s.Store.AddURL(ctx, url)
}

func (s *Service) GetURL(ctx context.Context, id string) (string, error) {

	return s.Store.GetURL(ctx, id)

}

func (s *Service) JSONHandle(ctx context.Context, url string) (string, error) {

	logger.Log.Info("start handle JSON")
	url = strings.TrimSpace(url)
	return s.Store.AddURL(ctx, url)
}

func (s *Service) JSONHandleBatch(ctx context.Context, data []models.RequestBatch) ([]models.ResponseBatch, error) {

	logger.Log.Debug(" service JSONHandleBatch")
	return s.Store.AddURLBatch(ctx, data)
}

func (s *Service) PingDB(ctx context.Context) error {
	logger.Log.Debug("PingDB handler start")
	return storepg.Ping(s.Config.DBPath)
}

func (s *Service) GetUserURLs(ctx context.Context) ([]models.UserURL, error) {
	logger.Log.Debug("GetUserURLs handler start")

	return s.Store.GetUserURLs(ctx)
}

func (s *Service) DeleteUserURLs(ctx context.Context, data []string) error {

	logger.Log.Debug("DeleteUserURLs service start")
	userID, ok := ctx.Value(auth.NameCookie).(string)
	if ok {
		go func() {
			for _, v := range data {
				var dw = models.DeleteURLUser{
					Data:   v,
					USerID: userID,
				}
				s.ch <- dw
			}
		}()
	}
	return nil
}

func (s *Service) Worker(ctx context.Context) {
	var count = 500
	data := []models.DeleteURLUser{}

	timer := time.NewTicker(1 * time.Second)
	for {
		select {
		case v := <-s.ch:
			data = append(data, v)
			if len(data) >= count {
				err := s.Store.DeleteUserURLs(ctx, data)
				if err == nil {
					data = []models.DeleteURLUser{}
				}
			}
		case <-timer.C:
			if len(data) > 1 {
				err := s.Store.DeleteUserURLs(ctx, data)
				if err == nil {
					data = []models.DeleteURLUser{}
				}
			}
		default:
			time.Sleep(1 * time.Second)
		}

	}

}
