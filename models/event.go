package models

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/google/uuid"
)

// immutable event
type Event struct {
	uuid      string
	username  string
	timestamp int64
	ip        string
}

// mutable event info
type EventInfo struct {
	UUID      string `json:"event_uuid"`
	Username  string `json:"username"`
	Timestamp int64  `json:"unix_timestamp"`
	IP        string `json:"ip_address"`
}

type ValidationError struct {
	value    string
	argument string
}

func (event *Event) UnmarshalJSON(data []byte) error {
	info := &EventInfo{}

	if err := json.Unmarshal(data, &info); err != nil {
		return err
	}

	uEvent, validationErr := NewEvent(*info)
	if validationErr != nil {
		return validationErr
	}

	*event = *uEvent

	return nil
}

func (event *Event) ToEventInfo() EventInfo {
	return EventInfo{
		event.uuid,
		event.username,
		event.timestamp,
		event.ip,
	}
}

func (event *Event) MarshalJSON() ([]byte, error) {
	return json.Marshal(event.ToEventInfo())
}

func (err ValidationError) Error() string {
	return fmt.Sprintf("value: %s is invalid for argument: %s", err.value, err.argument)
}

func NewValidationError(value, argument string) *ValidationError {
	return &ValidationError{value: value, argument: argument}
}

func NewEvent(eventInfo EventInfo) (*Event, error) {
	if !IsValidUUID(eventInfo.UUID) {
		return nil, NewValidationError(eventInfo.UUID, "UUID")
	}

	if !IsIpv4Net(eventInfo.IP) {
		return nil, NewValidationError(eventInfo.IP, "IP")
	}

	return &Event{
		uuid:      eventInfo.UUID,
		username:  eventInfo.Username,
		timestamp: eventInfo.Timestamp,
		ip:        eventInfo.IP,
	}, nil
}

func EventFromJSON(jsonStr string) (*Event, error) {
	event := Event{}
	err := json.Unmarshal([]byte(jsonStr), &event)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func IsIpv4Net(ip string) bool {
	return net.ParseIP(ip) != nil
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
