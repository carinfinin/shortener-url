package middleware

import (
	"context"
	"github.com/carinfinin/shortener-url/internal/app/auth"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	"net/http"
	"time"
)

func AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		var token string
		c, err := request.Cookie(string(auth.NameCookie))
		if err != nil {
			token = auth.GenerateToken()

			valueCookie, err := auth.EncodeToken(token)
			if err != nil {
				logger.Log.Debug("Error generate token EncodeToken :", err)
			}

			c := http.Cookie{
				Name:    string(auth.NameCookie),
				Value:   valueCookie,
				Expires: time.Now().Add(time.Hour * 24),
				Path:    "/",
			}
			http.SetCookie(writer, &c)

		} else {
			token, err = auth.DecodeCookie(c)
			if err != nil {
				logger.Log.Error(err)
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}

		}
		ctx := context.WithValue(request.Context(), auth.NameCookie, token)
		newReq := request.WithContext(ctx)
		next.ServeHTTP(writer, newReq)
	})
}
