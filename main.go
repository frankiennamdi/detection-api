package main

import (
	"fmt"
	"github.com/frankiennamdi/detection-api/core"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/frankiennamdi/detection-api/config"
	"github.com/frankiennamdi/detection-api/app"
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

	server := core.Server{Config: *appConfig}
	serverContext := server.Configure()
	Listen(serverContext, app.NewServiceContext(serverContext))
}

func Listen(serverContext *core.ServerContext, ctx *app.ServiceContext) {
	log.Printf(support.Info, fmt.Sprintf("service starting on Port : %d ...", serverContext.AppConfig().Server.Port))

	router := app.Router{ServiceContext: ctx}
	routes := router.InitRoutes()

	if err := http.ListenAndServe(fmt.Sprintf(":%d", serverContext.AppConfig().Server.Port), routes); err != nil {
		log.Fatal(err)
	}
}
