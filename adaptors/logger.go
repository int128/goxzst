package adaptors

import (
	"log"

	"github.com/int128/goxzst/adaptors/interfaces"
)

func NewLogger() adaptors.Logger {
	return &Logger{}
}

type Logger struct{}

func (*Logger) Logf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
