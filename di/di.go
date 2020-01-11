//go:generate wire
//+build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/int128/goxzst/adaptors"
	"github.com/int128/goxzst/adaptors/cmd"
	"github.com/int128/goxzst/adaptors/env"
	"github.com/int128/goxzst/adaptors/fs"
	"github.com/int128/goxzst/adaptors/logger"
	"github.com/int128/goxzst/usecases/archive"
	"github.com/int128/goxzst/usecases/build"
	"github.com/int128/goxzst/usecases/digest"
	"github.com/int128/goxzst/usecases/makeall"
	"github.com/int128/goxzst/usecases/templates"
)

func NewCmd() adaptors.Cmd {
	wire.Build(
		// adaptors
		cmd.Set,
		env.Set,
		fs.Set,
		logger.Set,

		// usecases
		makeall.Set,
		archive.Set,
		build.Set,
		digest.Set,
		templates.Set,
	)
	return nil
}
