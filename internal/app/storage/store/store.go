package store

import (
	"fmt"
	"github.com/carinfinin/shortener-url/internal/app/config"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	"github.com/carinfinin/shortener-url/internal/app/models"
	"github.com/carinfinin/shortener-url/internal/app/storage"
	"sync"
)

type Store struct {
	store map[string]string
	mu    sync.Mutex
	URL   string
}

func New(cfg *config.Config) (*Store, error) {
	logger.Log.Info("store starting")

	return &Store{
		store: map[string]string{},
		mu:    sync.Mutex{},
		URL:   cfg.URL,
	}, nil
}

func (s *Store) generateAndExistXMLID(length int64) string {
	xmlID := storage.GenerateXMLID(length)
	if _, ok := s.store[xmlID]; ok {
		return s.generateAndExistXMLID(length + 1)
	} else {
		return xmlID
	}
}

func (s *Store) AddURL(url string) (string, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	ID := s.generateAndExistXMLID(storage.LengthXMLID)
	for i, v := range s.store {
		if v == url {
			logger.Log.Error(" AddURL error : дублирование URL")

			return i, storage.ErrDouble
		}
	}
	s.store[ID] = url
	return ID, nil
}

func (s *Store) GetURL(xmlID string) (string, error) {
	v, ok := s.store[xmlID]
	if !ok {
		return "", fmt.Errorf("ключ не найден")
	}
	return v, nil
}
func (s *Store) Close() error {
	logger.Log.Info("closed store")
	s.store = nil
	return nil
}

func (s *Store) AddURLBatch(data []models.RequestBatch) ([]models.ResponseBatch, error) {

	var result []models.ResponseBatch
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, v := range data {
		if _, ok := s.store[v.ID]; ok {
			return nil, fmt.Errorf("incorrect id in data request")
		}
	}
	for _, v := range data {
		s.store[v.ID] = v.LongURL

		var tmp = models.ResponseBatch{
			ID:       v.ID,
			ShortURL: s.URL + "/" + v.ID,
		}
		result = append(result, tmp)

	}

	return result, nil
}

func (s *Store) Ping() error {
	return nil
}
