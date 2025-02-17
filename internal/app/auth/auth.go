package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	"github.com/google/uuid"
	"net/http"
)

type Auth struct {
	aesBlock cipher.Block
	aeGSM    cipher.AEAD
	nonce    []byte
}

var keyAuth = []byte("___allohomora___")

func New() (*Auth, error) {

	aesBlock, err := aes.NewCipher(keyAuth)
	if err != nil {
		logger.Log.Debug("Error generateToken :", err)
		return nil, err
	}
	aegsm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		logger.Log.Debug("Error NewGCM :", err)
		return nil, err
	}

	nonce, err := generateRandom(aegsm.NonceSize())
	if err != nil {
		logger.Log.Debug("Error generateRandom :", err)
		return nil, err
	}

	return &Auth{aesBlock, aegsm, nonce}, nil
}

func (auth *Auth) GenerateToken() (string, error) {
	id := uuid.Must(uuid.NewRandom())
	userID := id[:]
	//userID := []byte("000007")
	fmt.Println("GenerateToken", string(userID))

	dst := auth.aeGSM.Seal(nil, auth.nonce, userID, nil)
	fmt.Printf("encripted: %x\n", dst)
	fmt.Println(hex.EncodeToString(dst))

	return hex.EncodeToString(dst), nil
}

func (auth *Auth) DecodeCookie(cookie *http.Cookie) bool {
	if cookie.Value == "" {
		logger.Log.Error("cookie.Value  null")
		return true
	}
	dst, err := hex.DecodeString(cookie.Value)
	src2, err := auth.aeGSM.Open(nil, auth.nonce, dst, nil) // расшифровываем
	if err != nil {
		logger.Log.Error("Error decodeCookie  aegsm.Open:", src2)
	}
	fmt.Println(string(src2))
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
func (auth *Auth) Close() {
	auth = nil
}
