package router

import (
	"fmt"
	"github.com/carinfinin/shortener-url/internal/app/storage/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateURL(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name    string
		request string
		want    want
		url     string
	}{
		{
			name: "simple test #1",
			want: want{
				contentType: "text/plain",
				statusCode:  201,
			},
			url:     "http://localhost:8080",
			request: "/",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodPost, test.request, strings.NewReader("https://yandex.ru"))

			s := store.New()
			r := ConfigureRouter(s, test.url)
			w := httptest.NewRecorder()

			h := CreateURL(*r)
			h(w, request)

			result := w.Result()

			assert.Equal(t, test.want.statusCode, result.StatusCode)
			assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))

			newURL, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.NotNil(t, newURL)

			fmt.Println("test create URL")
		})
	}
}

func TestGetURL(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name    string
		data    string
		request string
		want    want
		url     string
	}{
		{
			name: "simple test #1",
			data: "https://www.google.com",
			want: want{
				statusCode: 307,
			},
			url:     "http://localhost:8080",
			request: "/",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			s := store.New()

			xmlID := s.AddURL(test.data)

			request := httptest.NewRequest(http.MethodGet, test.request+xmlID, nil)

			r := ConfigureRouter(s, test.url)
			w := httptest.NewRecorder()

			h := GetURL(*r)
			h(w, request)

			result := w.Result()
			fmt.Println(result)

			assert.Equal(t, test.want.statusCode, result.StatusCode)

			assert.Equal(t, test.data, result.Header.Get("Location"))

			newURL, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			fmt.Println(string(newURL))
			assert.NotNil(t, newURL)

			fmt.Println("Location", result.Header.Get("Location"))
		})
	}
}
