package main

import (
	"github.com/frankiennamdi/detection-api/config"
	"github.com/frankiennamdi/detection-api/services"
	"github.com/frankiennamdi/detection-api/support"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Printf(support.Info, "Starting")

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		sig := <-sigc
		log.Printf(support.Info, sig)
		log.Printf(support.Info, "Goodbye")
		os.Exit(0)
	}()

	appConfig := &config.AppConfig{}
	err := appConfig.Read()

	if err != nil {
		log.Panicf(support.Fatal, err)
	}

	server := services.Server{Config: *appConfig}
	serverContext := server.Configure()
	serverContext.Listen()
}
