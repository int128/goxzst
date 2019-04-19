// Package adaptors provides bridge between use-cases and external infrastructure.
package adaptors

import (
	"github.com/google/wire"
	"github.com/int128/goxzst/adaptors/interfaces"
)

var Set = wire.NewSet(
	Cmd{},
	Env{},
	FileSystem{},
	Logger{},
	wire.Bind((*adaptors.Env)(nil), (*Env)(nil)),
	wire.Bind((*adaptors.FileSystem)(nil), (*FileSystem)(nil)),
	wire.Bind((*adaptors.Logger)(nil), (*Logger)(nil)),
)
