package db

import (
	appConfig "github.com/frankiennamdi/detection-api/config"
	"github.com/frankiennamdi/detection-api/support"
	"github.com/oschwald/geoip2-golang"
	"log"
)

type MaxMindDb struct {
	config                   appConfig.AppConfig
	maxMindDbConnectionLimit chan int
}

type MaxMindDbRequired func(db *geoip2.Reader) error

func NewMaxMindDb(config appConfig.AppConfig) *MaxMindDb {
	return &MaxMindDb{config: config, maxMindDbConnectionLimit: make(chan int, config.IPGeoDbConfig.MaxConnection)}
}

func (maxMindDb *MaxMindDb) WithMaxMindDb(fnx MaxMindDbRequired) (err error) {
	maxMindDb.maxMindDbConnectionLimit <- 1
	db, err := geoip2.Open(support.Resolve(maxMindDb.config.IPGeoDbConfig.Location))

	if err != nil {
		<-maxMindDb.maxMindDbConnectionLimit
		return err
	}

	defer func() {
		<-maxMindDb.maxMindDbConnectionLimit

		if closeErr := db.Close(); closeErr != nil {
			log.Printf(support.Warn, closeErr)
			err = closeErr
		}
	}()

	if err := fnx(db); err != nil {
		return err
	}

	return nil
}
