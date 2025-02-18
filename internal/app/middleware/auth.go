package middleware

import (
	"context"
	"fmt"
	"github.com/carinfinin/shortener-url/internal/app/auth"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	"net/http"
)

func AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		c, err := request.Cookie(auth.NameCookie)
		if err != nil {
			token := auth.GenerateToken()

			valueCookie, err := auth.EncodeToken(token)
			if err != nil {
				logger.Log.Error("Error generate token EncodeToken :", err)
			}

			c := http.Cookie{
				Name:   auth.NameCookie,
				Value:  valueCookie,
				MaxAge: 300,
			}
			http.SetCookie(writer, &c)

			ctx := context.WithValue(request.Context(), auth.NameCookie, token)

			newReq := request.WithContext(ctx)
			next.ServeHTTP(writer, newReq)
		} else {
			token, err := auth.DecodeCookie(c)
			if err != nil {
				logger.Log.Error(err)
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}

			fmt.Println("---------------")
			fmt.Println(token)
			fmt.Println("---------------")

			ctx := context.WithValue(request.Context(), auth.NameCookie, token)
			newReq := request.WithContext(ctx)
			next.ServeHTTP(writer, newReq)
		}

	})
}
