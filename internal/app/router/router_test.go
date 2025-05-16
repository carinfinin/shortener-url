package router

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/carinfinin/shortener-url/internal/app/auth"
	"github.com/carinfinin/shortener-url/internal/app/config"
	"github.com/carinfinin/shortener-url/internal/app/models"
	"github.com/carinfinin/shortener-url/internal/app/service"
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

			token := auth.GenerateToken()
			ctx := context.WithValue(request.Context(), auth.NameCookie, token)
			newReq := request.WithContext(ctx)

			cfg := config.Config{URL: test.url}
			s, err := store.New(&cfg)
			require.NoError(t, err)
			service := service.New(s, &cfg)
			r := ConfigureRouter(service, &cfg)
			w := httptest.NewRecorder()

			r.CreateURL(w, newReq)

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
			cfg := config.Config{URL: test.url}

			s, err := store.New(&cfg)
			require.NoError(t, err)

			token := auth.GenerateToken()
			ctx := context.WithValue(context.Background(), auth.NameCookie, token)

			xmlID, err := s.AddURL(ctx, test.data)
			require.NoError(t, err)

			request := httptest.NewRequest(http.MethodGet, test.request+xmlID, nil)
			service := service.New(s, &cfg)
			r := ConfigureRouter(service, &cfg)

			w := httptest.NewRecorder()

			r.GetURL(w, request)

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

func TestJSONHandle(t *testing.T) {
	tests := []struct {
		name       string
		data       string
		request    string
		url        string
		statusCode int
	}{
		{
			name:       "simple test #1",
			data:       "{\n  \"url\": \"https://practicum.yandex.ru\"\n}",
			request:    "/api/shorten",
			url:        "http://localhost:8080",
			statusCode: 201,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			buf := bytes.NewBuffer([]byte(test.data))
			var req models.Request

			err := json.Unmarshal([]byte(test.data), &req)
			assert.NoError(t, err)

			request := httptest.NewRequest(http.MethodPost, test.request, buf)
			w := httptest.NewRecorder()
			cfg := config.Config{URL: test.url}

			s, err := store.New(&cfg)
			require.NoError(t, err)

			token := auth.GenerateToken()
			ctx := context.WithValue(request.Context(), auth.NameCookie, token)
			newReq := request.WithContext(ctx)
			service := service.New(s, &cfg)
			r := ConfigureRouter(service, &cfg)
			r.JSONHandle(w, newReq)
			result := w.Result()
			assert.Equal(t, test.statusCode, result.StatusCode)

			var res models.Response
			decoder := json.NewDecoder(result.Body)

			err = decoder.Decode(&res)

			assert.NoError(t, err)
			assert.NotNil(t, res.Result)
			err = result.Body.Close()
			assert.NoError(t, err)
		})
	}
}

func TestRouter_User(t *testing.T) {
	//JSONHandleBatch
	//GetUserURLs
	//DeleteUserURLs
	var url = "http://localhost:8080"
	js := []byte(`[{"correlation_id": "123", "original_url": "practicum.ru"}]`)

	buf := bytes.NewBuffer(js)

	request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", buf)
	w := httptest.NewRecorder()
	cfg := config.Config{URL: url}

	s, err := store.New(&cfg)
	require.NoError(t, err)

	token := auth.GenerateToken()
	ctx := context.WithValue(request.Context(), auth.NameCookie, token)
	newReq := request.WithContext(ctx)
	ss := service.New(s, &cfg)
	r := ConfigureRouter(ss, &cfg)
	r.JSONHandleBatch(w, newReq)
	result := w.Result()

	fmt.Println(result.StatusCode)
	assert.Equal(t, 201, result.StatusCode)
	err = result.Body.Close()
	assert.NoError(t, err)

	// get urls user
	request = httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	newReq = request.WithContext(ctx)
	w = httptest.NewRecorder()
	r.GetUserURLs(w, newReq)
	result = w.Result()
	fmt.Println(result.StatusCode)
	assert.Equal(t, 200, result.StatusCode)
	err = result.Body.Close()
	assert.NoError(t, err)

	// delete
	js = []byte(`["123"]`)
	buf = bytes.NewBuffer(js)
	request = httptest.NewRequest(http.MethodDelete, "/api/user/urls", buf)
	newReq = request.WithContext(ctx)
	w = httptest.NewRecorder()
	r.DeleteUserURLs(w, newReq)
	result = w.Result()
	fmt.Println(result.StatusCode)
	assert.Equal(t, 202, result.StatusCode)

	err = result.Body.Close()
	assert.NoError(t, err)
}

func BenchmarkRouter_CreateURL(b *testing.B) {

	body := strings.NewReader("https://yandex.ru")
	request := httptest.NewRequest(http.MethodPost, "/", body)

	token := auth.GenerateToken()
	ctx := context.WithValue(request.Context(), auth.NameCookie, token)
	newReq := request.WithContext(ctx)

	cfg := config.Config{URL: "http://localhost:8080"}
	s, _ := store.New(&cfg)

	service := service.New(s, &cfg)
	r := ConfigureRouter(service, &cfg)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		body.Seek(0, 0)
		r.CreateURL(w, newReq)
	}

}

func BenchmarkRouter_GetURL(b *testing.B) {

	cfg := config.Config{URL: "http://localhost:8080"}

	s, _ := store.New(&cfg)
	token := auth.GenerateToken()
	ctx := context.WithValue(context.Background(), auth.NameCookie, token)

	xmlID, _ := s.AddURL(ctx, "https://www.google.com")

	request := httptest.NewRequest(http.MethodGet, "/"+xmlID, nil)
	service := service.New(s, &cfg)
	r := ConfigureRouter(service, &cfg)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r.GetURL(w, request)
	}

}
