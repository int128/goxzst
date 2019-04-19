//+build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/int128/goxzst/adaptors"
	"github.com/int128/goxzst/usecases"
)

func NewCmd() *adaptors.Cmd {
	wire.Build(
		adaptors.Set,
		usecases.Set,
	)
	return nil
}
