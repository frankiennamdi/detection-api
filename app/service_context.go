package app

import (
	"fmt"
	"github.com/frankiennamdi/detection-api/app/services"
	"github.com/frankiennamdi/detection-api/core"
	"github.com/frankiennamdi/detection-api/repository"
	"github.com/frankiennamdi/detection-api/support"
	"log"
	"net/http"
)

// provides a context for the services of the server. It initializes all the services and
// provides a function for initializing the routes and listening for connections
type ServiceContext struct {
	detectionService core.DetectionService
	eventRepository  core.EventRepository
	server *core.ServerContext
}

func NewServiceContext(ctx *core.ServerContext) *ServiceContext {
	eventRepository := repository.NewSQLLiteEventsRepository(ctx.EventDb())
	detectionService := services.NewDetectionService(eventRepository,
		repository.NewMaxMindIPGeoInfoRepository(ctx.GeoIPDb()),
		services.DefaultCalculatorService{},
		ctx.AppConfig().SuspiciousSpeed)

	return &ServiceContext{
		detectionService: detectionService,
		eventRepository:  eventRepository,
		server: ctx,
	}
}

func (serviceContext *ServiceContext) DetectionService() core.DetectionService {
	return serviceContext.detectionService
}

func (serviceContext *ServiceContext) EventRepository() core.EventRepository {
	return serviceContext.eventRepository
}

func (serviceContext *ServiceContext) Listen() {
	router := Router{serviceContext: serviceContext}
	routes := router.InitRoutes()

	log.Printf(support.Info, fmt.Sprintf("service starting on Port : %d ...",
		serviceContext.server.AppConfig().Server.Port))

	if err := http.ListenAndServe(fmt.Sprintf(":%d",
		serviceContext.server.AppConfig().Server.Port), routes); err != nil {
		log.Panicf(support.Fatal, err)
	}
}

