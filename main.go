package main

import (
	"log"
	"os"

	"github.com/int128/goxzst/adaptors/interfaces"
	"github.com/int128/goxzst/di"
)

var version = "v0.0.0"

func main() {
	if err := di.Invoke(func(cmd adaptors.Cmd) {
		os.Exit(cmd.Run(os.Args, version))
	}); err != nil {
		log.Fatalf("Error: %+v", err)
	}
}
