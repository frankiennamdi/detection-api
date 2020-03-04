package services

import (
	"github.com/frankiennamdi/detection-api/models"
	"github.com/frankiennamdi/detection-api/support"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"log"
	"sort"
	"testing"
	"time"
)

func TestInsertEvent(t *testing.T) {
	testSetup := SetUp()
	defer testSetup.CleanUp()

	req := require.New(t)
	eventRepository := NewSQLLiteEventsRepository(testSetup.ServerContext().EventDb())
	event := newTestEvent(models.EventInfo{
		UUID:      "85ad929a-db03-4bf4-9541-8f728fa12e42",
		Username:  "john",
		Timestamp: 1514764800,
		IP:        "1.0.0.0",
	})
	req.NotNil(event)
	result, insertErr := eventRepository.Insert([]*models.Event{event})
	req.NoError(insertErr)
	req.Len(result, 1)
	rows, err := result[0].RowsAffected()
	req.NoError(err)
	req.Equal(1, int(rows))
}

func TestInsertEvent_Maintain_Original_When_Uuid_Is_Duplicated(t *testing.T) {
	initialTime := int64(1514764800)
	testSetup := SetUp()

	defer testSetup.CleanUp()

	req := require.New(t)
	eventRepository := NewSQLLiteEventsRepository(testSetup.ServerContext().EventDb())
	_, insertErr := eventRepository.Insert([]*models.Event{newTestEvent(models.EventInfo{
		UUID:      "85ad929a-db03-4bf4-9541-8f728fa12e42",
		Username:  "john",
		Timestamp: initialTime,
		IP:        "1.0.0.0",
	}), newTestEvent(models.EventInfo{
		UUID:      "85ad929a-db03-4bf4-9541-8f728fa12e42",
		Username:  "john",
		Timestamp: AddTime(initialTime, 1, time.Hour),
		IP:        "1.0.0.0",
	})})
	req.NoError(insertErr)

	events, err := eventRepository.FindByUsername("john")
	req.NoError(err)
	req.Len(events, 1)
}

func TestInsertEvent_Maintain_Original_When_User_Timestamp_Is_Duplicated(t *testing.T) {
	initialTime := int64(1514764800)
	testSetup := SetUp()

	defer testSetup.CleanUp()

	req := require.New(t)
	eventRepository := NewSQLLiteEventsRepository(testSetup.ServerContext().EventDb())
	_, insertErr := eventRepository.Insert([]*models.Event{newTestEvent(models.EventInfo{
		UUID:      uuid.New().String(),
		Username:  "john",
		Timestamp: initialTime,
		IP:        "1.0.0.0",
	}), newTestEvent(models.EventInfo{
		UUID:      uuid.New().String(),
		Username:  "john",
		Timestamp: initialTime,
		IP:        "1.0.0.0",
	})})
	req.NoError(insertErr)

	events, err := eventRepository.FindByUsername("john")
	req.NoError(err)
	req.Len(events, 1)
}

func TestInsertAndQueryEvent(t *testing.T) {
	testSetup := SetUp()
	defer testSetup.CleanUp()

	req := require.New(t)
	eventRepository := NewSQLLiteEventsRepository(testSetup.ServerContext().EventDb())
	eventInfo := models.EventInfo{
		UUID:      "85ad929a-db03-4bf4-9541-8f728fa12e42",
		Username:  "john",
		Timestamp: 1514764800,
		IP:        "1.0.0.0",
	}

	event := newTestEvent(eventInfo)
	req.NotNil(event)
	result, insertErr := eventRepository.Insert([]*models.Event{event})
	req.NoError(insertErr)
	req.Len(result, 1)
	rows, err := result[0].RowsAffected()
	req.NoError(err)
	req.Equal(1, int(rows))

	events, err := eventRepository.FindByUsername("john")
	req.NoError(err)

	actualEventInfo := events[0].ToEventInfo()
	req.Equal(eventInfo, actualEventInfo)
}

func TestInsertAndQueryEvent_That_Events_Are_Returned_In_Order(t *testing.T) {
	initialTime := int64(1514764800)
	testSetup := SetUp()

	defer testSetup.CleanUp()

	req := require.New(t)
	eventRepository := NewSQLLiteEventsRepository(testSetup.ServerContext().EventDb())
	events := []*models.Event{newTestEvent(models.EventInfo{
		UUID:      uuid.New().String(),
		Username:  "john",
		Timestamp: AddTime(initialTime, 10, time.Hour),
		IP:        "1.0.0.0",
	}), newTestEvent(models.EventInfo{
		UUID:      uuid.New().String(),
		Username:  "john",
		Timestamp: AddTime(initialTime, -10, time.Hour),
		IP:        "1.0.0.0",
	}), newTestEvent(models.EventInfo{
		UUID:      uuid.New().String(),
		Username:  "john",
		Timestamp: initialTime,
		IP:        "1.0.0.0",
	})}

	result, insertErr := eventRepository.Insert(events)
	req.NoError(insertErr)
	req.Len(result, len(events))

	actualEvents, err := eventRepository.FindByUsername("john")
	req.NoError(err)

	var actualTimeStamps []int64

	for _, event := range actualEvents {
		if event != nil {
			actualTimeStamps = append(actualTimeStamps, event.ToEventInfo().Timestamp)
		}
	}

	var expectedTimeStamps []int64

	for _, event := range events {
		expectedTimeStamps = append(expectedTimeStamps, event.ToEventInfo().Timestamp)
	}

	sort.Slice(expectedTimeStamps, func(i, j int) bool { return expectedTimeStamps[i] < expectedTimeStamps[j] })
	req.Equal(expectedTimeStamps, actualTimeStamps)
}

func newTestEvent(eventInfo models.EventInfo) *models.Event {
	event, err := models.NewEvent(eventInfo)
	if err != nil {
		log.Fatalf(support.Fatal, err)
	}

	return event
}
