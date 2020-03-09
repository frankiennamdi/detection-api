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
	}
}

func (serviceContext *ServiceContext) DetectionService() core.DetectionService {
	return serviceContext.detectionService
}

func (serviceContext *ServiceContext) EventRepository() core.EventRepository {
	return serviceContext.eventRepository
}

func (serviceContext *ServiceContext) Listen(port int) {
	router := Router{serviceContext: serviceContext}
	routes := router.InitRoutes()

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), routes); err != nil {
		log.Panicf(support.Fatal, err)
	}
}

