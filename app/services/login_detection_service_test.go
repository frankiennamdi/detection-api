package services

import (
	"database/sql"
	"github.com/frankiennamdi/detection-api/core"
	"github.com/frankiennamdi/detection-api/test"
	"log"
	"math"
	"net"
	"testing"
	"time"

	"github.com/frankiennamdi/detection-api/models"
	"github.com/frankiennamdi/detection-api/support"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type MockEventRepository struct {
	userEvents []*models.Event
}

type MockIPGeoInfoRepository struct {
	geoMap map[string]*models.GeoPoint
}

type MockCalculatorService struct{}

func (mockCalculatorService MockCalculatorService) HaversineDistance(fromPoint,
	toPoint models.GeoPoint) models.GeoDistance {
	log.Printf(support.Info, fromPoint)
	log.Printf(support.Info, toPoint)

	return models.GeoDistance{
		Miles: 0,
		Km:    0,
	}
}
func (mockCalculatorService MockCalculatorService) TimeDifferenceInHours(currentTimeStamp,
	previousTimeStamp int64) float64 {
	log.Printf(support.Info, currentTimeStamp)
	log.Printf(support.Info, previousTimeStamp)

	return 0
}

func (mockCalculatorService MockCalculatorService) SpeedToTravelDistanceInMPH(eventGeoInfoFrom,
	eventGeoInfoTo models.EventGeoInfo) float64 {
	return math.Abs(eventGeoInfoFrom.GeoPoint.Longitude - eventGeoInfoTo.GeoPoint.Longitude/
		float64(eventGeoInfoFrom.EventInfo.Timestamp-eventGeoInfoTo.EventInfo.Timestamp))
}

func (mockIPGeoInfoRepository MockIPGeoInfoRepository) FindGeoPoint(ip net.IP) (*models.GeoPoint, error) {
	if value, ok := mockIPGeoInfoRepository.geoMap[ip.String()]; ok {
		return value, nil
	}

	return nil, nil
}

func (mockEventRepo *MockEventRepository) FindRelatedEvents(event *models.Event, filter core.EventFilter) error {
	log.Printf(support.Info, event)

	for _, iterEvent := range mockEventRepo.userEvents {
		filter.Filter(iterEvent)
	}

	return nil
}

func (mockEventRepo *MockEventRepository) InsertEvents(events []*models.Event) ([]sql.Result, error) {
	log.Printf(support.Info, events)
	return nil, nil
}

func TestFindSuspiciousTravelInfo_Fail_When_No_Geo_Info_For_Current_Event(t *testing.T) {
	req := require.New(t)
	detectionService := NewLoginDetectionService(nil,
		&MockIPGeoInfoRepository{geoMap: make(map[string]*models.GeoPoint)},
		nil,
		500)
	event := newEvent(models.EventInfo{
		UUID:      uuid.New().String(),
		Username:  "john",
		Timestamp: 0,
		IP:        "1.0.0.0",
	})

	_, processErr := detectionService.findSuspiciousTravel(&models.RelatedEventInfo{
		CurrentEvent:    event,
		PreviousEvent:   nil,
		SubsequentEvent: nil,
	})

	req.Error(processErr)
}

func TestFindSuspiciousTravelInfo_When_No_Previous_Or_Subsequent_Event(t *testing.T) {
	req := require.New(t)
	detectionService := NewLoginDetectionService(nil,
		&MockIPGeoInfoRepository{geoMap: map[string]*models.GeoPoint{
			"1.0.0.0": {
				Latitude:       10,
				Longitude:      10,
				AccuracyRadius: 10,
			},
		}},
		nil,
		500)
	event := newEvent(models.EventInfo{
		UUID:      uuid.New().String(),
		Username:  "john",
		Timestamp: 0,
		IP:        "1.0.0.0",
	})

	result, err := detectionService.findSuspiciousTravel(&models.RelatedEventInfo{
		CurrentEvent:    event,
		PreviousEvent:   nil,
		SubsequentEvent: nil,
	})

	req.NoError(err)
	req.NotNil(result.CurrentGeo)
	req.Nil(result.SubsequentIPAccess)
	req.Nil(result.PrecedingIPAccess)
	req.Nil(result.TravelFromCurrentGeoSuspicious)
	req.Nil(result.TravelToCurrentGeoSuspicious)
}

func TestFindSuspiciousTravelInfo_When_No_Previous_But_Subsequent_Event(t *testing.T) {
	req := require.New(t)
	detectionService := NewLoginDetectionService(nil,
		&MockIPGeoInfoRepository{geoMap: map[string]*models.GeoPoint{
			"1.0.0.0": {
				Latitude:       10,
				Longitude:      10,
				AccuracyRadius: 10,
			},
			"1.1.0.0": {
				Latitude:       100,
				Longitude:      10,
				AccuracyRadius: 10,
			},
		}},
		&MockCalculatorService{},
		10)

	result, err := detectionService.findSuspiciousTravel(&models.RelatedEventInfo{
		CurrentEvent: newEvent(models.EventInfo{
			UUID:      uuid.New().String(),
			Username:  "john",
			Timestamp: 100,
			IP:        "1.0.0.0",
		}),
		PreviousEvent: nil,
		SubsequentEvent: newEvent(models.EventInfo{
			UUID:      uuid.New().String(),
			Username:  "john",
			Timestamp: 200,
			IP:        "1.1.0.0",
		}),
	})

	req.NoError(err)
	req.NotNil(result.CurrentGeo)
	req.NotNil(result.SubsequentIPAccess)
	req.Nil(result.PrecedingIPAccess)
	req.Equal(true, *result.TravelFromCurrentGeoSuspicious)
	req.Nil(result.TravelToCurrentGeoSuspicious)
}

func TestFindSuspiciousTravelInfo_When_No_Subsequent_But_Previous_Event(t *testing.T) {
	req := require.New(t)
	detectionService := NewLoginDetectionService(nil,
		&MockIPGeoInfoRepository{geoMap: map[string]*models.GeoPoint{
			"1.0.0.0": {
				Latitude:       10,
				Longitude:      10,
				AccuracyRadius: 10,
			},
			"1.1.0.0": {
				Latitude:       100,
				Longitude:      10,
				AccuracyRadius: 10,
			},
		}},
		&MockCalculatorService{},
		10)

	result, err := detectionService.findSuspiciousTravel(&models.RelatedEventInfo{
		CurrentEvent: newEvent(models.EventInfo{
			UUID:      uuid.New().String(),
			Username:  "john",
			Timestamp: 100,
			IP:        "1.0.0.0",
		}),
		PreviousEvent: newEvent(models.EventInfo{
			UUID:      uuid.New().String(),
			Username:  "john",
			Timestamp: 200,
			IP:        "1.1.0.0",
		}),
		SubsequentEvent: nil,
	})

	req.NoError(err)
	req.NotNil(result.CurrentGeo)
	req.Nil(result.SubsequentIPAccess)
	req.NotNil(result.PrecedingIPAccess)
	req.Equal(true, *result.TravelToCurrentGeoSuspicious)
	req.Nil(result.TravelFromCurrentGeoSuspicious)
}

func TestFindSuspiciousTravelInfo_When_All_Events_Present(t *testing.T) {
	req := require.New(t)
	detectionService := NewLoginDetectionService(nil,
		&MockIPGeoInfoRepository{geoMap: map[string]*models.GeoPoint{
			"1.0.0.0": {
				Latitude:       10,
				Longitude:      10,
				AccuracyRadius: 10,
			},
			"1.1.0.0": {
				Latitude:       100,
				Longitude:      10,
				AccuracyRadius: 10,
			},
			"1.2.0.0": {
				Latitude:       100,
				Longitude:      10,
				AccuracyRadius: 10,
			},
		}},
		&MockCalculatorService{},
		10)

	result, err := detectionService.findSuspiciousTravel(&models.RelatedEventInfo{
		CurrentEvent: newEvent(models.EventInfo{
			UUID:      uuid.New().String(),
			Username:  "john",
			Timestamp: 100,
			IP:        "1.0.0.0",
		}),
		PreviousEvent: newEvent(models.EventInfo{
			UUID:      uuid.New().String(),
			Username:  "john",
			Timestamp: 200,
			IP:        "1.1.0.0",
		}),
		SubsequentEvent: newEvent(models.EventInfo{
			UUID:      uuid.New().String(),
			Username:  "john",
			Timestamp: 200,
			IP:        "1.2.0.0",
		}),
	})

	req.NoError(err)
	req.NotNil(result.CurrentGeo)
	req.NotNil(result.SubsequentIPAccess)
	req.NotNil(result.PrecedingIPAccess)
	req.Equal(true, *result.TravelToCurrentGeoSuspicious)
	req.Equal(true, *result.TravelFromCurrentGeoSuspicious)
}

func TestFindRelatedEvents(t *testing.T) {
	req := require.New(t)
	currentEvent := newEvent(models.EventInfo{
		UUID:      uuid.New().String(),
		Username:  "john",
		Timestamp: 1514764800,
		IP:        "1.0.0.0",
	})

	var events []*models.Event

	closestPreviousEvent, previousEvents := generateEvents(currentEvent, -1)

	events = append(events, previousEvents...)

	closestSubEvent, subsequentEvents := generateEvents(currentEvent, 1)

	events = append(events, subsequentEvents...)
	events = append(events, currentEvent)

	detectionService := NewLoginDetectionService(&MockEventRepository{userEvents: events}, nil, nil, 500)
	relatedEvent, err := detectionService.findRelatedEvents(currentEvent)
	req.NoError(err)
	req.Equal(currentEvent, relatedEvent.CurrentEvent)
	req.Equal(closestPreviousEvent, relatedEvent.PreviousEvent)
	req.Equal(closestSubEvent, relatedEvent.SubsequentEvent)
}

func newEvent(eventInfo models.EventInfo) *models.Event {
	currentEvent, err := models.NewEvent(eventInfo)
	if err != nil {
		log.Panicf(support.Fatal, err)
	}

	return currentEvent
}

func generateEvents(currentEvent *models.Event, increment int) (*models.Event, []*models.Event) {
	closestEvent := newEvent(models.EventInfo{
		UUID:      uuid.New().String(),
		Username:  "john",
		Timestamp: test.AddTime(currentEvent.ToEventInfo().Timestamp, increment, time.Hour),
		IP:        "1.0.0.0",
	})

	var events []*models.Event
	events = append(events, closestEvent)

	for i := 2; i < 4; i++ {
		newEvent := newEvent(models.EventInfo{
			UUID:      uuid.New().String(),
			Username:  "john",
			Timestamp: test.AddTime(currentEvent.ToEventInfo().Timestamp, increment*i, time.Hour),
			IP:        "1.0.0.0",
		})

		events = append(events, newEvent)
	}

	return closestEvent, events
}
