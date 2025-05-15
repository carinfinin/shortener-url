package router_test

import (
	"encoding/json"
	"fmt"
	"github.com/carinfinin/shortener-url/internal/app/config"
	"github.com/carinfinin/shortener-url/internal/app/models"
	"github.com/carinfinin/shortener-url/internal/app/router"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
)

func ExampleRouter_JSONHandle() {
	cfg := config.Config{URL: "http://localhost:8080"}

	s := &router.MockService{}
	s.On("CreateURL", mock.Anything, "https://example.com").Return("3521", nil)
	r := router.ConfigureRouter(s, &cfg)

	reqBody := `{"url": "https://example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.JSONHandle(w, req)

	fmt.Println(w.Code)
	fmt.Println(w.Header().Get("Content-Type"))

	// Output:
	// 201
	// application/json
}

func ExampleRouter_JSONHandleBatch() {
	cfg := config.Config{URL: "http://localhost:8080"}

	reqBody := `[{"correlation_id": "453","original_url": "https://example.com"}]`
	resBody := `[{"correlation_id": "453","short_url": "http://localhost:8080/453"}]`
	res := []models.ResponseBatch{}
	req := []models.RequestBatch{}
	json.Unmarshal([]byte(resBody), &res)
	json.Unmarshal([]byte(reqBody), &req)
	s := &router.MockService{}
	s.On("JSONHandleBatch", mock.Anything, req).Return(res, nil)
	r := router.ConfigureRouter(s, &cfg)

	request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(reqBody))
	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.JSONHandleBatch(w, request)

	fmt.Println(w.Code)
	fmt.Println(w.Header().Get("Content-Type"))
	fmt.Println(w.Body)

	// Output:
	// 201
	// application/json
	// [{"correlation_id":"453","short_url":"http://localhost:8080/453"}]
}
