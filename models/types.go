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

/*
{
  "currentGeo": {
    "lat": 39.1702,
    "lon": -76.8538,
    "radius": 20
  },
  "travelToCurrentGeoSuspicious": true,
  "travelFromCurrentGeoSuspicious": false,
  "precedingIpAccess": {
    "ip": "24.242.71.20",
    "speed": 55,
    "lat": 30.3764,
    "lon": -97.7078,
    "radius": 5,
    "timestamp": 1514764800
  },
  "subsequentIpAccess": {
    "ip": "91.207.175.104",
    "speed": 27600,
    "lat": 34.0494,
    "lon": -118.2641,
    "radius": 200,
    "timestamp": 1514851200
  }
}
*/

type RelatedAccessInfo struct {
	IP             string  `json:"ip"`
	Speed          float64 `json:"speed"`
	Latitude       float64 `json:"lat"`
	Longitude      float64 `json:"lon"`
	AccuracyRadius uint16  `json:"radius"`
	Timestamp      int64   `json:"timestamp"`
}
