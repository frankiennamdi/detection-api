package db

import (
	"database/sql"
	"fmt"
	"log"

	appConfig "github.com/frankiennamdi/detection-api/config"
	"github.com/frankiennamdi/detection-api/support"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
)

// provide services for SQLite db
type SqLiteDb struct {
	config                  appConfig.AppConfig
	sqLiteDbConnectionLimit chan int
}

type SqLiteDbContext struct {
	db     *sql.DB
	config appConfig.AppConfig
}

type TransactionEnabled func(tx *sql.Tx) error

type SQLiteDbRequired func(context *SqLiteDbContext) error

func NewSqLiteDb(config appConfig.AppConfig) *SqLiteDb {
	return &SqLiteDb{config: config, sqLiteDbConnectionLimit: make(chan int, config.EventDb.MaxConnection)}
}

func (sqLiteDb *SqLiteDb) WithSqLiteDbContext(fnx SQLiteDbRequired) (err error) {
	sqLiteDb.sqLiteDbConnectionLimit <- 1
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?%s", sqLiteDb.config.EventDb.File, "mode=rwc"))

	if err != nil {
		<-sqLiteDb.sqLiteDbConnectionLimit
		return err
	}

	defer func() {
		if closeErr := db.Close(); err != nil {
			log.Printf(support.Warn, closeErr)
			<-sqLiteDb.sqLiteDbConnectionLimit

			err = closeErr
		} else {
			<-sqLiteDb.sqLiteDbConnectionLimit
		}
	}()

	fnxErr := fnx(&SqLiteDbContext{db: db,
		config: sqLiteDb.config,
	})

	if fnxErr != nil {
		return fnxErr
	}

	return nil
}

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
