package storefile

import (
	"fmt"
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

func New(path string) (*Store, error) {

	data, err := readAllinMemory(path)
	if err != nil {
		logger.Log.Error("error in readAllinMemory", err)
		return nil, err
	}

	producer, err := NewProducer(path)
	if err != nil {
		return nil, err
	}

	return &Store{
		store:    data,
		mu:       sync.Mutex{},
		path:     path,
		producer: producer,
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

	xmlID := s.generateAndExistXMLID(storage.LengthXMLID)
	line := models.Line{ID: xmlID, URL: url}

	err := s.producer.WriteLine(&line)
	if err != nil {
		return "", err
	}
	s.store[xmlID] = url

	s.mu.Unlock()
	return xmlID, nil
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
