package main

import (
	"log"
	"video-cutter/internal/app"
)

func main() {
	cfg := app.NewConfig()

	srv := app.NewServer(cfg)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server shutdown with error: %v", err)
	}
}
