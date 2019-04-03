package di

import (
	"github.com/int128/goxzst/adaptors"
	"github.com/int128/goxzst/usecases"
)

var dependencies = []interface{}{
	usecases.NewMake,
	usecases.NewCrossBuild,
	usecases.NewArchive,
	usecases.NewDigest,
	usecases.NewRenderTemplate,

	adaptors.NewCmd,
	adaptors.NewLogger,
	adaptors.NewEnv,
	adaptors.NewFilesystem,
}
