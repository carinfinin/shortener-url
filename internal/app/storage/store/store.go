package store

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Store struct {
	store map[string]string
	mu    sync.Mutex
}

func New() *Store {
	return &Store{
		store: make(map[string]string),
		mu:    sync.Mutex{},
	}
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

func (s *Store) AddURL(url string) string {
	s.mu.Lock()
	xmlID := s.generateAndExistXMLID(lengthXMLID)
	s.store[xmlID] = url
	s.mu.Unlock()
	return xmlID
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
