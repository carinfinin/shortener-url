package main

import (
	"github.com/carinfinin/shortener-url/internal/app/server"
)

func main() {

	s := server.New("8080")
	err := s.Start()

	if err != nil {
		panic(err)
	}
}
