package main

import (
	"os"

	"github.com/int128/goxzst/di"
)

var version = "HEAD"

func main() {
	cmd := di.NewCmd()
	os.Exit(cmd.Run(os.Args, version))
}
