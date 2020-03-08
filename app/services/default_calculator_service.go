package services

import (
	"github.com/frankiennamdi/detection-api/support"
	"math"
	"time"

	"github.com/frankiennamdi/detection-api/models"
)

// default calculator service implementation
// provide calculation services like distance and time differences
type DefaultCalculatorService struct{}

const (
	earthRadiusInKm   = float64(6371)
	earthRadiusInMile = float64(3958)
)

func (service DefaultCalculatorService) degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func (service DefaultCalculatorService) HaversineDistance(fromPoint,
	toPoint *models.GeoPoint) (*models.GeoDistance, error) {
	if fromPoint == nil || toPoint == nil {
		return nil, support.NewIllegalArgumentError("fromPoint or toPoint cannot be nil")
	}

	var deltaLat = service.degreesToRadians(toPoint.Latitude - fromPoint.Latitude)

	var deltaLon = service.degreesToRadians(toPoint.Longitude - fromPoint.Longitude)

	var a = math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(fromPoint.Latitude*(math.Pi/180))*math.Cos(toPoint.Latitude*(math.Pi/180))*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return models.NewGeoDistance(
		math.Round(earthRadiusInKm*c*100)/100, math.Round(earthRadiusInMile*c*100)/100), nil
}

func (service DefaultCalculatorService) TimeDifferenceInHours(currentTimeStamp, previousTimeStamp int64) float64 {
	return time.Unix(currentTimeStamp, 0).Sub(time.Unix(previousTimeStamp, 0)).Hours()
}

func (service DefaultCalculatorService) SpeedToTravelDistanceInMPH(
	eventGeoInfoFrom, eventGeoInfoTo *models.EventGeoInfo) (*float64, error) {
	if eventGeoInfoFrom == nil || eventGeoInfoTo == nil {
		return nil, support.NewIllegalArgumentError("to and from geo information cannot be nil")
	}

	distanceDiff, err := service.HaversineDistance(eventGeoInfoFrom.GeoPoint(), eventGeoInfoTo.GeoPoint())

	if err != nil {
		return nil, nil
	}

	timeDiff := service.TimeDifferenceInHours(eventGeoInfoFrom.EventInfo().Timestamp,
		eventGeoInfoTo.EventInfo().Timestamp)
	speed := math.Abs(math.Round(distanceDiff.Miles() / timeDiff))

	return &speed, nil
}
