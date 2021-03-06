package repository

import (
	"github.com/frankiennamdi/detection-api/test"
	"net"
	"testing"

	"github.com/frankiennamdi/detection-api/models"
	"github.com/stretchr/testify/require"
)

var IPToGeoPointTestCases = []struct {
	IP               string
	expectedGeoPoint *models.GeoPoint
}{
	{
		IP:               "91.207.175.104",
		expectedGeoPoint: &models.GeoPoint{Latitude: 34.0549, Longitude: -118.2578, AccuracyRadius: 200},
	},
}

func TestFindGeoPoint(t *testing.T) {
	testSetup := test.SetUp()

	defer testSetup.CleanUp()

	req := require.New(t)

	for _, input := range IPToGeoPointTestCases {
		repo := NewMaxMindIPGeoInfoRepository(testSetup.AppServerContext().GeoIPDb())
		geoPoint, err := repo.FindGeoPoint(net.ParseIP(input.IP))
		req.NoError(err)
		req.Equal(input.expectedGeoPoint, geoPoint)
	}
}
