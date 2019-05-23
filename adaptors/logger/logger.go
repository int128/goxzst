package logger

import (
	"log"

	"github.com/int128/goxzst/adaptors"
)

func NewLogger() adaptors.Logger {
	return &Logger{}
}

type Logger struct{}

func (*Logger) Logf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
