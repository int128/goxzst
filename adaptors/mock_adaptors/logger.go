package mock_adaptors

import (
	"testing"

	"github.com/int128/goxzst/adaptors"
)

// NewLogger returns a Logger for testing.
func NewLogger(t *testing.T) adaptors.Logger {
	return &logger{t}
}

type logger struct {
	t *testing.T
}

func (l *logger) Logf(format string, v ...interface{}) {
	l.t.Logf(format, v...)
}
