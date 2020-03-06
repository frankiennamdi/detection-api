package repository

import (
	"github.com/frankiennamdi/detection-api/test"
	"log"
	"testing"
	"time"

	"github.com/frankiennamdi/detection-api/models"
	"github.com/frankiennamdi/detection-api/support"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestInsertEvent(t *testing.T) {
	testSetup := test.SetUp()
	defer testSetup.CleanUp()

	req := require.New(t)
	eventRepository := NewSQLLiteEventsRepository(testSetup.AppServerContext().EventDb())
	event := newTestEvent(models.EventInfo{
		UUID:      "85ad929a-db03-4bf4-9541-8f728fa12e42",
		Username:  "john",
		Timestamp: 1514764800,
		IP:        "1.0.0.0",
	})
	req.NotNil(event)
	result, insertErr := eventRepository.InsertEvents([]*models.Event{event})
	req.NoError(insertErr)
	req.Len(result, 1)
	rows, err := result[0].RowsAffected()
	req.NoError(err)
	req.Equal(1, int(rows))
}

func TestInsertEvent_Maintain_Original_When_Uuid_Is_Duplicated(t *testing.T) {
	initialTime := int64(1514764800)
	testSetup := test.SetUp()

	defer testSetup.CleanUp()

	req := require.New(t)
	eventRepository := NewSQLLiteEventsRepository(testSetup.AppServerContext().EventDb())
	events := []*models.Event{newTestEvent(models.EventInfo{
		UUID:      "85ad929a-db03-4bf4-9541-8f728fa12e42",
		Username:  "john",
		Timestamp: initialTime,
		IP:        "1.0.0.0",
	}), newTestEvent(models.EventInfo{
		UUID:      "85ad929a-db03-4bf4-9541-8f728fa12e42",
		Username:  "john",
		Timestamp: test.AddTime(initialTime, 1, time.Hour),
		IP:        "1.0.0.0",
	})}
	_, insertErr := eventRepository.InsertEvents(events)
	req.NoError(insertErr)

	filter := NewRelatedEventsFilter(events[0])
	err := eventRepository.FindRelatedEvents(events[0], filter)
	req.NoError(err)
	req.Equal(events[0].ToEventInfo(), filter.GetRelatedEvents().CurrentEvent.ToEventInfo())
}

func TestInsertEvent_Maintain_Original_When_User_Timestamp_Is_Duplicated(t *testing.T) {
	initialTime := int64(1514764800)
	testSetup := test.SetUp()

	defer testSetup.CleanUp()

	req := require.New(t)
	eventRepository := NewSQLLiteEventsRepository(testSetup.AppServerContext().EventDb())
	events := []*models.Event{newTestEvent(models.EventInfo{
		UUID:      uuid.New().String(),
		Username:  "john",
		Timestamp: initialTime,
		IP:        "1.0.0.0",
	}), newTestEvent(models.EventInfo{
		UUID:      uuid.New().String(),
		Username:  "john",
		Timestamp: initialTime,
		IP:        "1.0.0.0",
	})}
	_, insertErr := eventRepository.InsertEvents(events)
	req.NoError(insertErr)

	filter := NewRelatedEventsFilter(events[0])
	err := eventRepository.FindRelatedEvents(events[0], filter)

	req.NoError(err)
	req.Nil(filter.GetRelatedEvents().SubsequentEvent)
	req.Nil(filter.GetRelatedEvents().PreviousEvent)
	req.NotNil(filter.GetRelatedEvents().CurrentEvent)
	req.Equal(events[0], filter.GetRelatedEvents().CurrentEvent)
}

func TestInsertAndQueryEvent(t *testing.T) {
	testSetup := test.SetUp()
	defer testSetup.CleanUp()

	req := require.New(t)
	eventRepository := NewSQLLiteEventsRepository(testSetup.AppServerContext().EventDb())
	eventInfo := models.EventInfo{
		UUID:      "85ad929a-db03-4bf4-9541-8f728fa12e42",
		Username:  "john",
		Timestamp: 1514764800,
		IP:        "1.0.0.0",
	}

	event := newTestEvent(eventInfo)
	req.NotNil(event)
	result, insertErr := eventRepository.InsertEvents([]*models.Event{event})
	req.NoError(insertErr)
	req.Len(result, 1)
	rows, err := result[0].RowsAffected()
	req.NoError(err)
	req.Equal(1, int(rows))

	filter := NewRelatedEventsFilter(event)
	filterErr := eventRepository.FindRelatedEvents(event, filter)
	req.NoError(filterErr)

	actualEventInfo := filter.GetRelatedEvents().CurrentEvent.ToEventInfo()
	req.Equal(eventInfo, actualEventInfo)
}

func TestInsertAndQueryEvent_That_Events_Are_Returned_In_Order(t *testing.T) {
	initialTime := int64(1514764800)
	testSetup := test.SetUp()

	defer testSetup.CleanUp()

	req := require.New(t)
	eventRepository := NewSQLLiteEventsRepository(testSetup.AppServerContext().EventDb())
	events := []*models.Event{newTestEvent(models.EventInfo{
		UUID:      uuid.New().String(),
		Username:  "john",
		Timestamp: test.AddTime(initialTime, 10, time.Hour),
		IP:        "1.0.0.0",
	}), newTestEvent(models.EventInfo{
		UUID:      uuid.New().String(),
		Username:  "john",
		Timestamp: test.AddTime(initialTime, -10, time.Hour),
		IP:        "1.0.0.0",
	}), newTestEvent(models.EventInfo{
		UUID:      uuid.New().String(),
		Username:  "john",
		Timestamp: initialTime,
		IP:        "1.0.0.0",
	})}

	result, insertErr := eventRepository.InsertEvents(events)
	req.NoError(insertErr)
	req.Len(result, len(events))

	filter := NewRelatedEventsFilter(events[2])
	err := eventRepository.FindRelatedEvents(events[2], filter)
	req.NoError(err)
	req.Equal(events[1].ToEventInfo(), filter.GetRelatedEvents().PreviousEvent.ToEventInfo())
	req.Equal(events[0].ToEventInfo(), filter.GetRelatedEvents().SubsequentEvent.ToEventInfo())
	req.Equal(events[2].ToEventInfo(), filter.GetRelatedEvents().CurrentEvent.ToEventInfo())
}

func TestFindEventByUsername_In_Empty_Database(t *testing.T) {
	testSetup := test.SetUp()

	defer testSetup.CleanUp()

	req := require.New(t)
	eventRepository := NewSQLLiteEventsRepository(testSetup.AppServerContext().EventDb())
	event, err := models.NewEvent(models.EventInfo{
		UUID:      uuid.New().String(),
		Username:  "john",
		Timestamp: 1514764800,
		IP:        "1.0.0.0",
	})
	req.NoError(err)

	filter := NewRelatedEventsFilter(event)
	filterErr := eventRepository.FindRelatedEvents(event, filter)
	req.NoError(filterErr)
	req.Equal(event, filter.GetRelatedEvents().CurrentEvent)
}

func newTestEvent(eventInfo models.EventInfo) *models.Event {
	event, err := models.NewEvent(eventInfo)
	if err != nil {
		log.Panicf(support.Fatal, err)
	}

	return event
}
