package config

import (
	"github.com/frankiennamdi/detection-api/support"
	"github.com/frankiennamdi/go-configuration/configuration"
	"io/ioutil"
)

type EventDbConfig struct {
	File          string `config:"file"`
	Name          string `config:"name"`
	MigrationLoc  string `config:"migrationLoc"`
	MaxConnection int    `config:"maxConnection"`
}

type IPGeoDbConfig struct {
	Location      string `config:"location"`
	MaxConnection int    `config:"maxConnection"`
}

type ServerConfig struct {
	Port int `config:"port"`
}

type AppConfig struct {
	EventDb         EventDbConfig `config:"eventDb"`
	Server          ServerConfig  `config:"server"`
	IPGeoDbConfig   IPGeoDbConfig `config:"ipGeoDbConfig"`
	SuspiciousSpeed float64       `config:"suspiciousSpeed"`
}

func (appConfig *AppConfig) Read() error {
	configFile := support.GetEnvironmentValue("CONFIG_FILE", "resources/config.yml")
	yamlData, err := ioutil.ReadFile(support.Resolve(configFile))

	if err != nil {
		panic(err)
	}

	binder := configuration.New()
	if err := binder.InitializeConfigFromYaml(yamlData, appConfig); err != nil {
		return err
	}

	return nil
}
