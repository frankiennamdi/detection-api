package main

import (
	"github.com/frankiennamdi/detection-api/core"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/frankiennamdi/detection-api/app"
	"github.com/frankiennamdi/detection-api/config"
	"github.com/frankiennamdi/detection-api/support"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf(support.Info, "starting")

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		sig := <-sigc
		log.Printf(support.Info, sig)
		log.Printf(support.Info, "goodbye")
		os.Exit(0)
	}()

	appConfig := &config.AppConfig{}
	err := appConfig.Read()

	if err != nil {
		log.Panicf(support.Fatal, err)
	}

	server := core.NewServer(*appConfig)
	serverContext := server.Configure()
	serviceContext := app.NewServiceContext(serverContext)
	serviceContext.Listen()
}
