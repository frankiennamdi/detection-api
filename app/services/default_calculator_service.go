package services

import (
	"math"
	"time"

	"github.com/frankiennamdi/detection-api/models"
)

type DefaultCalculatorService struct{}

const (
	earthRadiusInKm   = float64(6371)
	earthRadiusInMile = float64(3958)
)

func (service DefaultCalculatorService) degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func (service DefaultCalculatorService) HaversineDistance(fromPoint, toPoint models.GeoPoint) models.GeoDistance {
	var deltaLat = service.degreesToRadians(toPoint.Latitude - fromPoint.Latitude)

	var deltaLon = service.degreesToRadians(toPoint.Longitude - fromPoint.Longitude)

	var a = math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(fromPoint.Latitude*(math.Pi/180))*math.Cos(toPoint.Latitude*(math.Pi/180))*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return models.GeoDistance{
		Km:    math.Round(earthRadiusInKm*c*100) / 100,
		Miles: math.Round(earthRadiusInMile*c*100) / 100,
	}
}

func (service DefaultCalculatorService) TimeDifferenceInHours(currentTimeStamp, previousTimeStamp int64) float64 {
	return time.Unix(currentTimeStamp, 0).Sub(time.Unix(previousTimeStamp, 0)).Hours()
}

func (service DefaultCalculatorService) SpeedToTravelDistanceInMPH(
	eventGeoInfoFrom, eventGeoInfoTo models.EventGeoInfo) float64 {
	distanceDiff := service.HaversineDistance(*eventGeoInfoFrom.GeoPoint, *eventGeoInfoTo.GeoPoint)

	timeDiff := service.TimeDifferenceInHours(eventGeoInfoFrom.EventInfo.Timestamp,
		eventGeoInfoTo.EventInfo.Timestamp)

	return math.Abs(math.Round(distanceDiff.Miles / timeDiff))
}
