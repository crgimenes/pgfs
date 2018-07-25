package main

import (
	"log"

	"github.com/crgimenes/pgfs/config"
	"github.com/crgimenes/pgfs/filesystem"
)

func main() {
	config.Load()
	err := filesystem.Mount(config.Get().Mountpoint)
	if err != nil {
		log.Fatal(err)
	}
}
