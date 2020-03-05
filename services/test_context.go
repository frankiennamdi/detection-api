package services

import (
	"github.com/frankiennamdi/detection-api/config"
	"github.com/frankiennamdi/detection-api/support"
	"log"
	"path/filepath"
)

type Test struct {
	appConfig    config.AppConfig
	temporaryDir *support.TemporaryDir
	context      *ServerContext
}

func SetUp() *Test {
	log.Println("setting up")

	temporaryDir := support.NewTemporaryDir("", "sqlite3-driver-test")
	path := support.Resolve("migrations")
	appConfig := config.AppConfig{
		EventDb: config.EventDbConfig{
			File:          filepath.Join(temporaryDir.Path(), "sqlite3.db"),
			Name:          "event_db",
			MigrationLoc:  path,
			MaxConnection: 100,
		},
		IPGeoDbConfig: config.IPGeoDbConfig{
			Location:      "resources/geo-database/GeoLite2-City.mmdb",
			MaxConnection: 100},
		SuspiciousSpeed: 500,
	}

	server := Server{Config: appConfig}
	serverContext := server.Configure()

	return &Test{
		appConfig:    appConfig,
		temporaryDir: temporaryDir,
		context:      serverContext,
	}
}

func (test *Test) CleanUp() {
	log.Println("cleaning up")
	test.temporaryDir.Clean()
}

func (test *Test) AppConfig() config.AppConfig {
	return test.appConfig
}

func (test *Test) ServerContext() *ServerContext {
	return test.context
}
