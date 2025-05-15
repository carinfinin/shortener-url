package middleware

import (
	"bytes"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestLogger(t *testing.T) {
	logOutput := bytes.NewBuffer(nil)
	logger.Log.Out = logOutput
	logger.Log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	middleware := RequestLogger(handler)
	middleware.ServeHTTP(rr, req)

	logStr := logOutput.String()
	assert.Contains(t, logStr, "url=/test")
	assert.Contains(t, logStr, "method=GET")
	assert.Contains(t, logStr, "got incoming HTTP request")
}

func TestLoggingResponseWriter(t *testing.T) {
	rr := httptest.NewRecorder()
	resData := &responseData{}

	lw := loggingResponseWriter{
		ResponseWriter: rr,
		responseData:   resData,
	}

	lw.WriteHeader(http.StatusTeapot)
	assert.Equal(t, http.StatusTeapot, resData.status)
	assert.Equal(t, http.StatusTeapot, rr.Code)

	size, err := lw.Write([]byte("test"))
	assert.NoError(t, err)
	assert.Equal(t, 4, size)
	assert.Equal(t, 4, resData.size)
}
