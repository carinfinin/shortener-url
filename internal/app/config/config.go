package config

import (
	"flag"
	"os"
)

type Config struct {
	Addr     string
	URL      string
	LogLevel string
}

func New() *Config {
	var config Config

	flag.StringVar(&config.Addr, "a", "localhost:8080", "url address server")
	flag.StringVar(&config.URL, "b", "http://localhost:8080", "result short url")
	flag.StringVar(&config.LogLevel, "level", "info", "log level")
	flag.Parse()

	if envServerAddr := os.Getenv("SERVER_ADDRESS"); envServerAddr != "" {
		config.Addr = envServerAddr
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		config.URL = envBaseURL
	}
	return &config
}
