package services

import (
	"fmt"
	"github.com/frankiennamdi/detection-api/config"
	"github.com/frankiennamdi/detection-api/db"
	"github.com/frankiennamdi/detection-api/support"
	_ "github.com/golang-migrate/migrate/v4"                  // indirect
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3" // indirect
	_ "github.com/golang-migrate/migrate/v4/source/file"      // indirect
	"log"
	"net/http"
)

type Server struct {
	Config config.AppConfig
}

type ServerContext struct {
	sqLitDb          *db.SqLiteDb
	maxMindDb        *db.MaxMindDb
	config           config.AppConfig
	detectionService DetectionService
	eventRepository  EventRepository
}

func (serverContext *ServerContext) EventDb() *db.SqLiteDb {
	return serverContext.sqLitDb
}

func (serverContext *ServerContext) MaxMindDb() *db.MaxMindDb {
	return serverContext.maxMindDb
}

func (server *Server) Configure() *ServerContext {
	sqLiteDb := db.NewSqLiteDb(server.Config)
	err := server.configureEventDb(sqLiteDb)

	if err != nil {
		log.Panicf(support.Fatal, err)
	}

	maxMindDb := db.NewMaxMindDb(server.Config)
	eventRepository := NewSQLLiteEventsRepository(sqLiteDb)
	detectionService := NewLoginDetectionService(eventRepository,
		NewMaxMindIPGeoInfoRepository(maxMindDb),
		DefaultCalculatorService{},
		server.Config.SuspiciousSpeed)

	return &ServerContext{
		sqLitDb:          sqLiteDb,
		maxMindDb:        maxMindDb,
		config:           server.Config,
		detectionService: detectionService,
		eventRepository:  eventRepository,
	}
}

func (serverContext *ServerContext) Listen() {
	log.Printf(support.Info, fmt.Sprintf("service starting on Port : %d ...", serverContext.config.Server.Port))

	router := Router{serverContext: serverContext}
	routes := router.InitRoutes()

	if err := http.ListenAndServe(fmt.Sprintf(":%d", serverContext.config.Server.Port), routes); err != nil {
		log.Fatal(err)
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
