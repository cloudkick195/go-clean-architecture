package main

import (
	"go_clean_architecture/config"
	"go_clean_architecture/server"
	"log"
)

func main() {
	config.InitConfig()
	server.NewServer()
	s := server.NewServer()
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
