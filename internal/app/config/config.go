package config

import (
	"encoding/json"
	"flag"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	"os"
)

// Config contains application settings.
type Config struct {
	Addr     string `json:"server_address"`
	URL      string `json:"base_url"`
	LogLevel string `json:"-"`
	FilePath string `json:"file_storage_path"`
	DBPath   string `json:"database_dsn"`
	TLS      bool   `json:"enable_https"`
	path     string `json:"-"`
}

// New constructor for type Config.
func New() *Config {
	var config Config

	flag.StringVar(&config.Addr, "a", "localhost:8080", "url address server")
	flag.StringVar(&config.URL, "b", "http://localhost:8080", "result short url")
	flag.StringVar(&config.LogLevel, "level", "info", "log level")
	/* data.json */
	flag.StringVar(&config.FilePath, "f", "", "file path")
	/* host=localhost user=user password=password dbname=shortener_url sslmode=disable */
	flag.StringVar(&config.DBPath, "d", "", "db path")
	flag.BoolVar(&config.TLS, "s", false, "tls")
	flag.StringVar(&config.path, "config", "", "config file path")
	flag.StringVar(&config.path, "c", "", "config file path")
	flag.Parse()

	if envServerAddr := os.Getenv("SERVER_ADDRESS"); envServerAddr != "" {
		config.Addr = envServerAddr
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		config.URL = envBaseURL
	}
	if envFilePath := os.Getenv("FILE_STORAGE_PATH"); envFilePath != "" {
		config.FilePath = envFilePath
	}
	if DBPath := os.Getenv("DATABASE_DSN"); DBPath != "" {
		config.DBPath = DBPath
	}
	if tls := os.Getenv("ENABLE_HTTPS"); tls != "" {
		config.TLS = tls == "true"
	}
	if path := os.Getenv("CONFIG"); path != "" {
		config.path = path

	}

	if config.path != "" {
		err := readConfigJSON(config.path, &config)
		if err != nil {
			logger.Log.Errorf("error read config path %s: %s", config.path, err)
		}
	}

	return &config
}

func readConfigJSON(fname string, cfg *Config) error {
	b, err := os.ReadFile(fname)
	if err != nil {
		return err
	}
	var cNew Config
	err = json.Unmarshal(b, &cNew)
	if err != nil {
		return err
	}

	if cfg.Addr == "" {
		cfg.Addr = cNew.Addr
	}
	if cfg.URL == "" {
		cfg.URL = cNew.URL
	}
	if cfg.FilePath == "" {
		cfg.FilePath = cNew.FilePath
	}
	if cfg.DBPath == "" {
		cfg.DBPath = cNew.DBPath
	}
	if !cfg.TLS {
		cfg.TLS = cNew.TLS
	}

	return nil
}
