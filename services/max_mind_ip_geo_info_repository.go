package services

import (
	"github.com/frankiennamdi/detection-api/db"
	"github.com/frankiennamdi/detection-api/models"
	"net"
)

func NewMaxMindIPGeoInfoRepository(maxMindDb *db.MaxMindDb) *MaxMindIPGeoInfoRepository {
	return &MaxMindIPGeoInfoRepository{maxMindDb: maxMindDb}
}

func (repo *MaxMindIPGeoInfoRepository) FindGeoPoint(ip net.IP) (geoPoint *models.GeoPoint, err error) {
	context, err := repo.maxMindDb.NewContext()

	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := context.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	city, err := context.Db().City(ip)

	if err != nil {
		return nil, err
	}

	return &models.GeoPoint{
		Latitude:       city.Location.Latitude,
		Longitude:      city.Location.Longitude,
		AccuracyRadius: city.Location.AccuracyRadius,
	}, nil
}
