package service

import (
	"context"
	"github.com/carinfinin/shortener-url/internal/app/auth"
	"github.com/carinfinin/shortener-url/internal/app/config"
	"github.com/carinfinin/shortener-url/internal/app/models"
	"github.com/carinfinin/shortener-url/internal/app/storage"
	"github.com/carinfinin/shortener-url/internal/app/storage/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestService_CreateURLPositive(t *testing.T) {

	r := MockRepository{}
	r.On("AddURL", mock.Anything, "test.ru").Return("short", nil)

	cfg := config.Config{}
	s := New(&r, &cfg)

	url, err := s.CreateURL(context.Background(), "test.ru")

	assert.NoError(t, err, err)
	assert.Equal(t, url, "short")

}
func TestService_CreateURLNegative(t *testing.T) {

	r := MockRepository{}
	r.On("AddURL", mock.Anything, "test.ru").Return("", storage.ErrDeleteURL)

	cfg := config.Config{}
	s := New(&r, &cfg)

	url, err := s.CreateURL(context.Background(), "test.ru")

	assert.Error(t, err, &storage.ErrDeleteURL)
	assert.Empty(t, url)

}
func TestService_CreateURLNegative2(t *testing.T) {

	r := MockRepository{}
	r.On("AddURL", mock.Anything, "test.ru").Return("short", storage.ErrDouble)

	cfg := config.Config{}
	s := New(&r, &cfg)

	id, err := s.CreateURL(context.Background(), "test.ru")

	assert.Error(t, err, &storage.ErrDeleteURL)
	assert.Equal(t, id, "short")

}

func TestService_GetURL(t *testing.T) {
	t.Run("positive", func(t *testing.T) {

		r := MockRepository{}
		r.On("GetURL", mock.Anything, "1245").Return("test.ru", nil)
		cfg := config.Config{}
		s := New(&r, &cfg)
		url, err := s.GetURL(context.Background(), "1245")
		assert.NoError(t, err, err)
		assert.Equal(t, url, "test.ru")
	})

	t.Run("negative", func(t *testing.T) {

		r := MockRepository{}
		r.On("GetURL", mock.Anything, "1245").Return("", storage.ErrDeleteURL)
		cfg := config.Config{}
		s := New(&r, &cfg)
		url, err := s.GetURL(context.Background(), "1245")
		assert.Error(t, err, &storage.ErrDeleteURL)
		assert.Empty(t, url)
	})
}

func TestJSONHandleBatch(t *testing.T) {
	cfg := config.New()
	r, err := store.New(cfg)
	require.NoError(t, err, err)

	token := auth.GenerateToken()
	ctx := context.WithValue(context.Background(), auth.NameCookie, token)

	s := New(r, cfg)
	data := make([]models.RequestBatch, 0)
	data = append(data, models.RequestBatch{
		ID:      "123",
		LongURL: "practicum.ru",
	})

	result, err := s.JSONHandleBatch(ctx, data)
	assert.NoError(t, err, err)

	assert.Equal(t, result, []models.ResponseBatch{
		{
			ID:       "123",
			ShortURL: "http://localhost:8080/123",
		},
	})

}
