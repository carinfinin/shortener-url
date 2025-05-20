package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
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

func TestReadConfigJSON(t *testing.T) {
	var fname = "data.json"

	tests := []struct {
		name  string
		data  string
		err   bool
		res   Config
		rfile string
	}{
		{
			name:  "success",
			data:  `{"server_address": "localhost:8080","base_url": "http://localhost","file_storage_path": "","database_dsn": "","enable_https": true}`,
			err:   false,
			rfile: fname,
			res: Config{
				Addr:     "localhost:8080",
				URL:      "http://localhost",
				FilePath: "",
				DBPath:   "",
				TLS:      true,
			},
		},
		{
			name:  "error",
			data:  `{"server_address": "localhost:8080","base_url": "http://localhost","file_storage_path": "","database_dsn": "","enable_https": true`,
			err:   true,
			rfile: fname,
			res:   Config{},
		},
		{
			name:  "error2",
			data:  `{"server_address": "localhost:8080","base_url": "http://localhost","file_storage_path": "","database_dsn": "","enable_https": true}`,
			err:   true,
			rfile: "sdf.er",
			res:   Config{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			f, err := os.Create(fname)
			assert.NoError(t, err, err)

			_, err = f.WriteString(test.data)
			assert.NoError(t, err, err)
			f.Close()

			cfg := Config{}

			err = readConfigJSON(test.rfile, &cfg)
			if !test.err {
				assert.NoError(t, err, err)
				assert.Equal(t, test.res, cfg)
			} else {
				assert.Error(t, err)
			}

			err = os.Remove(fname)
			assert.NoError(t, err, err)

		})
	}

}
