package storefile

import (
	"fmt"
	"github.com/carinfinin/shortener-url/internal/app/config"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	"github.com/carinfinin/shortener-url/internal/app/models"
	"github.com/carinfinin/shortener-url/internal/app/storage"
	"sync"
)

type Store struct {
	store    map[string]string
	mu       sync.Mutex
	path     string
	producer storage.ProducerInterface
	URL      string
}

// TODO нужно консьюмер наверное убрать
func readAllinMemory(path string) (map[string]string, error) {
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
	ID := storage.GenerateXMLID(length)
	if _, ok := s.store[ID]; ok {
		return s.generateAndExistXMLID(length + 1)
	} else {
		return ID
	}
}

func (s *Store) AddURL(url string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ID := s.generateAndExistXMLID(storage.LengthXMLID)
	line := models.Line{ID: ID, URL: url}

	err := s.producer.WriteLine(&line)
	if err != nil {
		return "", err
	}
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
	err := s.producer.Close()
	if err != nil {
		logger.Log.Error("error closed store", err)
		return err
	}
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

		line := models.Line{ID: v.ID, URL: v.LongURL}

		err := s.producer.WriteLine(&line)
		if err != nil {
			return nil, err
		}
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
