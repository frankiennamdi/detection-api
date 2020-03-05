package services

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Router struct {
	serverContext *ServerContext
}

func (router Router) setRoutes(routes *mux.Router) *mux.Router {
	detectionController := EventDetectionController{
		DetectionService: router.serverContext.detectionService,
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
