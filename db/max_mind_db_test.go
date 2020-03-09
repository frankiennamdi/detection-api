package db

import (
	"github.com/frankiennamdi/detection-api/config"
	"github.com/oschwald/geoip2-golang"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWithMaxMindDb_When_File_Does_Not_Exit(t *testing.T) {
	db := NewMaxMindDb(config.AppConfig{
		EventDb: config.EventDbConfig{},
		Server:  config.ServerConfig{},
		IPGeoDbConfig: config.IPGeoDbConfig{
			Location: "bad",
		},
		SuspiciousSpeed: 0,
	})

	err := db.WithMaxMindDb(func(db *geoip2.Reader) error {
		return nil
	});
	req := require.New(t)
	req.Error(err)
}
