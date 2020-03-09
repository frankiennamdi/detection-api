package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

var IPValidationTestsCases = []struct {
	IP             string
	expectedResult bool
}{
	{IP: "", expectedResult: false},

	{IP: "1.2.0.0", expectedResult: true},
}

var UUIDValidationTestsCases = []struct {
	UUID           string
	expectedResult bool
}{
	{UUID: "", expectedResult: false},
	{UUID: "85ad929a-db03-4bf4-9541-8f728fa12e42", expectedResult: true},
}

var newEventTestsCases = []struct {
	eventInfo     EventInfo
	expectedError error
}{
	{eventInfo: EventInfo{
		UUID:      "85ad929a-db03-4bf4-9541-8f728fa12e42",
		Username:  "john",
		Timestamp: 1514764800,
		IP:        "1.0.0.0",
	}, expectedError: nil},
	{eventInfo: EventInfo{
		UUID:      "85ad929a",
		Username:  "john",
		Timestamp: 1514764800,
		IP:        "1.0.0.0",
	}, expectedError: NewValidationError("85ad929a", "UUID")},
}

var newEventFromJSONTestCases = []struct {
	eventInfoJSON string
	expectedError error
}{
	{
		eventInfoJSON: `{
			"event_uuid":"85ad929a-db03-4bf4-9541-8f728fa12e42", 
			"username":"john",
			"unix_timestamp":1514764800, 
			"ip_address":"1.0.0.0"
		}`,
		expectedError: nil,
	},
}

func TestNewEvent(t *testing.T) {
	for _, input := range newEventTestsCases {
		req := require.New(t)
		event, err := NewEvent(input.eventInfo)
		req.Equal(input.expectedError, err)

		if err != nil {
			req.Nil(event)
		} else {
			req.NotNil(event)
			req.Equal(input.eventInfo.IP, event.ip)
			req.Equal(input.eventInfo.UUID, event.uuid)
			req.Equal(input.eventInfo.Username, event.username)
		}
	}
}

func TestNewEventFromJson(t *testing.T) {
	for _, input := range newEventFromJSONTestCases {
		req := require.New(t)
		event, err := EventFromJSON(input.eventInfoJSON)
		req.Equal(input.expectedError, err)

		if err != nil {
			req.Nil(event)
		} else {
			req.NotNil(event)
			eventInfo := EventInfo{}
			marshalErr := json.Unmarshal([]byte(input.eventInfoJSON), &eventInfo)
			req.NoError(marshalErr)
			req.Equal(eventInfo.IP, event.ip)
			req.Equal(eventInfo.UUID, event.uuid)
			req.Equal(eventInfo.Username, event.username)
		}
	}
}

func TestIsIpv4Net(t *testing.T) {
	for _, input := range IPValidationTestsCases {
		req := require.New(t)
		req.Equal(IsIpv4Net(input.IP), input.expectedResult)
	}
}

func TestIsValidUUID(t *testing.T) {
	for _, input := range UUIDValidationTestsCases {
		req := require.New(t)
		req.Equal(IsValidUUID(input.UUID), input.expectedResult)
	}
}
