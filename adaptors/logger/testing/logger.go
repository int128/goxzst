package testing

import (
	"testing"

	"github.com/int128/goxzst/adaptors/logger"
)

// New returns the logger for testing.
func New(t *testing.T) logger.Interface {
	return &testingLogger{t}
}

type testingLogger struct {
	t *testing.T
}

func (l *testingLogger) Logf(format string, v ...interface{}) {
	l.t.Logf(format, v...)
}
