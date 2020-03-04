package db

import (
	"database/sql"
	"fmt"
	appConfig "github.com/frankiennamdi/detection-api/config"
	"github.com/frankiennamdi/detection-api/support"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
)

type SqLiteDb struct {
	Config appConfig.AppConfig
}

type SqLiteDbContext struct {
	db     *sql.DB
	config appConfig.AppConfig
}

func (sqLiteDb *SqLiteDb) NewCreateContext() (*SqLiteDbContext, error) {
	return sqLiteDb.newContext("mode=rwc")
}

func (sqLiteDb *SqLiteDb) NewReadWriteContext() (*SqLiteDbContext, error) {
	return sqLiteDb.newContext("mode=rw")
}

func (sqLiteDb *SqLiteDb) newContext(parameters string) (*SqLiteDbContext, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?%s", sqLiteDb.Config.EventDb.File, parameters))
	if err != nil {
		return nil, err
	}

	return &SqLiteDbContext{db: db,
		config: sqLiteDb.Config,
	}, nil
}

func (sqLiteDbContext *SqLiteDbContext) Close() error {
	if err := sqLiteDbContext.db.Close(); err != nil {
		return err
	}

	return nil
}

type TransactionEnabled func(tx *sql.Tx) error

func (sqLiteDbContext *SqLiteDbContext) WithTransaction(fnx TransactionEnabled) error {
	tx, beginErr := sqLiteDbContext.db.Begin()
	if beginErr != nil {
		return beginErr
	}

	fnxErr := fnx(tx)

	if fnxErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}

		return fnxErr
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		return commitErr
	}

	return nil
}

func MigrateUp(dbContext *SqLiteDbContext) error {
	driver, err := sqlite3.WithInstance(dbContext.db, &sqlite3.Config{
		DatabaseName: dbContext.config.EventDb.Name,
	})
	if err != nil {
		return err
	}
	migrations, migrationInitErr := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", support.Resolve(dbContext.config.EventDb.MigrationLoc)),
		dbContext.config.EventDb.Name, driver)

	if migrationInitErr != nil {
		return migrationInitErr
	}

	upgradeErr := migrations.Up()
	if upgradeErr != nil {
		return upgradeErr
	}

	return nil
}

func (sqLiteDbContext *SqLiteDbContext) AppConfig() appConfig.AppConfig {
	return sqLiteDbContext.config
}

func (sqLiteDbContext *SqLiteDbContext) Database() *sql.DB {
	return sqLiteDbContext.db
}
