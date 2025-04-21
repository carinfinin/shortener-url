package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	"github.com/google/uuid"
	"net/http"
)

var keyAuth = []byte("___allohomora___")

// NameToken type for context.
type NameToken string

const NameCookie NameToken = "token"

var ErrorUserNotFound = errors.New("userID not found or invalid")

type Token string

// GenerateToken generates a unique token.
func GenerateToken() string {
	id := uuid.Must(uuid.NewRandom())
	userID := id[:]
	return hex.EncodeToString(userID)
}

// EncodeToken decodes token.
func EncodeToken(token string) (string, error) {
	userID, err := hex.DecodeString(token)
	if err != nil {
		return "", err
	}
	aesBlock, err := aes.NewCipher(keyAuth)
	if err != nil {
		logger.Log.Error("Error generateToken NewCipher :", err)
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

	dst := aegsm.Seal(nonce, nonce, userID, nil)
	return hex.EncodeToString(dst), nil
}
func DecodeCookie(cookie *http.Cookie) (string, error) {
	if cookie.Value == "" {
		return "", fmt.Errorf("cookie.Value  null")
	}
	return DecryptToken(cookie.Value)
}

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func DecryptToken(encryptedToken string) (string, error) {

	encryptedData, err := hex.DecodeString(encryptedToken)
	if err != nil {
		logger.Log.Error("Error decoding hex string:", err)
		return "", err
	}

	aesBlock, err := aes.NewCipher(keyAuth)
	if err != nil {
		logger.Log.Error("Error creating AES cipher:", err)
		return "", err
	}

	aegsm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		logger.Log.Error("Error creating GCM:", err)
		return "", err
	}

	nonceSize := aegsm.NonceSize()
	if len(encryptedData) < nonceSize {
		logger.Log.Error("Encrypted data too short")
		return "", fmt.Errorf("encrypted data too short")
	}

	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]

	plaintext, err := aegsm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		logger.Log.Error("Error decrypting data:", err)
		return "", err
	}

	return hex.EncodeToString(plaintext), nil
}
