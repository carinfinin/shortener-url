package main

import (
	"fmt"
	"github.com/carinfinin/shortener-url/internal/app/config"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	"github.com/carinfinin/shortener-url/internal/app/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

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

	go func() {
		if err := s.Start(); err != nil {
			logger.Log.Info("server failed")
		}
	}()

	logger.Log.Info("server started")

	<-exit
	s.Store.Close()
	logger.Log.Info("stopping server")

}
