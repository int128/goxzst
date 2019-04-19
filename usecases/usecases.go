// Package usecases provides use-cases.
package usecases

import (
	"github.com/google/wire"
	"github.com/int128/goxzst/usecases/interfaces"
)

var Set = wire.NewSet(
	Make{},
	Archive{},
	CrossBuild{},
	Digest{},
	RenderTemplate{},
	wire.Bind((*usecases.Make)(nil), (*Make)(nil)),
	wire.Bind((*usecases.Archive)(nil), (*Archive)(nil)),
	wire.Bind((*usecases.CrossBuild)(nil), (*CrossBuild)(nil)),
	wire.Bind((*usecases.Digest)(nil), (*Digest)(nil)),
	wire.Bind((*usecases.RenderTemplate)(nil), (*RenderTemplate)(nil)),
)
