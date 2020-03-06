package core

import (
	"github.com/frankiennamdi/detection-api/config"
	"github.com/frankiennamdi/detection-api/db"
	"github.com/frankiennamdi/detection-api/support"
	_ "github.com/golang-migrate/migrate/v4"                  // indirect
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3" // indirect
	_ "github.com/golang-migrate/migrate/v4/source/file"      // indirect
	"log"
)

type Server struct {
	Config config.AppConfig
}

type ServerContext struct {
	sqLitDb          *db.SqLiteDb
	maxMindDb        *db.MaxMindDb
	config           config.AppConfig
}

func (serverContext *ServerContext) EventDb() *db.SqLiteDb {
	return serverContext.sqLitDb
}

func (serverContext *ServerContext) GeoIPDb() *db.MaxMindDb {
	return serverContext.maxMindDb
}

func (serverContext *ServerContext) AppConfig() config.AppConfig {
	return serverContext.config
}

func (server *Server) Configure() *ServerContext {
	sqLiteDb := db.NewSqLiteDb(server.Config)
	err := server.configureEventDb(sqLiteDb)

	if err != nil {
		log.Panicf(support.Fatal, err)
	}
	maxMindDb := db.NewMaxMindDb(server.Config)
	return &ServerContext{
		sqLitDb:          sqLiteDb,
		maxMindDb:        maxMindDb,
		config:           server.Config,
	}
}

func (server *Server) configureEventDb(sqLiteDb *db.SqLiteDb) (err error) {
	log.Printf(support.Info, "configuring EventDb")

	fnxErr := sqLiteDb.WithSqLiteDbContext(func(context *db.SqLiteDbContext) error {
		migrationErr := db.MigrateUp(context)
		if migrationErr != nil {
			return err
		}
		return nil
	})

	log.Printf(support.Info, "configuring EventDb Complete")

	return fnxErr
}
