package main

import (
	"fmt"
	"github.com/carinfinin/shortener-url/internal/app/config"
	"github.com/carinfinin/shortener-url/internal/app/server"
)

func main() {

	config := config.New()

	s := server.New(config)
	err := s.Start()

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
