package db

import (
	"github.com/frankiennamdi/detection-api/config"
	"github.com/frankiennamdi/detection-api/support"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

func TestWithSqLiteDbContext_When_Config_Is_Incomplete(t *testing.T) {
	db := NewSqLiteDb(config.AppConfig{
		EventDb:         config.EventDbConfig{
			File: "bad",
		},
		Server:          config.ServerConfig{},
		IPGeoDbConfig:   config.IPGeoDbConfig{},
		SuspiciousSpeed: 0,
	})

	defer func() {
		err := os.Remove("bad")
		if err != nil {
			log.Printf(support.Warn, err)
		}
	}()

	err := db.WithSqLiteDbContext(func(context *SqLiteDbContext) error {
		err := MigrateUp(context)
		if err != nil {
			return err
		}
		return nil
	}, "mode=rw")
	req := require.New(t)
	req.Error(err)
}
