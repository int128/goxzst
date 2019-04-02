package main

import (
	"log"
	"os"

	"github.com/int128/goxzst/adaptors/interfaces"
	"github.com/int128/goxzst/di"
)

func main() {
	if err := di.Invoke(func(cmd adaptors.Cmd) {
		os.Exit(cmd.Run(os.Args))
	}); err != nil {
		log.Fatalf("Error: %+v", err)
	}
}
