//+build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/int128/goxzst/adaptors"
	"github.com/int128/goxzst/adaptors/cmd"
	"github.com/int128/goxzst/adaptors/env"
	"github.com/int128/goxzst/adaptors/fs"
	"github.com/int128/goxzst/adaptors/logger"
	"github.com/int128/goxzst/usecases"
	"github.com/int128/goxzst/usecases/archive"
	"github.com/int128/goxzst/usecases/build"
	"github.com/int128/goxzst/usecases/digest"
	"github.com/int128/goxzst/usecases/makeall"
	"github.com/int128/goxzst/usecases/templates"
)

var adaptorsSet = wire.NewSet(
	cmd.Cmd{},
	env.Env{},
	fs.FileSystem{},
	logger.Logger{},
	wire.Bind((*adaptors.Cmd)(nil), (*cmd.Cmd)(nil)),
	wire.Bind((*adaptors.Env)(nil), (*env.Env)(nil)),
	wire.Bind((*adaptors.FileSystem)(nil), (*fs.FileSystem)(nil)),
	wire.Bind((*adaptors.Logger)(nil), (*logger.Logger)(nil)),
)

var usecasesSet = wire.NewSet(
	makeall.Make{},
	archive.Archive{},
	build.CrossBuild{},
	digest.Digest{},
	templates.RenderTemplate{},
	wire.Bind((*usecases.Make)(nil), (*makeall.Make)(nil)),
	wire.Bind((*usecases.Archive)(nil), (*archive.Archive)(nil)),
	wire.Bind((*usecases.CrossBuild)(nil), (*build.CrossBuild)(nil)),
	wire.Bind((*usecases.Digest)(nil), (*digest.Digest)(nil)),
	wire.Bind((*usecases.RenderTemplate)(nil), (*templates.RenderTemplate)(nil)),
)

func NewCmd() adaptors.Cmd {
	wire.Build(
		adaptorsSet,
		usecasesSet,
	)
	return nil
}
