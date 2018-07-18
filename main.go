package main

import (
	"log"

	_ "github.com/crgimenes/goconfig/toml"
	"github.com/crgimenes/pgfs/config"
	"github.com/crgimenes/pgfs/fuse"
)

func main() {
	err := fuse.Run(config.Get().Mountpoint)
	if err != nil {
		log.Fatal(err)
	}
}
