package storage

import (
	"context"
	"errors"
	"github.com/carinfinin/shortener-url/internal/app/models"
	"math/rand"
	"time"
)

// ErrDouble ошибка возвращаемая при дубликате.
var ErrDouble = errors.New("duplicate url")

// ErrDeleteURL ошибка возвращаемая при удалённом елементе.
var ErrDeleteURL = errors.New("deleted url")

// Repository интерфейс базы данных.
//
//go:generate mockgen -source=storage.go -destination=mocks/storage_mock.go -package=mocks
type Repository interface {
	AddURL(ctx context.Context, url string) (string, error)
	GetURL(ctx context.Context, xmlID string) (string, error)
	AddURLBatch(ctx context.Context, data []models.RequestBatch) ([]models.ResponseBatch, error)
	GetUserURLs(ctx context.Context) ([]models.UserURL, error)
	DeleteUserURLs(ctx context.Context, data []models.DeleteURLUser) error
	Close() error
}

// ProducerInterface интерфейс для storefile.
type ProducerInterface interface {
	WriteLine(line *models.AuthLine) error
	Close(data map[string]models.AuthLine) error
}

// LengthXMLID константа длинна символов генерации короткого урл.
const LengthXMLID int64 = 10

// GenerateXMLID генерирует короткий урл
func GenerateXMLID(l int64) string {

	rand.New(rand.NewSource(time.Now().UnixNano()))
	letters := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

	b := make([]byte, l)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
