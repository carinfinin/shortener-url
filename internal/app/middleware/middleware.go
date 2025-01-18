package middleware

import (
	"github.com/carinfinin/shortener-url/internal/app/logger"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
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

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, req)
		log.Printf("%s %s %s", req.Method, req.RequestURI, time.Since(start))
	})
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
