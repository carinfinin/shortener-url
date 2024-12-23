package main

import (
	"fmt"
	"github.com/carinfinin/shortener-url/internal/app/server"
)

func main() {

	// storage
	// server
	s := server.New("8081")
	fmt.Println("create server struct")
	err := s.Start()

	if err != nil {
		panic(err)

	}
}
