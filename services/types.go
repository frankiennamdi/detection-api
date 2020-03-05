package services

import (
	"database/sql"
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
}

type CalculatorService interface {
	HaversineDistance(fromPoint, toPoint models.GeoPoint) models.GeoDistance
	TimeDifferenceInHours(currentTimeStamp, previousTimeStamp int64) float64
	SpeedToTravelDistanceInMPH(eventGeoInfoFrom, eventGeoInfoTo models.EventGeoInfo) float64
}
