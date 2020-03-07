package repository

import "github.com/frankiennamdi/detection-api/models"

type RelatedEventsFilter struct {
	currEvent     *models.Event
	relatedEvents *models.RelatedEventInfo
}

func NewRelatedEventsFilter(currEvent *models.Event) *RelatedEventsFilter {
	return &RelatedEventsFilter{currEvent: currEvent,
		relatedEvents: &models.RelatedEventInfo{
			CurrentEvent: currEvent,
		}}
}

func (filter *RelatedEventsFilter) Filter(event *models.Event) {
	currEventInfo := filter.currEvent.ToEventInfo()
	eventInfo := event.ToEventInfo()

	if eventInfo.UUID != currEventInfo.UUID {
		if eventInfo.Timestamp < currEventInfo.Timestamp {
			// earlier events
			if filter.relatedEvents.PreviousEvent == nil {
				filter.relatedEvents.PreviousEvent = event
			} else if filter.relatedEvents.PreviousEvent.ToEventInfo().Timestamp < eventInfo.Timestamp {
				filter.relatedEvents.PreviousEvent = event
			}
		} else if eventInfo.Timestamp > currEventInfo.Timestamp {
			// later events
			if filter.relatedEvents.SubsequentEvent == nil {
				filter.relatedEvents.SubsequentEvent = event
			} else if filter.relatedEvents.SubsequentEvent.ToEventInfo().Timestamp > eventInfo.Timestamp {
				filter.relatedEvents.SubsequentEvent = event
			}
		}
	}
}

func (filter *RelatedEventsFilter) GetRelatedEvents() *models.RelatedEventInfo {
	return filter.relatedEvents
}
