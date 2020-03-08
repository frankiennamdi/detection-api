package services

import (
	"fmt"
	"log"
	"testing"

	"github.com/frankiennamdi/detection-api/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var testsCases = []struct {
	from                 *models.GeoPoint
	to                   *models.GeoPoint
	expectedDistanceInKm float64
}{
	{
		&models.GeoPoint{Latitude: 22.55, Longitude: 43.12},
		&models.GeoPoint{Latitude: 13.45, Longitude: 100.28},
		6094.54,
	},
	{
		&models.GeoPoint{Latitude: 51.510357, Longitude: -0.116773},
		&models.GeoPoint{Latitude: 38.889931, Longitude: -77.009003},
		5897.66,
	},
	{
		&models.GeoPoint{Latitude: 39.1702, Longitude: -76.8538},
		&models.GeoPoint{Latitude: 34.0494, Longitude: -118.2641},
		3707.35,
	},
}

func TestCalculateDistance(t *testing.T) {
	calculator := DefaultCalculatorService{}
	req := require.New(t)

	for _, input := range testsCases {
		meters, err := calculator.HaversineDistance(input.from, input.to)
		req.NoError(err)

		log.Printf("%+v", meters)
		req.Equal(input.expectedDistanceInKm, meters.Km(),
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
	currentEvent := models.NewEventGeoInfo(&models.EventInfo{
		UUID:      uuid.New().String(),
		Username:  "John",
		Timestamp: 1514764800,
		IP:        "1.0.0.0",
	}, &models.GeoPoint{
		Latitude:       39.1702,
		Longitude:      -76.8538,
		AccuracyRadius: 10,
	})

	subsequentEvent := models.NewEventGeoInfo(&models.EventInfo{
		UUID:      uuid.New().String(),
		Username:  "john",
		Timestamp: 1514851200,
		IP:        "2.0.0.0",
	}, &models.GeoPoint{
		Latitude:       34.0494,
		Longitude:      -118.2641,
		AccuracyRadius: 10,
	})
	speed, err := calculator.SpeedToTravelDistanceInMPH(subsequentEvent, currentEvent)
	req := require.New(t)
	req.NoError(err)
	req.Equal(float64(96), *speed)
}

func TestTimeDifferenceInMinutes(t *testing.T) {
	calculator := DefaultCalculatorService{}
	diff := calculator.TimeDifferenceInHours(1514851200, 1514764800)
	req := require.New(t)
	req.Equal(float64(24), diff)
}

func TestHaversineDistance_When_From_Or_To_Point_Are_Nil(t *testing.T) {
	calculator := DefaultCalculatorService{}
	req := require.New(t)

	_, err1 := calculator.HaversineDistance(nil, nil)
	req.Error(err1)

	_, err2 := calculator.HaversineDistance(nil, &models.GeoPoint{})
	req.Error(err2)

	_, err3 := calculator.HaversineDistance(&models.GeoPoint{}, nil)
	req.Error(err3)
}

func TestSpeedToTravelDistanceInMPH_When_Info_Is_Nil(t *testing.T) {
	calculator := DefaultCalculatorService{}
	req := require.New(t)

	_, err1 := calculator.SpeedToTravelDistanceInMPH(nil, nil)
	req.Error(err1)

	_, err2 := calculator.SpeedToTravelDistanceInMPH(nil, &models.EventGeoInfo{})
	req.Error(err2)

	_, err3 := calculator.SpeedToTravelDistanceInMPH(&models.EventGeoInfo{}, nil)
	req.Error(err3)
}
