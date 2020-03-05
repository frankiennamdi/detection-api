package models

type GeoPoint struct {
	Latitude       float64 `json:"lat"`
	Longitude      float64 `json:"lon"`
	AccuracyRadius uint16  `json:"radius"`
}

type GeoDistance struct {
	Miles float64
	Km    float64
}

type RelatedEventInfo struct {
	CurrentEvent    *Event
	PreviousEvent   *Event
	SubsequentEvent *Event
}

type EventGeoInfo struct {
	EventInfo *EventInfo
	GeoPoint  *GeoPoint
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
