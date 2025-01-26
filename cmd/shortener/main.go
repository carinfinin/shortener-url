package main

import (
	"fmt"
	"github.com/carinfinin/shortener-url/internal/app/config"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	"github.com/carinfinin/shortener-url/internal/app/server"
)

func main() {

	config := config.New()

	err := logger.Configure(config.LogLevel)
	if err != nil {
		panic(err)
	}
	logger.Log.Info("server starting")

	s, err := server.New(config)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	err = s.Start()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	logger.Log.Info("server started")

}
