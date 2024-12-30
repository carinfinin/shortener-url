package config

import "flag"

type Config struct {
	Addr string
	URL  string
}

func New() *Config {
	var config Config

	flag.StringVar(&config.Addr, "a", "localhost:8080", "url address server")
	flag.StringVar(&config.URL, "b", "http://localhost:8080/", "result short url")
	flag.Parse()
	return &config
}
