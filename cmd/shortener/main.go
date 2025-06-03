package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/carinfinin/shortener-url/internal/app/config"
	"github.com/carinfinin/shortener-url/internal/app/grpcserver"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	"github.com/carinfinin/shortener-url/internal/app/server"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	cfg := config.New()

	err := logger.Configure(cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	logger.Log.Info("server starting")

	s, err := server.New(cfg)
	if err != nil {
		logger.Log.Error("server starting error:", err)
		panic(err)
	}

	go func() {
		err = grpcserver.Start(cfg)
		if err != nil {
			logger.Log.Error("grps server starting error:", err)
			panic(err)
		}
	}()

	go func() {
		if er := s.Start(); er != nil {
			logger.Log.Error("server failed error: ", er)
		}
	}()

	printGlobalVar()

	logger.Log.Info("server started")

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err = s.Stop(shutdownCtx); err != nil {
		logger.Log.Error("error stop server: ", err)
	}
	// закрытие воркеров в сервисе
	s.Service.Close()

	if err = s.Store.Close(); err != nil {
		logger.Log.Error("error stop store: ", err)
	}
	logger.Log.Info("stopping server")
}
