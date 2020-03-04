package db

import (
	appConfig "github.com/frankiennamdi/detection-api/config"
	"github.com/frankiennamdi/detection-api/support"
	"github.com/oschwald/geoip2-golang"
)

type MaxMindDb struct {
	Config appConfig.AppConfig
}

type MaxMindDbContext struct {
	db *geoip2.Reader
}

func (maxMindDb *MaxMindDb) NewContext() (*MaxMindDbContext, error) {
	db, err := geoip2.Open(support.Resolve(maxMindDb.Config.IPGeoDbConfig.Location))

	if err != nil {
		return nil, err
	}

	return &MaxMindDbContext{db: db}, err
}

func (maxMindDbContext *MaxMindDbContext) Close() error {
	if err := maxMindDbContext.db.Close(); err != nil {
		return err
	}

	return nil
}

func (maxMindDbContext *MaxMindDbContext) Db() *geoip2.Reader {
	return maxMindDbContext.db
}
