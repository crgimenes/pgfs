package main

import (
	"log"

	"github.com/crgimenes/goconfig"
	"github.com/crgimenes/pgfs/fuse"
)

type config struct {
	Mountpoint string `json:"m" cfg:"m"`
}

func main() {
	cfg := config{}
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
