package core

import (
	"database/sql"
	"net"

	"github.com/frankiennamdi/detection-api/models"
)

type EventRepository interface {
	InsertEvents(events []*models.Event) ([]sql.Result, error)
	FindRelatedEvents(event *models.Event, filter EventFilter) error
}

type EventFilter interface {
	Filter(event *models.Event)
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
