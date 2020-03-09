package models

type GeoDistance struct {
	miles float64
	km    float64
}

func NewGeoDistance(km float64, miles float64) *GeoDistance {
	return &GeoDistance{
		miles: miles,
		km:    km,
	}
}

func (geoDistance GeoDistance) Miles() float64 {
	return geoDistance.miles
}

func (geoDistance GeoDistance) Km() float64 {
	return geoDistance.km
}

type EventGeoInfo struct {
	eventInfo *EventInfo
	geoPoint  *GeoPoint
}

func NewEventGeoInfo(eventInfo *EventInfo, geoPoint *GeoPoint) *EventGeoInfo {
	return &EventGeoInfo{
		eventInfo: eventInfo,
		geoPoint:  geoPoint,
	}
}

func (eventGeoInfo EventGeoInfo) GeoPoint() *GeoPoint {
	return eventGeoInfo.geoPoint
}

func (eventGeoInfo EventGeoInfo) EventInfo() *EventInfo {
	return eventGeoInfo.eventInfo
}

type RelatedEventInfo struct {
	CurrentEvent    *Event
	PreviousEvent   *Event
	SubsequentEvent *Event
}

type GeoPoint struct {
	Latitude       float64 `json:"lat"`
	Longitude      float64 `json:"lon"`
	AccuracyRadius uint16  `json:"radius"`
}

type SuspiciousTravelResult struct {
	CurrentGeo                     *GeoPoint          `json:"currentGeo,omitempty"`
	TravelToCurrentGeoSuspicious   *bool              `json:"travelToCurrentGeoSuspicious,omitempty"`
	TravelFromCurrentGeoSuspicious *bool              `json:"travelFromCurrentGeoSuspicious,omitempty"`
	PrecedingIPAccess              *RelatedAccessInfo `json:"precedingIpAccess,omitempty"`
	SubsequentIPAccess             *RelatedAccessInfo `json:"subsequentIpAccess,omitempty"`
}

type RelatedAccessInfo struct {
	IP             string  `json:"ip"`
	Speed          float64 `json:"speed"`
	Latitude       float64 `json:"lat"`
	Longitude      float64 `json:"lon"`
	AccuracyRadius uint16  `json:"radius"`
	Timestamp      int64   `json:"timestamp"`
}
