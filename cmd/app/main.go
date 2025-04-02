package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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
	application := app.New(cfg)
	go application.GRPC.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)
	<-stop
	application.GRPC.Stop()
}
