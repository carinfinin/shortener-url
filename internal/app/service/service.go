package service

import (
	"context"
	"github.com/carinfinin/shortener-url/internal/app/auth"
	"github.com/carinfinin/shortener-url/internal/app/config"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	"github.com/carinfinin/shortener-url/internal/app/models"
	"github.com/carinfinin/shortener-url/internal/app/storage/storepg"
	"strings"
	"time"
)

// Repository интерфейс базы данных.
//

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=Repository --filename=repositorymock_test.go --inpackage
type Repository interface {
	AddURL(ctx context.Context, url string) (string, error)
	GetURL(ctx context.Context, xmlID string) (string, error)
	AddURLBatch(ctx context.Context, data []models.RequestBatch) ([]models.ResponseBatch, error)
	GetUserURLs(ctx context.Context) ([]models.UserURL, error)
	DeleteUserURLs(ctx context.Context, data []models.DeleteURLUser) error
	Close() error
}

// Service реализует интерфейс IService.
type Service struct {
	store  Repository
	Config *config.Config
	ch     chan models.DeleteURLUser
	close  context.CancelFunc
}

// Close закрытие сервиса
func (s *Service) Close() {
	s.close()
}

// New конструктор для Service.
func New(store Repository, cfg *config.Config) *Service {
	s := &Service{
		store:  store,
		Config: cfg,
		ch:     make(chan models.DeleteURLUser),
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.close = cancel

	go s.worker(ctx)
	return s
}

// CreateURL создаёт урл.
func (s *Service) CreateURL(ctx context.Context, url string) (string, error) {
	url = strings.TrimSpace(url)
	return s.store.AddURL(ctx, url)
}

// GetURL получает урл.
func (s *Service) GetURL(ctx context.Context, id string) (string, error) {

	return s.store.GetURL(ctx, id)

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
		case <-ctx.Done():
			err := s.store.DeleteUserURLs(ctx, data)
			if err != nil {
				logger.Log.Error("worker error", err)
			}
			return
		}
	}
}
