package app

import (
	"github.com/frankiennamdi/detection-api/app/services"
	"github.com/frankiennamdi/detection-api/core"
	"github.com/frankiennamdi/detection-api/repository"
)

type ServiceContext struct {
	detectionService core.DetectionService
	eventRepository  core.EventRepository
}

func NewServiceContext(ctx *core.ServerContext) *ServiceContext {
	eventRepository := repository.NewSQLLiteEventsRepository(ctx.EventDb())
	detectionService := services.NewLoginDetectionService(eventRepository,
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

