package test

import (
	"github.com/frankiennamdi/detection-api/core"
	"log"
	"path/filepath"

	"github.com/frankiennamdi/detection-api/config"
	"github.com/frankiennamdi/detection-api/support"
)

type Test struct {
	appConfig        config.AppConfig
	temporaryDir     *support.TemporaryDir
	appServerContext *core.ServerContext
}

func SetUp() *Test {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
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

	server := core.Server{Config: appConfig}
	appServerContext := server.Configure()

	return &Test{
		appConfig:        appConfig,
		temporaryDir:     temporaryDir,
		appServerContext: appServerContext,
	}
}

func (test *Test) CleanUp() {
	log.Println("cleaning up")
	test.temporaryDir.Clean()
}

func (test *Test) AppConfig() config.AppConfig {
	return test.appConfig
}

func (test *Test) AppServerContext() *core.ServerContext {
	return test.appServerContext
}
