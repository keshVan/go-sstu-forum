package main

import (
	"log"

	"github.com/keshvan/go-sstu-forum/config"
	"github.com/keshvan/go-sstu-forum/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
