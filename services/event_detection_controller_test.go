package services

import (
	"github.com/frankiennamdi/detection-api/support"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEventDetectionHandler(t *testing.T) {
	testSetup := SetUp()
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
		DetectionService: testSetup.ServerContext().detectionService,
	}
	requestRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(detectionController.EventDetectionHandler)
	handler.ServeHTTP(requestRecorder, request)
	req.Equal(http.StatusOK, requestRecorder.Code)
	req.Equal(`{"currentGeo":{"lat":30.5334,"lon":-95.4559,"radius":1000}}`,
		requestRecorder.Body.String())
}

func TestEventDetectionHandler_With_Subsequent_Event(t *testing.T) {
	testSetup := SetUp()
	defer testSetup.CleanUp()

	detectionController := EventDetectionController{
		DetectionService: testSetup.ServerContext().detectionService,
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
	req.True(AssertThatJSONEqual(`{"currentGeo":{"lat":30.5334,"lon":-95.4559,"radius":1000}}`,
		requestRecorder.Body.String()))

	log.Printf(support.Info, requestRecorder.Body.String())

	requestRecorder = newRecordedRequest(detectionController, newPostRequest(`{
		"username": "bob",
		"unix_timestamp": 1514851200,
		"event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e43",
		"ip_address": "91.207.175.104"
	}`))

	req.Equal(http.StatusOK, requestRecorder.Code)
	log.Printf(support.Info, requestRecorder.Body.String())

	expectedResult := `{
						"currentGeo":{"lat":34.0549,"lon":-118.2578,"radius":200},
						"travelToCurrentGeoSuspicious":false,
						"precedingIpAccess":{
												"ip":"206.81.252.6",
												"speed":56,"lat":30.5334,
												"lon":-95.4559,
												"radius":1000,
												"timestamp":1514764800}
						}`

	req.True(AssertThatJSONEqual(expectedResult, requestRecorder.Body.String()))
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
