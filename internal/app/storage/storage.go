package storage

import (
	"context"
	"errors"
	"github.com/carinfinin/shortener-url/internal/app/models"
	"math/rand"
	"time"
)

var ErrDouble = errors.New("duplicate url")
var ErrDeleteURL = errors.New("deleted url")

type Repository interface {
	AddURL(ctx context.Context, url string) (string, error)
	GetURL(ctx context.Context, xmlID string) (string, error)
	AddURLBatch(ctx context.Context, data []models.RequestBatch) ([]models.ResponseBatch, error)
	GetUserURLs(ctx context.Context) ([]models.UserURL, error)
	DeleteUserURLs(ctx context.Context, data []models.DeleteURLUser) error
	Close() error
}

type ProducerInterface interface {
	WriteLine(line *models.AuthLine) error
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
