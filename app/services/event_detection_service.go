package services

import (
	"fmt"
	"github.com/frankiennamdi/detection-api/core"
	"github.com/frankiennamdi/detection-api/repository"
	"github.com/frankiennamdi/detection-api/support"
	"net"

	"github.com/frankiennamdi/detection-api/models"
)

// service for detecting event characteristics relative to other events
type EventDetectionService struct {
	eventRepository     core.EventRepository
	ipGeoInfoRepository core.IPGeoInfoRepository
	calculatorService   core.CalculatorService
	suspiciousSpeed     float64
}

func NewDetectionService(
	eventRepository core.EventRepository,
	ipGeoInfoRepository core.IPGeoInfoRepository,
	calculatorService core.CalculatorService,
	suspiciousSpeed float64) *EventDetectionService {
	return &EventDetectionService{
		eventRepository:     eventRepository,
		ipGeoInfoRepository: ipGeoInfoRepository,
		calculatorService:   calculatorService,
		suspiciousSpeed:     suspiciousSpeed,
	}
}

func (service EventDetectionService) ProcessEvent(currEvent *models.Event) (*models.SuspiciousTravelResult, error) {
	if currEvent == nil {
		return nil, support.NewIllegalArgumentError("currEvent cannot be nil")
	}

	relatedEventInfo, err := service.findRelatedEvents(currEvent)
	if err != nil {
		return nil, err
	}

	suspiciousTravelResult, err := service.findSuspiciousTravel(relatedEventInfo)

	if err != nil {
		return nil, err
	}

	return suspiciousTravelResult, err
}

func (service EventDetectionService) findSuspiciousTravel(
	relatedEventInfo *models.RelatedEventInfo) (*models.SuspiciousTravelResult, error) {
	result := &models.SuspiciousTravelResult{}

	currEventInfo := relatedEventInfo.CurrentEvent.ToEventInfo()
	currEventGeoInfo, err := service.ipGeoInfoRepository.FindGeoPoint(net.ParseIP(currEventInfo.IP))

	if err != nil {
		return nil, err
	}

	if currEventGeoInfo == nil {
		return nil, fmt.Errorf("cannot find geo information for event: %+v", currEventInfo)
	}

	result.CurrentGeo = currEventGeoInfo

	if relatedEventInfo.PreviousEvent != nil {
		preEventInfo := relatedEventInfo.PreviousEvent.ToEventInfo()
		preEventGeoInfo, err := service.ipGeoInfoRepository.FindGeoPoint(net.ParseIP(preEventInfo.IP))

		if err != nil {
			return nil, err
		}

		if preEventGeoInfo != nil {
			travelToCurrentGeoSpeed, err := service.calculatorService.SpeedToTravelDistanceInMPH(
				models.NewEventGeoInfo(&currEventInfo, currEventGeoInfo),
				models.NewEventGeoInfo(&preEventInfo, preEventGeoInfo))

			if err != nil {
				return nil, err
			}

			value := *travelToCurrentGeoSpeed >= service.suspiciousSpeed
			result.TravelToCurrentGeoSuspicious = &value
			result.PrecedingIPAccess = &models.RelatedAccessInfo{
				IP:             preEventInfo.IP,
				Speed:          *travelToCurrentGeoSpeed,
				Latitude:       preEventGeoInfo.Latitude,
				Longitude:      preEventGeoInfo.Longitude,
				AccuracyRadius: preEventGeoInfo.AccuracyRadius,
				Timestamp:      preEventInfo.Timestamp,
			}
		}
	}

	if relatedEventInfo.SubsequentEvent != nil {
		subEventInfo := relatedEventInfo.SubsequentEvent.ToEventInfo()
		subEventGeoInfo, err := service.ipGeoInfoRepository.FindGeoPoint(net.ParseIP(subEventInfo.IP))

		if err != nil {
			return nil, err
		}

		if subEventGeoInfo != nil {
			travelFromCurrentGeoSpeed, err := service.calculatorService.SpeedToTravelDistanceInMPH(
				models.NewEventGeoInfo(&currEventInfo, currEventGeoInfo),
				models.NewEventGeoInfo(&subEventInfo, subEventGeoInfo))

			if err != nil {
				return nil, err
			}

			value := *travelFromCurrentGeoSpeed >= service.suspiciousSpeed
			result.TravelFromCurrentGeoSuspicious = &value
			result.SubsequentIPAccess = &models.RelatedAccessInfo{
				IP:             subEventInfo.IP,
				Speed:          *travelFromCurrentGeoSpeed,
				Latitude:       subEventGeoInfo.Latitude,
				Longitude:      subEventGeoInfo.Longitude,
				AccuracyRadius: subEventGeoInfo.AccuracyRadius,
				Timestamp:      subEventInfo.Timestamp,
			}
		}
	}

	return result, nil
}

func (service EventDetectionService) findRelatedEvents(currEvent *models.Event) (*models.RelatedEventInfo, error) {
	filter := repository.NewRelatedEventsFilter(currEvent)
	err := service.eventRepository.InsertAndFindRelatedEvents(currEvent, filter)

	if err != nil {
		return nil, err
	}

	return filter.GetRelatedEvents(), nil
}
