package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/frankiennamdi/detection-api/test"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/frankiennamdi/detection-api/models"
	"github.com/frankiennamdi/detection-api/support"
	"github.com/stretchr/testify/require"
)

func TestEventDetectionHandler(t *testing.T) {
	testSetup := test.SetUp()
	defer testSetup.CleanUp()

	body := `{
		"username": "bob",
		"unix_timestamp": 1514764800,
		"event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e42",
		"ip_address": "206.81.252.6"
	}`
	request, err := http.NewRequest(http.MethodPost, "/api/events", strings.NewReader(body))
	req := require.New(t)
	req.NoError(err)

	detectionController := EventDetectionController{
		DetectionService: NewServiceContext(testSetup.AppServerContext()).DetectionService(),
	}
	requestRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(detectionController.EventDetectionHandler)
	handler.ServeHTTP(requestRecorder, request)
	req.Equal(http.StatusOK, requestRecorder.Code)

	expected := &models.SuspiciousTravelResult{}

	expected.CurrentGeo = &models.GeoPoint{
		Latitude:       30.5334,
		Longitude:      -95.4559,
		AccuracyRadius: 1000,
	}
	assertThatBodyResultEquals(expected, requestRecorder.Body, req)
}

func TestEventDetectionHandler_With_Subsequent_Event(t *testing.T) {
	testSetup := test.SetUp()
	defer testSetup.CleanUp()

	detectionController := EventDetectionController{
		DetectionService: NewServiceContext(testSetup.AppServerContext()).DetectionService(),
	}

	req := require.New(t)

	var requestRecorder *httptest.ResponseRecorder

	requestRecorder = newRecordedRequest(detectionController, newPostRequest(`{
		"username": "bob",
		"unix_timestamp": 1514764800,
		"event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e42",
		"ip_address": "206.81.252.6"
	}`))

	req.Equal(http.StatusOK, requestRecorder.Code)
	assertThatBodyResultEquals(&models.SuspiciousTravelResult{
		CurrentGeo: &models.GeoPoint{
			Latitude:       30.5334,
			Longitude:      -95.4559,
			AccuracyRadius: 1000,
		},
	}, requestRecorder.Body, req)

	requestRecorder = newRecordedRequest(detectionController, newPostRequest(`{
		"username": "bob",
		"unix_timestamp": 1514851200,
		"event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e43",
		"ip_address": "91.207.175.104"
	}`))

	req.Equal(http.StatusOK, requestRecorder.Code)
	assertThatBodyResultEquals(&models.SuspiciousTravelResult{
		CurrentGeo: &models.GeoPoint{
			Latitude:       34.0549,
			Longitude:      -118.2578,
			AccuracyRadius: 200,
		},
		TravelToCurrentGeoSuspicious: boolean(false),
		PrecedingIPAccess: &models.RelatedAccessInfo{
			IP:             "206.81.252.6",
			Speed:          56,
			Latitude:       30.5334,
			Longitude:      -95.4559,
			AccuracyRadius: 1000,
			Timestamp:      1514764800,
		},
	}, requestRecorder.Body, req)
}

func TestEventDetectionHandler_With_Previous_And_Subsequent_Event(t *testing.T) {
	testSetup := test.SetUp()
	defer testSetup.CleanUp()

	detectionController := EventDetectionController{
		DetectionService: NewServiceContext(testSetup.AppServerContext()).DetectionService(),
	}

	req := require.New(t)

	var requestRecorder *httptest.ResponseRecorder

	requestRecorder = newRecordedRequest(detectionController, newPostRequest(`{
 		"username": "bob",
 		"unix_timestamp": 1514851200,
 		"event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e43",
		"ip_address": "24.242.71.20"
	}`))

	req.Equal(http.StatusOK, requestRecorder.Code)
	assertThatBodyResultEquals(&models.SuspiciousTravelResult{
		CurrentGeo: &models.GeoPoint{
			Latitude:       30.3773,
			Longitude:      -97.71,
			AccuracyRadius: 5,
		},
	}, requestRecorder.Body, req)

	requestRecorder = newRecordedRequest(detectionController, newPostRequest(`{
		"username": "bob",
		"unix_timestamp": 1514761200,
		"event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e49",
 		"ip_address": "91.207.175.104"
 	}`))

	req.Equal(http.StatusOK, requestRecorder.Code)

	assertThatBodyResultEquals(&models.SuspiciousTravelResult{
		CurrentGeo: &models.GeoPoint{
			Latitude:       34.0549,
			Longitude:      -118.2578,
			AccuracyRadius: 200,
		},
		TravelFromCurrentGeoSuspicious: boolean(false),
		SubsequentIPAccess: &models.RelatedAccessInfo{
			IP:             "24.242.71.20",
			Speed:          49,
			Latitude:       30.3773,
			Longitude:      -97.71,
			AccuracyRadius: 5,
			Timestamp:      1514851200,
		},
	}, requestRecorder.Body, req)

	requestRecorder = newRecordedRequest(detectionController, newPostRequest(`{
		"username": "bob",
		"unix_timestamp": 1514764800,
		"event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e48",
		"ip_address": "206.81.252.6"
	}`))

	req.Equal(http.StatusOK, requestRecorder.Code)
	assertThatBodyResultEquals(&models.SuspiciousTravelResult{
		CurrentGeo: &models.GeoPoint{
			Latitude:       30.5334,
			Longitude:      -95.4559,
			AccuracyRadius: 1000,
		},
		TravelToCurrentGeoSuspicious:   boolean(true),
		TravelFromCurrentGeoSuspicious: boolean(false),
		PrecedingIPAccess: &models.RelatedAccessInfo{
			IP:             "91.207.175.104",
			Speed:          1351,
			Latitude:       34.0549,
			Longitude:      -118.2578,
			AccuracyRadius: 200,
			Timestamp:      1514761200,
		},
		SubsequentIPAccess: &models.RelatedAccessInfo{
			IP:             "24.242.71.20",
			Speed:          6,
			Latitude:       30.3773,
			Longitude:      -97.71,
			AccuracyRadius: 5,
			Timestamp:      1514851200,
		},
	}, requestRecorder.Body, req)
}

func newRecordedRequest(detectionController EventDetectionController,
	request *http.Request) *httptest.ResponseRecorder {
	requestRecorder := httptest.NewRecorder()

	handler := http.HandlerFunc(detectionController.EventDetectionHandler)
	handler.ServeHTTP(requestRecorder, request)

	return requestRecorder
}

func newPostRequest(body string) *http.Request {
	request, err := http.NewRequest(http.MethodPost, "/api/events", strings.NewReader(body))
	if err != nil {
		log.Panicf(support.Fatal, err)
	}

	return request
}

func assertThatBodyResultEquals(expected *models.SuspiciousTravelResult, actual *bytes.Buffer,
	assert *require.Assertions) {
	result, err := unmarshalToSuspiciousTravelResult(fmt.Sprint(actual))
	assert.NoError(err)
	assert.Equal(expected.TravelFromCurrentGeoSuspicious, result.TravelFromCurrentGeoSuspicious)
	assert.Equal(expected.TravelFromCurrentGeoSuspicious, result.TravelFromCurrentGeoSuspicious)
	assert.Equal(expected.CurrentGeo, result.CurrentGeo)
	assert.Equal(expected.PrecedingIPAccess, result.PrecedingIPAccess)
	assert.Equal(expected.SubsequentIPAccess, result.SubsequentIPAccess)
}

func boolean(value bool) *bool {
	return &value
}

func unmarshalToSuspiciousTravelResult(body string) (*models.SuspiciousTravelResult, error) {
	var result *models.SuspiciousTravelResult
	err := json.Unmarshal([]byte(body), &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}
