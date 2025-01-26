package store

import (
	"fmt"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	"math/rand"
	"sync"
	"time"
)

type Store struct {
	store map[string]string
	mu    sync.Mutex
	path  string
}

type Line struct {
	URL string `json:"url"`
	ID  string `json:"id"`
}

func New(path string) (*Store, error) {
	logger.Log.Info("stare starting")

	consumer, err := NewConsumer(path)
	if err != nil {
		return nil, err
	}
	defer consumer.Close()

	data, err := consumer.ReadAll()
	if err != nil {
		return nil, err
	}

	logger.Log.Info("stare started")

	return &Store{
		store: data,
		mu:    sync.Mutex{},
		path:  path,
	}, nil
}

const lengthXMLID int64 = 10

func (s *Store) generateAndExistXMLID(length int64) string {
	xmlID := generateXMLID(length)
	if _, ok := s.store[xmlID]; ok {
		return s.generateAndExistXMLID(length + 1)
	} else {
		return xmlID
	}
}

func (s *Store) AddURL(url string) (string, error) {
	s.mu.Lock()

	producer, err := NewProducer(s.path)
	if err != nil {
		return "", err
	}
	defer producer.Close()
	xmlID := s.generateAndExistXMLID(lengthXMLID)
	line := Line{ID: xmlID, URL: url}

	err = producer.WriteLine(&line)
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

func generateXMLID(l int64) string {

	rand.New(rand.NewSource(time.Now().UnixNano()))

	letters := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

	b := make([]byte, l)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
