package logger

import (
	"log"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	wire.Struct(new(Logger), "*"),
	wire.Bind(new(Interface), new(*Logger)),
)

type Interface interface {
	Logf(format string, v ...interface{})
}

func init() {
	log.SetFlags(log.Lmicroseconds)
}

type Logger struct{}

func (*Logger) Logf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
