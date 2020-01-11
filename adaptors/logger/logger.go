package logger

import (
	"log"

	"github.com/google/wire"
	"github.com/int128/goxzst/adaptors"
)

var Set = wire.NewSet(
	wire.Struct(new(Logger), "*"),
	wire.Bind(new(adaptors.Logger), new(*Logger)),
)

type Logger struct{}

func (*Logger) Logf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
