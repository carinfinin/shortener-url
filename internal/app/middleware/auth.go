package middleware

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	"net/http"
)

func AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		nameCookie := "token"

		c, err := request.Cookie(nameCookie)
		if err != nil || !decodeCookie(c) {
			valueCookie, err := generateToken()
			if err != nil {
				logger.Log.Error("Error generate token :", err)
			}

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

var keyAuth = []byte("___allohomora___")

func generateToken() (string, error) {
	//id := uuid.Must(uuid.NewRandom())
	//userID := id[:]
	userID := []byte("000007")
	fmt.Println(string(userID))

	aesBlock, err := aes.NewCipher(keyAuth)
	if err != nil {
		logger.Log.Error("Error generateToken :", err)
		return "", err
	}

	aegsm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		logger.Log.Error("Error NewGCM :", err)
		return "", err
	}

	nonce, err := generateRandom(aegsm.NonceSize())
	if err != nil {
		logger.Log.Error("Error generateRandom :", err)
		return "", err
	}

	dst := aegsm.Seal(nil, nonce, userID, nil)
	fmt.Printf("encripted: %x\n", dst)
	fmt.Println(string(dst))

	return fmt.Sprint("%x", dst), nil
}

func decodeCookie(cookie *http.Cookie) bool {
	dst := make([]byte, aes.BlockSize)
	aesBlock, err := aes.NewCipher(keyAuth)
	if err != nil {
		logger.Log.Error("Error decodeCookie :", err)
		return true
	}
	if cookie.Value == "" {
		logger.Log.Error("cookie.Value  null")
		return true
	}

	aegsm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		logger.Log.Error("Error NewGCM decodeCookie :", err)
		return true
	}
	nonce, err := generateRandom(aegsm.NonceSize())
	if err != nil {
		logger.Log.Error("Error generateRandom :", err)
		return true
	}

	src2, err := aegsm.Open(nil, nonce, dst, nil) // расшифровываем
	if err != nil {
		logger.Log.Error("Error decodeCookie  aegsm.Open:", err)

		return true
	}
	fmt.Printf("decripted: %s\n", src2)

	return false
}

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
