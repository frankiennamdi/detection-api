package repository

import (
	"net"

	"github.com/frankiennamdi/detection-api/db"
	"github.com/frankiennamdi/detection-api/models"
	"github.com/oschwald/geoip2-golang"
)

// provides services for mapping ip to geo point
type MaxMindIPGeoInfoRepository struct {
	maxMindDb *db.MaxMindDb
}

func NewMaxMindIPGeoInfoRepository(maxMindDb *db.MaxMindDb) *MaxMindIPGeoInfoRepository {
	return &MaxMindIPGeoInfoRepository{maxMindDb: maxMindDb}
}

func (repo *MaxMindIPGeoInfoRepository) FindGeoPoint(ip net.IP) (geoPoint *models.GeoPoint, err error) {
	var result *models.GeoPoint

	fnxErr := repo.maxMindDb.WithMaxMindDb(func(db *geoip2.Reader) error {
		city, err := db.City(ip)
		if err != nil {
			return err
		}
		result = &models.GeoPoint{
			Latitude:       city.Location.Latitude,
			Longitude:      city.Location.Longitude,
			AccuracyRadius: city.Location.AccuracyRadius,
		}
		return nil
	})

	if fnxErr != nil {
		return nil, fnxErr
	}

	return result, nil
}
