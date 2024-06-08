package main

import (
	"log"
	"url-shortener/internal/config"
)

func main() {
	// init config
	cfg := config.MustLoad()
	log.Println(cfg)
	// init logger

	// init storage

	// init router

	// init server
}
