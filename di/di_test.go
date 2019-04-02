package di_test

import (
	"testing"

	"github.com/int128/goxzst/adaptors/interfaces"
	"github.com/int128/goxzst/di"
)

func TestInvoke(t *testing.T) {
	err := di.Invoke(func(cmd adaptors.Cmd) {
		if cmd == nil {
			t.Errorf("cmd wants non-nin but nil")
		}
	})
	if err != nil {
		t.Errorf("Invoke returned error: %+v", err)
	}
}
