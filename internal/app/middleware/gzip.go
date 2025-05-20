package middleware

import (
	"github.com/carinfinin/shortener-url/internal/app/compress"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	"net/http"
	"strings"
)

// CompressGzipWriter compresses and write data in gzip format.
func CompressGzipWriter(h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		w := writer

		if strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
			logger.Log.Info("Accept-Encoding == gzip")
			logger.Log.Info(request.Header)

			cw := compress.NewCompressWriter(writer)
			defer cw.Close()
			w = cw
			w.Header().Set("Content-Encoding", "gzip")
		}

		h.ServeHTTP(w, request)
	})
}

// CompressGzipReader reads compressed data.
func CompressGzipReader(h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if strings.Contains(request.Header.Get("Content-Encoding"), "gzip") {
			logger.Log.Info("Content-Encoding == gzip")
			logger.Log.Info(request.Header)

			cr, err := compress.NewCompressReader(request.Body)
			logger.Log.Info(cr)

			if err != nil {
				logger.Log.Info("error newCompressReader", err)
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
			request.Body = cr
			defer cr.Close()
		}

		h.ServeHTTP(writer, request)
	})
}
