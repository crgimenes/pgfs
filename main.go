package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/crgimenes/pgfs/fuse"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s MOUNTPOINT\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
		os.Exit(2)
	}
	mountpoint := flag.Arg(0)

	err := fuse.Run(mountpoint)
	if err != nil {
		log.Fatal(err)
	}
}
