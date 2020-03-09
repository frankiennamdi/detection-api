package repository

import (
	"database/sql"
	"github.com/frankiennamdi/detection-api/core"

	"github.com/frankiennamdi/detection-api/db"
	"github.com/frankiennamdi/detection-api/models"
)

// provides services for storing and retrieving events from SQLite database
type SqLiteEventsRepository struct {
	sqLiteDb *db.SqLiteDb
}

func NewSQLLiteEventsRepository(sqLiteDb *db.SqLiteDb) *SqLiteEventsRepository {
	return &SqLiteEventsRepository{sqLiteDb: sqLiteDb}
}

func (eventRepository SqLiteEventsRepository) InsertAndFindRelatedEvents(event *models.Event,
	filter core.EventFilter) error {
	fnxErr := eventRepository.sqLiteDb.WithSqLiteDbContext(func(context *db.SqLiteDbContext) (err error) {
		if _, err := eventRepository.insertEvents([]*models.Event{event}, context); err != nil {
			return err
		}

		if err := eventRepository.findAndFilter(event, filter, context); err != nil {
			return err
		}

		return nil
	}, "mode=rw")

	return fnxErr
}

func (eventRepository SqLiteEventsRepository) FindRelatedEvents(event *models.Event,
	filter core.EventFilter) error {
	fnxErr := eventRepository.sqLiteDb.WithSqLiteDbContext(func(context *db.SqLiteDbContext) (err error) {
		if err = eventRepository.findAndFilter(event, filter, context); err != nil {
			return err
		}

		return nil
	}, "mode=rw")

	return fnxErr
}

func (eventRepository SqLiteEventsRepository) InsertEvents(events []*models.Event) ([]sql.Result, error) {
	var results []sql.Result

	fnxErr := eventRepository.sqLiteDb.WithSqLiteDbContext(func(context *db.SqLiteDbContext) (err error) {
		results, err = eventRepository.insertEvents(events, context)
		if err != nil {
			return err
		}
		return nil
	}, "mode=rw")

	if fnxErr != nil {
		return nil, fnxErr
	}

	return results, nil
}

func (eventRepository SqLiteEventsRepository) findAndFilter(event *models.Event,
	filter core.EventFilter,
	context *db.SqLiteDbContext) (err error) {
	stmt, err := context.Database().Prepare("SELECT * FROM events WHERE username = ? ORDER BY " +
		"timestamp ASC")

	if err != nil {
		return err
	}

	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	rows, err := stmt.Query(event.ToEventInfo().Username)

	if err != nil {
		return err
	}

	for rows.Next() {
		var eventInfo models.EventInfo
		err = rows.Scan(&eventInfo.UUID, &eventInfo.Username, &eventInfo.Timestamp, &eventInfo.IP)

		if err != nil {
			return err
		}

		event, err := models.NewEvent(eventInfo)
		if err != nil {
			return err
		}

		filter.Filter(event)
	}

	err = rows.Err()

	if err != nil {
		return err
	}

	return nil
}

func (eventRepository SqLiteEventsRepository) insertEvents(events []*models.Event,
	context *db.SqLiteDbContext) ([]sql.Result, error) {
	var results []sql.Result

	transactionErr := context.WithTransaction(func(tx *sql.Tx) (err error) {
		stmt, err := tx.Prepare("INSERT OR IGNORE INTO events(uuid, username, timestamp, ip) " +
			"VALUES(?, ?, ?, ?)")

		if err != nil {
			return err
		}

		for _, event := range events {
			eventInfo := event.ToEventInfo()
			result, stmtErr := stmt.Exec(eventInfo.UUID, eventInfo.Username, eventInfo.Timestamp, eventInfo.IP)
			if stmtErr != nil {
				return stmtErr
			}
			results = append(results, result)
		}

		defer func() {
			if closeErr := stmt.Close(); closeErr != nil {
				err = closeErr
			}
		}()

		return nil
	})

	if transactionErr != nil {
		return nil, transactionErr
	}

	return results, nil
}
