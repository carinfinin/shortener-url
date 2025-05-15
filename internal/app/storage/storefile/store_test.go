package storefile

import (
	"context"
	"fmt"
	"github.com/carinfinin/shortener-url/internal/app/auth"
	"github.com/carinfinin/shortener-url/internal/app/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestStore_AddURL(t *testing.T) {

	tests := []struct {
		name   string
		err    bool
		url    string
		userID string
	}{
		{
			name:   "positive",
			err:    false,
			url:    "practicum.ru",
			userID: "1",
		},
		{
			name:   "fail",
			err:    true,
			url:    "practicum.ru",
			userID: "",
		},
	}

	cfg := config.New()
	cfg.FilePath = "test.json"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s, err := New(cfg)
			assert.NoError(t, err, err)

			ctx := context.WithValue(context.Background(), auth.NameCookie, tt.userID)

			short, err := s.AddURL(ctx, tt.url)

			fmt.Println(err)
			if tt.err {
				assert.ErrorIs(t, err, auth.ErrorUserNotFound)
			} else {
				assert.NoError(t, err, err)
			}

			url, err := s.GetURL(ctx, short)
			if tt.err {
				assert.EqualError(t, err, "key not found")

			} else {
				assert.NoError(t, err, err)
			}

			if !tt.err {
				assert.Equal(t, url, tt.url)
			}
			err = s.Close()
			assert.NoError(t, err, err)
		})
	}

	os.Remove("test.json")

}
