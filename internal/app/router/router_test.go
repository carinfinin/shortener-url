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
		Response    []byte
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name: "simple test #1",
			want: want{
				contentType: "text/plain",
				statusCode:  201,
				Response:    []byte("ee"),
			},
			request: "/",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodPost, test.request, strings.NewReader("https://yandex.ru"))

			s := store.New()
			r := ConfigureRouter(s)
			w := httptest.NewRecorder()

			h := CreateURL(*r)
			h(w, request)

			result := w.Result()

			assert.Equal(t, test.want.statusCode, result.StatusCode)
			assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))

			newUrl, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.NotNil(t, newUrl)

			fmt.Println("test create URL")
		})
	}
}

func TestGetURL(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		Response    []byte
	}
	tests := []struct {
		name    string
		data    string
		request string
		want    want
	}{
		{
			name: "simple test #1",
			data: "https://google.ru",
			want: want{
				statusCode: 307,
				Response:   []byte("ee"),
			},
			request: "/",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			s := store.New()

			xmlID := s.AddURL(test.data)

			request := httptest.NewRequest(http.MethodGet, test.request+xmlID, nil)

			r := ConfigureRouter(s)
			w := httptest.NewRecorder()

			h := GetURL(*r)
			h(w, request)

			result := w.Result()

			assert.Equal(t, test.want.statusCode, result.StatusCode)

			assert.Equal(t, test.data, result.Header.Get("Location"))

			newUrl, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			fmt.Println(string(newUrl))
			assert.NotNil(t, newUrl)

			fmt.Println("Location", result.Header.Get("Location"))
		})
	}
}
