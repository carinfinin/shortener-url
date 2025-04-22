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

// IService сервисный слой реализует бизнес логику.
type IService interface {
	CreateURL(ctx context.Context, url string) (string, error)
	GetURL(ctx context.Context, id string) (string, error)
	JSONHandle(ctx context.Context, url string) (string, error)
	JSONHandleBatch(ctx context.Context, data []models.RequestBatch) ([]models.ResponseBatch, error)
	PingDB(ctx context.Context) error
	GetUserURLs(ctx context.Context) ([]models.UserURL, error)
	DeleteUserURLs(ctx context.Context, data []string) error
}

// Service реализует интерфейс IService.
type Service struct {
	store  storage.Repository
	Config *config.Config
	ch     chan models.DeleteURLUser
}

// New конструктор для Service.
func New(store storage.Repository, cfg *config.Config) *Service {
	s := &Service{
		store:  store,
		Config: cfg,
		ch:     make(chan models.DeleteURLUser),
	}

	go s.worker(context.TODO())
	return s
}

// CreateURL создаёт урл.
func (s *Service) CreateURL(ctx context.Context, url string) (string, error) {
	return s.store.AddURL(ctx, url)
}

// GetURL получает урл.
func (s *Service) GetURL(ctx context.Context, id string) (string, error) {

	return s.store.GetURL(ctx, id)

}

// JSONHandle создаёт урл принимая json.
func (s *Service) JSONHandle(ctx context.Context, url string) (string, error) {

	logger.Log.Info("start handle JSON")
	url = strings.TrimSpace(url)
	return s.store.AddURL(ctx, url)
}

// JSONHandleBatch создаёт пачку урлов принимая json.
func (s *Service) JSONHandleBatch(ctx context.Context, data []models.RequestBatch) ([]models.ResponseBatch, error) {

	logger.Log.Debug(" service JSONHandleBatch")
	return s.store.AddURLBatch(ctx, data)
}

// PingDB проверяет доступность бд.
func (s *Service) PingDB(ctx context.Context) error {
	logger.Log.Debug("PingDB handler start")
	return storepg.Ping(s.Config.DBPath)
}

// GetUserURLs получпет урлы пользователя.
func (s *Service) GetUserURLs(ctx context.Context) ([]models.UserURL, error) {
	logger.Log.Debug("GetUserURLs handler start")

	return s.store.GetUserURLs(ctx)
}

// DeleteUserURLs удаляет урлы пользователя.
func (s *Service) DeleteUserURLs(ctx context.Context, data []string) error {

	logger.Log.Debug("DeleteUserURLs service start")
	userID, ok := ctx.Value(auth.NameCookie).(string)
	if !ok {
		logger.Log.Debug("DeleteUserURLs ErrorUserNotFound")
		return auth.ErrorUserNotFound
	}
	go func() {
		for _, v := range data {
			var dw = models.DeleteURLUser{
				Data:   v,
				USerID: userID,
			}
			s.ch <- dw
		}
	}()
	return nil
}

func (s *Service) worker(ctx context.Context) {
	var count = 100
	data := []models.DeleteURLUser{}

	timer := time.NewTicker(10 * time.Second)
	for {
		select {
		case v := <-s.ch:
			data = append(data, v)
			if len(data) >= count {
				err := s.store.DeleteUserURLs(ctx, data)
				if err != nil {
					logger.Log.Error("worker error", err)
				}
				data = data[:0]
				timer.Reset(10 * time.Second)
			}
		case <-timer.C:
			if len(data) > 0 {
				err := s.store.DeleteUserURLs(ctx, data)
				if err != nil {
					logger.Log.Error("worker error", err)
				}
				data = data[:0]
			}
		}
	}
}
