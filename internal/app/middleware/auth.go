package middleware

import (
	"fmt"
	"github.com/carinfinin/shortener-url/internal/app/auth"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	"net/http"
)

func AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		nameCookie := "token"
		auth, err := auth.New()
		defer auth.Close()

		if err != nil {
			logger.Log.Error("auth coll new error :", err)
		}

		c, err := request.Cookie(nameCookie)
		if err != nil || !auth.DecodeCookie(c) {
			valueCookie, err := auth.GenerateToken()
			if err != nil {
				logger.Log.Error("Error generate token :", err)
			}

			fmt.Println("valueCookie: ", valueCookie)

			c := http.Cookie{
				Name:   nameCookie,
				Value:  valueCookie,
				MaxAge: 300,
			}
			http.SetCookie(writer, &c)
		}
		fmt.Println(c)

		next.ServeHTTP(writer, request)
	})
}
