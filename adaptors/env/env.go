package env

import (
	"os"
	"os/exec"

	"github.com/google/wire"
	"github.com/int128/goxzst/adaptors"
	"github.com/pkg/errors"
)

var Set = wire.NewSet(
	wire.Struct(new(Env), "*"),
	wire.Bind(new(adaptors.Env), new(*Env)),
)

type Env struct{}

func (*Env) Getwd() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return dir, nil
}

func (*Env) LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

// Exec runs and waits for a process.
// It inherits env vars of the current process.
// It sets Stdout and Stderr to the os defaults.
func (*Env) Exec(in adaptors.ExecIn) error {
	cmd := exec.Command(in.Name, in.Args...)
	cmd.Env = append(os.Environ(), in.ExtraEnv...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "error while exec")
	}
	return nil
}
