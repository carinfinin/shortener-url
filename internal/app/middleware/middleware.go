package middleware

import (
	"fmt"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(data []byte) (int, error) {
	size, err := r.ResponseWriter.Write(data)
	r.responseData.size = size
	return size, err
}
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		next.ServeHTTP(writer, request)

		logger.Log.WithFields(logrus.Fields{
			"url":    request.RequestURI,
			"method": request.Method,
			"time":   time.Since(start),
		}).Info("got incoming HTTP request")

	})
}

func ResponseLogger(next http.Handler) http.Handler {
	fn := func(writer http.ResponseWriter, request *http.Request) {

		resData := responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: writer,
			responseData:   &resData,
		}

		next.ServeHTTP(&lw, request)

		logger.Log.WithFields(logrus.Fields{
			"status": resData.status,
			"size":   resData.size,
		}).Info("got incoming HTTP response")

	}
	return http.HandlerFunc(fn)
}

func AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		fmt.Println(request.Cookies())
		nameCookie := "token"
		valueCookie := generateID()

		c, err := request.Cookie(nameCookie)
		if err != nil || !decodeCookie(c) {
			c := http.Cookie{
				Name:   nameCookie,
				Value:  valueCookie,
				MaxAge: 300,
			}
			http.SetCookie(writer, &c)
		}
		next.ServeHTTP(writer, request)
	})
}

func generateID() string {
	return strconv.Itoa(222222222)
}
func decodeCookie(cookie *http.Cookie) bool {
	if cookie.Value != "" {
		return true
	}
	return false

}
