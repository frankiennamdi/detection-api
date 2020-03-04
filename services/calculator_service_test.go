package services

import (
	"fmt"
	"github.com/frankiennamdi/detection-api/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

var testsCases = []struct {
	from                 models.GeoPoint
	to                   models.GeoPoint
	expectedDistanceInKm float64
}{
	{
		models.GeoPoint{Latitude: 22.55, Longitude: 43.12},
		models.GeoPoint{Latitude: 13.45, Longitude: 100.28},
		6094.54,
	},
	{
		models.GeoPoint{Latitude: 51.510357, Longitude: -0.116773},
		models.GeoPoint{Latitude: 38.889931, Longitude: -77.009003},
		5897.66,
	},
	{
		models.GeoPoint{Latitude: 39.1702, Longitude: -76.8538},
		models.GeoPoint{Latitude: 34.0494, Longitude: -118.2641},
		3707.35,
	},
}

func TestCalculateDistance(t *testing.T) {
	calculator := DefaultCalculatorService{}
	req := require.New(t)

	for _, input := range testsCases {
		meters := calculator.HaversineDistance(input.from, input.to)
		log.Printf("%+v", meters)
		req.Equal(input.expectedDistanceInKm, meters.Km,
			fmt.Sprintf("FAIL: Want distance from %v to %v to be: %v but we got %v",
				input.from,
				input.to,
				input.expectedDistanceInKm,
				meters,
			))
	}
}

func TestSpeedToTravelDistanceInMPH(t *testing.T) {
	calculator := DefaultCalculatorService{}
	currentEvent := models.EventGeoInfo{
		EventInfo: &models.EventInfo{
			UUID:      uuid.New().String(),
			Username:  "John",
			Timestamp: 1514764800,
			IP:        "1.0.0.0",
		},
		GeoPoint: &models.GeoPoint{
			Latitude:       39.1702,
			Longitude:      -76.8538,
			AccuracyRadius: 10,
		},
	}

	subsequentEvent := models.EventGeoInfo{
		EventInfo: &models.EventInfo{
			UUID:      uuid.New().String(),
			Username:  "john",
			Timestamp: 1514851200,
			IP:        "2.0.0.0",
		},
		GeoPoint: &models.GeoPoint{
			Latitude:       34.0494,
			Longitude:      -118.2641,
			AccuracyRadius: 10,
		},
	}
	speed := calculator.SpeedToTravelDistanceInMPH(subsequentEvent, currentEvent)
	req := require.New(t)
	req.Equal(float64(96), speed)
}

func TestTimeDifferenceInMinutes(t *testing.T) {
	calculator := DefaultCalculatorService{}
	diff := calculator.TimeDifferenceInHours(1514851200, 1514764800)
	req := require.New(t)
	req.Equal(float64(24), diff)
}
