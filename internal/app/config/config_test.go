package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {

	t.Setenv("SERVER_ADDRESS", "localhost:8080")
	t.Setenv("BASE_URL", "http://localhost:8080")
	t.Setenv("FILE_STORAGE_PATH", "db.test")
	t.Setenv("DATABASE_DSN", "host=localhost user=user password=password dbname=shortener_url sslmode=disable")
	t.Setenv("ENABLE_HTTPS", "true")

	c := New()
	assert.Equal(t, c, &Config{
		URL:      "http://localhost:8080",
		Addr:     "localhost:8080",
		LogLevel: "info",
		FilePath: "db.test",
		DBPath:   "host=localhost user=user password=password dbname=shortener_url sslmode=disable",
		TLS:      true,
	})
	fmt.Println(c)

}
