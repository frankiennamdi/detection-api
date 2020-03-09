package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

// store the routes and handlers
type Router struct {
	serviceContext *ServiceContext
}

func (router Router) setRoutes(routes *mux.Router) *mux.Router {
	detectionController := EventDetectionController{
		detectionService: router.serviceContext.DetectionService(),
	}

	routes.HandleFunc("/api/health-check", StatusHandler).Methods(http.MethodGet)
	routes.HandleFunc("/api/events", detectionController.EventDetectionHandler).Methods(http.MethodPost)

	return routes
}

func (router Router) InitRoutes() *mux.Router {
	routes := mux.NewRouter().StrictSlash(false)
	configuredRoutes := router.setRoutes(routes)

	return configuredRoutes
}
