package services

import (
	"database/sql"
	"github.com/frankiennamdi/detection-api/config"
	"github.com/frankiennamdi/detection-api/db"
	"github.com/frankiennamdi/detection-api/models"
	"net"
)

type EventRepository interface {
	Insert(events []*models.Event) ([]sql.Result, error)
	FindByUsername(username string) ([]*models.Event, error)
}

type IPGeoInfoRepository interface {
	FindGeoPoint(IP net.IP) (*models.GeoPoint, error)
}

type DetectionService interface {
	ProcessEvent(currEvent *models.Event) (*models.SuspiciousTravelResult, error)
	FindSuspiciousTravel(relatedEventInfo *models.RelatedEventInfo) (*models.SuspiciousTravelResult, error)
	FindRelatedEvents(currEvent *models.Event) (*models.RelatedEventInfo, error)
}

type Server struct {
	Config config.AppConfig
}

type MaxMindIPGeoInfoRepository struct {
	maxMindDb *db.MaxMindDb
}
