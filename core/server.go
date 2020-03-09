package core

import (
	"log"

	"github.com/frankiennamdi/detection-api/config"
	"github.com/frankiennamdi/detection-api/db"
	"github.com/frankiennamdi/detection-api/support"
	_ "github.com/golang-migrate/migrate/v4"                  // indirect
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3" // indirect
	_ "github.com/golang-migrate/migrate/v4/source/file"      // indirect
)

// represent the server, the starting point of the system
type Server struct {
	config config.AppConfig
}

func NewServer(config config.AppConfig) *Server {
	return &Server{config: config}
}

// represents the initialized and configures core components of the system like database, and application configuration
type ServerContext struct {
	sqLitDb   *db.SqLiteDb
	maxMindDb *db.MaxMindDb
	config    config.AppConfig
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
	sqLiteDb := db.NewSqLiteDb(server.config)
	err := server.configureEventDb(sqLiteDb)

	if err != nil {
		log.Panicf(support.Fatal, err)
	}

	maxMindDb := db.NewMaxMindDb(server.config)

	return &ServerContext{
		sqLitDb:   sqLiteDb,
		maxMindDb: maxMindDb,
		config:    server.config,
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

	if !support.FileExists(server.config.EventDb.File) {
		log.Panic("database file does not exist")
	}

	log.Printf(support.Info, "configuring EventDb Complete")

	return fnxErr
}
