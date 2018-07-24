package main

import (
	"log"

	_ "github.com/crgimenes/goconfig/toml"
	"github.com/crgimenes/pgfs/config"
	"github.com/crgimenes/pgfs/filesystem"
)

func main() {
	err := filesystem.Mount(config.Get().Mountpoint)
	if err != nil {
		log.Fatal(err)
	}
}
