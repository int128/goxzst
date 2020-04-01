//+build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/int128/goxzst/adaptors/cmd"
	"github.com/int128/goxzst/adaptors/env"
	"github.com/int128/goxzst/adaptors/fs"
	"github.com/int128/goxzst/adaptors/logger"
	"github.com/int128/goxzst/usecases/archive"
	"github.com/int128/goxzst/usecases/crossbuild"
	"github.com/int128/goxzst/usecases/digest"
	"github.com/int128/goxzst/usecases/rendertemplate"
	"github.com/int128/goxzst/usecases/xzs"
	"github.com/int128/goxzst/usecases/xzst"
)

func NewCmd() cmd.Interface {
	wire.Build(
		// adaptors
		cmd.Set,
		env.Set,
		fs.Set,
		logger.Set,

		// usecases
		xzst.Set,
		xzs.Set,
		archive.Set,
		crossbuild.Set,
		digest.Set,
		rendertemplate.Set,
	)
	return nil
}
