package storage

import (
	"github.com/carinfinin/shortener-url/internal/app/models"
	"math/rand"
	"time"
)

type Repositories interface {
	AddURL(url string) (string, error)
	GetURL(xmlID string) (string, error)
	Close() error
}

type ProducerInterface interface {
	WriteLine(line *models.Line) error
	Close() error
}

// to do  обвить метод close

const LengthXMLID int64 = 10

func GenerateXMLID(l int64) string {

	rand.New(rand.NewSource(time.Now().UnixNano()))
	letters := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

	b := make([]byte, l)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
