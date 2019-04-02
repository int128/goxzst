package di

import (
	"github.com/int128/goxzst/adaptors/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

func Invoke(f func(cmd adaptors.Cmd)) error {
	c := dig.New()
	for _, d := range dependencies {
		if err := c.Provide(d); err != nil {
			return errors.Wrapf(err, "error while providing %+v", d)
		}
	}
	if err := c.Invoke(f); err != nil {
		return errors.Wrapf(err, "error while invoke")
	}
	return nil
}
