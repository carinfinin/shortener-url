package auth

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func BenchmarkEncodeToken(b *testing.B) {
	token := GenerateToken()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EncodeToken(token)
	}
}

func BenchmarkGenerateToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateToken()
	}
}

func TestDecodeCookie(t *testing.T) {

	t.Run("Positive test", func(t *testing.T) {
		token := GenerateToken()
		cv, err := EncodeToken(token)
		require.NoError(t, err, err)

		c := http.Cookie{
			Name:    string(NameCookie),
			Value:   cv,
			Expires: time.Now().Add(time.Hour * 24),
			Path:    "/",
		}

		v, err := DecodeCookie(&c)
		require.NoError(t, err, err)

		assert.Equal(t, token, v)
	})

	t.Run("Fail test", func(t *testing.T) {
		token := GenerateToken()

		c := http.Cookie{
			Name:    string(NameCookie),
			Value:   token,
			Expires: time.Now().Add(time.Hour * 24),
			Path:    "/",
		}

		v, err := DecodeCookie(&c)
		assert.Error(t, err)
		assert.Equal(t, "", v)

	})

}
