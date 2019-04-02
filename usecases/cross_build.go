package usecases

import (
	"fmt"
	"os"
	"os/exec"
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
	Logger adaptors.Logger
}

func (u *CrossBuild) Do(in usecases.CrossBuildIn) error {
	cmd := exec.Command("go", "build", "-o", in.OutputFilename)
	cmd.Args = append(cmd.Args, in.GoBuildArgs...)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GOOS=%s", in.Target.GOOS),
		fmt.Sprintf("GOARCH=%s", in.Target.GOARCH))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	u.Logger.Logf("go %s", strings.Join(cmd.Args, " "))
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "go build error")
	}
	return nil
}
