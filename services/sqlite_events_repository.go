package services

import (
	"database/sql"
	"github.com/frankiennamdi/detection-api/db"
	"github.com/frankiennamdi/detection-api/models"
)

type SqLiteEventsRepository struct {
	sqLiteDb *db.SqLiteDb
}

func NewSQLLiteEventsRepository(sqLiteDb *db.SqLiteDb) *SqLiteEventsRepository {
	return &SqLiteEventsRepository{sqLiteDb: sqLiteDb}
}

func (eventRepository *SqLiteEventsRepository) FindByUsername(username string) ([]*models.Event, error) {
	var events []*models.Event

	fnxErr := eventRepository.sqLiteDb.WithSqLiteDbContext(func(context *db.SqLiteDbContext) (err error) {
		events, err = eventRepository.findUser(username, context)

		if err != nil {
			return err
		}

		return nil
	})

	return events, fnxErr
}

func (eventRepository *SqLiteEventsRepository) Insert(events []*models.Event) ([]sql.Result, error) {
	var results []sql.Result

	fnxErr := eventRepository.sqLiteDb.WithSqLiteDbContext(func(context *db.SqLiteDbContext) (err error) {
		results, err = eventRepository.insertEvents(events, context)
		if err != nil {
			return err
		}
		return nil
	})

	if fnxErr != nil {
		return nil, fnxErr
	}

	return results, nil
}

func (eventRepository *SqLiteEventsRepository) findUser(username string,
	context *db.SqLiteDbContext) (events []*models.Event, err error) {
	events = make([]*models.Event, 0)
	stmt, err := context.Database().Prepare("SELECT * FROM events WHERE username = ? ORDER BY " +
		"timestamp ASC")

	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	rows, err := stmt.Query(username)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var eventInfo models.EventInfo
		err = rows.Scan(&eventInfo.UUID, &eventInfo.Username, &eventInfo.Timestamp, &eventInfo.IP)

		if err != nil {
			return nil, err
		}

		event, err := models.NewEvent(eventInfo)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return events, nil
}

func (eventRepository *SqLiteEventsRepository) insertEvents(events []*models.Event,
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
