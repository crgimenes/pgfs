package main

import (
	"log"

	"github.com/crgimenes/goconfig"
	_ "github.com/crgimenes/goconfig/toml"
	"github.com/crgimenes/pgfs/fuse"
)

type config struct {
	Mountpoint string `json:"m" toml:"mountpoint" cfg:"m"`
}

func main() {
	cfg := config{}
	goconfig.File = "pgfs.toml"
	err := goconfig.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Mountpoint == "" {
		log.Fatalln("mount point is required use -m parameter")
	}

	err = fuse.Run(cfg.Mountpoint)
	if err != nil {
		log.Fatal(err)
	}
}
