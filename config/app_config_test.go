package config

import (
	"log"
	"os"
	"testing"

	"github.com/frankiennamdi/detection-api/support"
	"github.com/stretchr/testify/require"
)

func TestConfigWithEnvironment(t *testing.T) {
	req := require.New(t)

	setEnv("DB_NAME", "test_db")

	defer unsetEnv("DB_NAME")

	appConfig := &AppConfig{}
	err := appConfig.Read()
	req.NoError(err)
	req.Equal("test_db", appConfig.EventDb.Name)
	req.Equal("resources/event-db/event_db.db", appConfig.EventDb.File)
	req.Equal("migrations", appConfig.EventDb.MigrationLoc)
}

func unsetEnv(key string) {
	if err := os.Unsetenv(key); err != nil {
		log.Panicf(support.Fatal, err)
	}
}

func setEnv(key string, value string) {
	if err := os.Setenv(key, value); err != nil {
		log.Panicf(support.Fatal, err)
	}
}
