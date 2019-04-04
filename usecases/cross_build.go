package usecases

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/int128/goxzst/adaptors/interfaces"
	"github.com/int128/goxzst/usecases/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

func NewCrossBuild(i CrossBuild) usecases.CrossBuild {
	return &i
}

type CrossBuild struct {
	dig.In
	Env        adaptors.Env
	Filesystem adaptors.Filesystem
	Logger     adaptors.Logger
}

func (u *CrossBuild) Do(in usecases.CrossBuildIn) error {
	if err := u.Filesystem.MkdirAll(filepath.Dir(in.OutputFilename)); err != nil {
		return errors.Wrapf(err, "error while creating the output directory")
	}

	args := append([]string{"build", "-o", in.OutputFilename}, in.GoBuildArgs...)
	env := []string{
		fmt.Sprintf("GOOS=%s", in.Platform.GOOS),
		fmt.Sprintf("GOARCH=%s", in.Platform.GOARCH),
	}

	u.Logger.Logf("%s go %s", strings.Join(env, " "), strings.Join(args, " "))
	if err := u.Env.Exec(adaptors.ExecIn{
		Name:     "go",
		Args:     args,
		ExtraEnv: env,
	}); err != nil {
		return errors.Wrapf(err, "go build error")
	}
	return nil
}
