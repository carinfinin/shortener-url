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

type Store struct {
	store    map[string]models.AuthLine
	mu       sync.Mutex
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
		mu:       sync.Mutex{},
		path:     cfg.FilePath,
		producer: producer,
		URL:      cfg.URL,
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

func (s *Store) AddURL(ctx context.Context, url string) (string, error) {
	s.mu.Lock()

	defer s.mu.Unlock()
	xmlID := s.generateAndExistXMLID(storage.LengthXMLID)

	userID, ok := ctx.Value(auth.NameCookie).(string)
	if !ok {
		return "", auth.ErrorUserNotFound
	}
	line := models.AuthLine{ShortURL: xmlID, OriginalURL: url, UserID: userID}

	//err := s.producer.WriteLine(&line)
	//if err != nil {
	//	return "", err
	//}
	s.store[xmlID] = line

	return xmlID, nil
}

func (s *Store) GetURL(ctx context.Context, xmlID string) (string, error) {
	v, ok := s.store[xmlID]
	if !ok {
		return "", fmt.Errorf("key not found")
	}
	if v.IsDeleted {
		logger.Log.Debug("deleted url")
		return "", storage.ErrDeleteURL
	}

	return v.OriginalURL, nil
}

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

		//err := s.producer.WriteLine(&line)
		//if err != nil {
		//	return nil, err
		//}
		s.store[v.ID] = line

		var tmp = models.ResponseBatch{
			ID:       v.ID,
			ShortURL: s.URL + "/" + v.ID,
		}
		result = append(result, tmp)

	}

	return result, nil
}
func (s *Store) GetUserURLs(ctx context.Context) ([]models.UserURL, error) {
	result := []models.UserURL{}
	userID, ok := ctx.Value(auth.NameCookie).(string)
	if !ok {
		return nil, auth.ErrorUserNotFound
	}
	for _, v := range s.store {
		if v.UserID == userID {
			tmp := models.UserURL{
				ShortURL:    s.URL + "/" + v.ShortURL,
				OriginalURL: v.OriginalURL,
			}
			result = append(result, tmp)
		}
	}

	return result, nil
}

func (s *Store) DeleteUserURLs(ctx context.Context, data []models.DeleteURLUser) error {

	for _, v := range data {

		if line, ok := s.store[v.Data]; ok && line.UserID == v.USerID {
			line.IsDeleted = true
			s.store[v.Data] = line
		}
	}
	return nil
}
