package store

import (
	"fmt"
	"math/rand"
)

type Store struct {
	store map[string]string
}

func New() *Store {
	return &Store{
		store: make(map[string]string),
	}
}

const lengthXMLID int64 = 10

func (s *Store) AddURL(url string) string {
	xmlID := generateXMLID(lengthXMLID)
	s.store[xmlID] = url

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
	startChar := "a"
	temp := ""
	var i int64 = 1
	for {
		myRand := random(0, 26)
		newChar := string(startChar[0] + byte(myRand))
		temp = temp + newChar
		if i == l {
			break
		}
		i++
	}
	return temp
}
func random(min, max int) int {
	return rand.Intn(max-min) + min
}
