package di

import (
	"github.com/int128/goxzst/adaptors"
	"github.com/int128/goxzst/usecases"
)

var dependencies = []interface{}{
	usecases.NewMake,
	usecases.NewCrossBuild,
	usecases.NewCreateZip,
	usecases.NewCreateSHA,
	usecases.NewRenderTemplate,
	adaptors.NewCmd,
	adaptors.NewLogger,
}
