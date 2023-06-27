package config

import (
	"log"

	"crg.eti.br/go/config"
	_ "crg.eti.br/go/config/ini"
)

// Config stores the configuration
type Config struct {
	Mountpoint     string `json:"m" toml:"mountpoint" cfg:"m" cfgRequired:"true"`
	DataSourceName string `json:"dsn" toml:"dsn" cfg:"dsn" cfgRequired:"true"`
	SchemaName     string `json:"schema" toml:"schema" cfg:"s" cfgDefault:"public" cfgRequired:"true"`
}

var cfg Config

// Load config parameters
func Load() {
	cfg = Config{}
	config.PrefixEnv = "fs"
	config.File = "pgfs.ini"
	err := config.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
}

// Get returns settings
func Get() Config {
	return cfg
}
