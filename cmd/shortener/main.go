package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/carinfinin/shortener-url/internal/app/config"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	"github.com/carinfinin/shortener-url/internal/app/server"
	"os"
	"os/signal"
	"syscall"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func printGlobalVar() {

	if buildVersion == "" {
		buildVersion = "N/A"
	}
	fmt.Printf("Build version: %s\n", buildVersion)

	if buildDate == "" {
		buildDate = "N/A"
	}
	fmt.Printf("Build date: %s\n", buildDate)

	if buildCommit == "" {
		buildCommit = "N/A"
	}
	fmt.Printf("Build commit: %s\n", buildCommit)
}

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	cfg := config.New()

	err := logger.Configure(cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	logger.Log.Info("server starting")

	s, err := server.New(cfg)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	go func() {
		if err := s.Start(); err != nil {
			logger.Log.Info("server failed")
		}
	}()

	printGlobalVar()

	logger.Log.Info("server started")

	<-ctx.Done()
	s.Stop(ctx)
	s.Store.Close()
	logger.Log.Info("stopping server")

}
