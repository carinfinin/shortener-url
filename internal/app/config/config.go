package config

import (
	"flag"
	"os"
)

type Config struct {
	Addr     string
	URL      string
	LogLevel string
	FilePath string
	DBPath   string
}

func New() *Config {
	var config Config

	flag.StringVar(&config.Addr, "a", "localhost:8080", "url address server")
	flag.StringVar(&config.URL, "b", "http://localhost:8080", "result short url")
	flag.StringVar(&config.LogLevel, "level", "info", "log level")
	/* data.json */
	flag.StringVar(&config.FilePath, "f", "", "file path")
	/* host=localhost user=user password=password dbname=shortener_url sslmode=disable */
	flag.StringVar(&config.DBPath, "d", "", "db path")
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
		config.FilePath = DBPath
	}
	return &config
}
