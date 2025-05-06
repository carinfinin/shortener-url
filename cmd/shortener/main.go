package main

import (
	_ "embed"
	"fmt"
	"github.com/carinfinin/shortener-url/internal/app/config"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	"github.com/carinfinin/shortener-url/internal/app/server"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

//go:embed version.txt
var gv string

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func printGlobalVar() {

	lines := strings.Split(gv, "\n")

	if buildVersion == "" && len(lines) > 0 {
		buildVersion = lines[0]
	} else {
		buildVersion = "N/A"
	}
	fmt.Printf("Build version: %s\n", buildVersion)

	if buildDate == "" && len(lines) > 1 {
		buildDate = lines[1]
	} else {
		buildDate = "N/A"
	}
	fmt.Printf("Build date: %s\n", buildDate)

	if buildCommit == "" && len(lines) > 2 {
		buildCommit = lines[2]
	} else {
		buildCommit = "N/A"
	}
	fmt.Printf("Build commit: %s\n", buildCommit)
}

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

	printGlobalVar()

	logger.Log.Info("server started")

	<-exit
	s.Store.Close()
	logger.Log.Info("stopping server")

}
