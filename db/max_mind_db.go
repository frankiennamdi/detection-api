package db

import (
	"log"

	appConfig "github.com/frankiennamdi/detection-api/config"
	"github.com/frankiennamdi/detection-api/support"
	"github.com/oschwald/geoip2-golang"
)

// provide services for the MaxMind db
type MaxMindDb struct {
	config                   appConfig.AppConfig
	maxMindDbConnectionLimit chan int
}

type MaxMindDbRequired func(db *geoip2.Reader) error

func NewMaxMindDb(config appConfig.AppConfig) *MaxMindDb {
	maxConnections := 200
	if config.IPGeoDbConfig.MaxConnection > 0 {
		maxConnections = config.IPGeoDbConfig.MaxConnection
	}
	return &MaxMindDb{config: config, maxMindDbConnectionLimit: make(chan int, maxConnections)}
}

func (maxMindDb *MaxMindDb) WithMaxMindDb(fnx MaxMindDbRequired) (err error) {
	maxMindDb.maxMindDbConnectionLimit <- 1
	db, err := geoip2.Open(support.Resolve(maxMindDb.config.IPGeoDbConfig.Location))

	if err != nil {
		<-maxMindDb.maxMindDbConnectionLimit
		return err
	}

	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf(support.Warn, closeErr)
			err = closeErr
		}
		<-maxMindDb.maxMindDbConnectionLimit
	}()

	if err := fnx(db); err != nil {
		return err
	}

	return nil
}
