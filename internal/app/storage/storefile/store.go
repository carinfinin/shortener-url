package storefile

import (
	"context"
	"fmt"
	"github.com/carinfinin/shortener-url/internal/app/auth"
	"github.com/carinfinin/shortener-url/internal/app/config"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	"github.com/carinfinin/shortener-url/internal/app/models"
	"github.com/carinfinin/shortener-url/internal/app/storage"
	"sync"
)

// Store хранилище в файлле реализует интерфейс Repository.
type Store struct {
	store    map[string]models.AuthLine
	mu       sync.RWMutex
	path     string
	producer storage.ProducerInterface
	URL      string
}

func readAllinMemory(path string) (map[string]models.AuthLine, error) {
	logger.Log.Info("coll function readAllinMemory")
	consumer, err := NewConsumer(path)
	if err != nil {
		return nil, err
	}
	defer consumer.Close()

	return consumer.ReadAll()
}

// New конструктор для  Store.
func New(cfg *config.Config) (*Store, error) {
	logger.Log.Info("start store in file")

	data, err := readAllinMemory(cfg.FilePath)
	if err != nil {
		logger.Log.Error("error in readAllinMemory", err)
		return nil, err
	}

	producer, err := NewProducer(cfg.FilePath)
	if err != nil {
		return nil, err
	}

	return &Store{
		store:    data,
		mu:       sync.RWMutex{},
		path:     cfg.FilePath,
		producer: producer,
		URL:      cfg.URL,
	}, nil
}

func (s *Store) generateAndExistXMLID(length int64) string {
	ID := storage.GenerateXMLID(length)
	if _, ok := s.store[ID]; ok {
		return s.generateAndExistXMLID(length + 1)
	} else {
		return ID
	}
}

// AddURL записывает в хранилище урл.
func (s *Store) AddURL(ctx context.Context, url string) (string, error) {

	ID := s.generateAndExistXMLID(storage.LengthXMLID)

	userID, ok := ctx.Value(auth.NameCookie).(string)
	if !ok {
		return "", auth.ErrorUserNotFound
	}

	for _, v := range s.store {
		if v.OriginalURL == url {
			return v.ShortURL, storage.ErrDouble
		}
	}

	line := models.AuthLine{ShortURL: ID, OriginalURL: url, UserID: userID}

	s.mu.Lock()
	s.store[ID] = line
	s.mu.Unlock()

	return ID, nil
}

// GetURL получает из хранилища урл.
func (s *Store) GetURL(ctx context.Context, ID string) (string, error) {
	s.mu.RLock()

	v, ok := s.store[ID]
	if !ok {
		return "", fmt.Errorf("key not found")
	}
	s.mu.RUnlock()

	if v.IsDeleted {
		logger.Log.Debug("deleted url")
		return "", storage.ErrDeleteURL
	}

	return v.OriginalURL, nil
}

// Close закрывает хранилище.
func (s *Store) Close() error {
	err := s.producer.Close(s.store)
	if err != nil {
		logger.Log.Error("error closed store", err)
		return err
	}
	logger.Log.Info("closed store")
	s.store = nil
	return nil
}

// AddURLBatch добавляет добавляет пачку урлов.
func (s *Store) AddURLBatch(ctx context.Context, data []models.RequestBatch) ([]models.ResponseBatch, error) {
	var result []models.ResponseBatch
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, v := range data {
		if _, ok := s.store[v.ID]; ok {
			return nil, fmt.Errorf("incorrect id in data request")
		}
	}
	for _, v := range data {

		userID, ok := ctx.Value(auth.NameCookie).(string)
		if !ok {
			return nil, auth.ErrorUserNotFound
		}
		line := models.AuthLine{ShortURL: v.ID, OriginalURL: v.LongURL, UserID: userID}

		s.store[v.ID] = line

		var tmp = models.ResponseBatch{
			ID:       v.ID,
			ShortURL: s.URL + "/" + v.ID,
		}
		result = append(result, tmp)

	}

	return result, nil
}

// GetUserURLs получает урлы пользователя.
func (s *Store) GetUserURLs(ctx context.Context) ([]models.UserURL, error) {
	result := []models.UserURL{}
	userID, ok := ctx.Value(auth.NameCookie).(string)
	if !ok {
		return nil, auth.ErrorUserNotFound
	}

	s.mu.RLock()
	for _, v := range s.store {
		if v.UserID == userID {
			tmp := models.UserURL{
				ShortURL:    s.URL + "/" + v.ShortURL,
				OriginalURL: v.OriginalURL,
			}
			result = append(result, tmp)
		}
	}
	s.mu.RUnlock()

	return result, nil
}

// DeleteUserURLs удаляет  урлы пользователя.
func (s *Store) DeleteUserURLs(ctx context.Context, data []models.DeleteURLUser) error {

	for _, v := range data {

		if line, ok := s.store[v.Data]; ok && line.UserID == v.USerID {
			line.IsDeleted = true
			s.store[v.Data] = line
		}
	}
	return nil
}
