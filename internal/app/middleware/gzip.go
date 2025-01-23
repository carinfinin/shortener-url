package middleware

import (
	"github.com/carinfinin/shortener-url/internal/app/compress"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	"net/http"
	"strings"
)

func CompressGzip(h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		w := writer

		if strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
			logger.Log.Info("Accept-Encoding == gzip")

			cw := compress.NewCompressWriter(writer)
			w = cw
			defer cw.Close()
		}

		if strings.Contains(request.Header.Get("Content-Encoding"), "gzip") {
			logger.Log.Info("Content-Encoding == gzip")

			cr, err := compress.NewCompressReader(request.Body)
			if err != nil {
				logger.Log.Info("error newCompressReader", err)
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
			request.Body = cr
			defer cr.Close()
		}

		h.ServeHTTP(w, request)
	})
}
