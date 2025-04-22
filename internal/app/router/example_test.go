package router_test

import (
	"fmt"
	"github.com/carinfinin/shortener-url/internal/app/config"
	"github.com/carinfinin/shortener-url/internal/app/router"
	"github.com/carinfinin/shortener-url/internal/app/service/mocks"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
)

func ExampleRouter_JSONHandle() {
	cfg := config.Config{URL: "http://localhost:8080"}

	s := &mocks.IService{}
	s.On("JSONHandle", mock.Anything, "https://example.com").Return("3521", nil)
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
