package services

import (
	"fmt"
	"github.com/frankiennamdi/detection-api/core"
	"github.com/frankiennamdi/detection-api/repository"
	"net"

	"github.com/frankiennamdi/detection-api/models"
)

type LoginDetectionService struct {
	eventRepository     core.EventRepository
	ipGeoInfoRepository core.IPGeoInfoRepository
	calculatorService   core.CalculatorService
	suspiciousSpeed     float64
}

func NewLoginDetectionService(
	eventRepository core.EventRepository,
	ipGeoInfoRepository core.IPGeoInfoRepository,
	calculatorService core.CalculatorService,
	suspiciousSpeed float64) *LoginDetectionService {
	return &LoginDetectionService{
		eventRepository:     eventRepository,
		ipGeoInfoRepository: ipGeoInfoRepository,
		calculatorService:   calculatorService,
		suspiciousSpeed:     suspiciousSpeed,
	}
}

func (service LoginDetectionService) ProcessEvent(currEvent *models.Event) (*models.SuspiciousTravelResult, error) {
	if _, err := service.eventRepository.InsertEvents([]*models.Event{currEvent}); err != nil {
		return nil, err
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

func (service LoginDetectionService) findSuspiciousTravel(
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
			travelToCurrentGeoSpeed := service.calculatorService.SpeedToTravelDistanceInMPH(models.EventGeoInfo{
				EventInfo: &currEventInfo,
				GeoPoint:  currEventGeoInfo,
			}, models.EventGeoInfo{
				EventInfo: &preEventInfo,
				GeoPoint:  preEventGeoInfo,
			})
			value := travelToCurrentGeoSpeed >= service.suspiciousSpeed
			result.TravelToCurrentGeoSuspicious = &value
			result.PrecedingIPAccess = &models.RelatedAccessInfo{
				IP:             preEventInfo.IP,
				Speed:          travelToCurrentGeoSpeed,
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
			travelFromCurrentGeoSpeed := service.calculatorService.SpeedToTravelDistanceInMPH(models.EventGeoInfo{
				EventInfo: &currEventInfo,
				GeoPoint:  currEventGeoInfo,
			}, models.EventGeoInfo{
				EventInfo: &subEventInfo,
				GeoPoint:  subEventGeoInfo,
			})
			value := travelFromCurrentGeoSpeed >= service.suspiciousSpeed
			result.TravelFromCurrentGeoSuspicious = &value
			result.SubsequentIPAccess = &models.RelatedAccessInfo{
				IP:             subEventInfo.IP,
				Speed:          travelFromCurrentGeoSpeed,
				Latitude:       subEventGeoInfo.Latitude,
				Longitude:      subEventGeoInfo.Longitude,
				AccuracyRadius: subEventGeoInfo.AccuracyRadius,
				Timestamp:      subEventInfo.Timestamp,
			}
		}
	}

	return result, nil
}

func (service LoginDetectionService) findRelatedEvents(currEvent *models.Event) (*models.RelatedEventInfo, error) {
	filter := repository.NewRelatedEventsFilter(currEvent)
	err := service.eventRepository.FindRelatedEvents(currEvent, filter)

	if err != nil {
		return nil, err
	}

	return filter.GetRelatedEvents(), nil
}
