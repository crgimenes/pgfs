package config

import (
	"log"

	"github.com/crgimenes/goconfig"
)

// Config stores the configuration
type Config struct {
	Mountpoint     string `json:"m" toml:"mountpoint" cfg:"m" cfgRequired:"true"`
	DataSourceName string `json:"dsn" toml:"dsn" cfg:"dsn" cfgRequired:"true"`
	SchemaName     string `json:"schema" toml:"schema" cfg:"s" cfgDefault:"public" cfgRequired:"true"`
}

var cfg Config

func init() {
	if err := load(); err != nil {
		log.Fatal(err)
	}
}

func load() (err error) {
	cfg = Config{}
	goconfig.File = "pgfs.toml"
	err = goconfig.Parse(&cfg)
	return
}

// Get returns settings
func Get() Config {
	return cfg
}
