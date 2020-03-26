package crossbuild

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/wire"
	"github.com/int128/goxzst/adaptors/env"
	"github.com/int128/goxzst/adaptors/fs"
	"github.com/int128/goxzst/adaptors/logger"
	"github.com/int128/goxzst/models/build"
)

var Set = wire.NewSet(
	wire.Struct(new(CrossBuild), "*"),
	wire.Bind(new(Interface), new(*CrossBuild)),
)

//go:generate mockgen -destination mock_crossbuild/mock_crossbuild.go github.com/int128/goxzst/usecases/crossbuild Interface

type Interface interface {
	Do(in Input) error
}

type Input struct {
	OutputFilename string
	GoBuildArgs    []string
	Platform       build.Platform
}

type CrossBuild struct {
	Env        env.Interface
	FileSystem fs.Interface
	Logger     logger.Interface
}

func (u *CrossBuild) Do(in Input) error {
	if err := u.FileSystem.MkdirAll(filepath.Dir(in.OutputFilename)); err != nil {
		return fmt.Errorf("error while creating the output directory: %w", err)
	}

	args := append([]string{"build", "-o", in.OutputFilename}, in.GoBuildArgs...)
	envVars := []string{
		fmt.Sprintf("GOOS=%s", in.Platform.GOOS),
		fmt.Sprintf("GOARCH=%s", in.Platform.GOARCH),
	}

	u.Logger.Logf("%s go %s", strings.Join(envVars, " "), strings.Join(args, " "))
	if err := u.Env.Exec(env.Exec{
		Name:     "go",
		Args:     args,
		ExtraEnv: envVars,
	}); err != nil {
		return fmt.Errorf("go build error: %w", err)
	}
	return nil
}
