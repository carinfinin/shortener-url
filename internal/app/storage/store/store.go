package store

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

// Store in memory хранилище реализует интерфейс Repository.
type Store struct {
	store map[string]models.AuthLine
	mu    sync.RWMutex
	URL   string
}

// New конструктор для  Store.
func New(cfg *config.Config) (*Store, error) {
	logger.Log.Info("store starting")

	return &Store{
		store: map[string]models.AuthLine{},
		mu:    sync.RWMutex{},
		URL:   cfg.URL,
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

	for _, v := range s.store {
		if v.OriginalURL == url {
			return v.ShortURL, storage.ErrDouble
		}
	}

	ID := s.generateAndExistXMLID(storage.LengthXMLID)

	userID, ok := ctx.Value(auth.NameCookie).(string)
	if !ok || userID == "" {
		return "", auth.ErrorUserNotFound
	}

	line := models.AuthLine{
		UserID:      userID,
		OriginalURL: url,
		ShortURL:    ID,
	}
	s.mu.Lock()
	s.store[ID] = line
	s.mu.Unlock()

	return ID, nil
}

// GetURL получает из хранилища урл.
func (s *Store) GetURL(ctx context.Context, id string) (string, error) {
	s.mu.RLock()
	v, ok := s.store[id]
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
	userID, ok := ctx.Value(auth.NameCookie).(string)
	if !ok {
		return nil, auth.ErrorUserNotFound
	}

	for _, v := range data {

		line := models.AuthLine{
			UserID:      userID,
			OriginalURL: v.LongURL,
			ShortURL:    v.ID,
		}
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
