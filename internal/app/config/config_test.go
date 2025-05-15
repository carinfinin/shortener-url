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

	c := New()
	assert.Equal(t, c, &Config{
		URL:      "http://localhost:8080",
		Addr:     "localhost:8080",
		LogLevel: "info",
		FilePath: "db.test",
	})
	fmt.Println(c)

}
