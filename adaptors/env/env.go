package env

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	wire.Struct(new(Env), "*"),
	wire.Bind(new(Interface), new(*Env)),
)

//go:generate mockgen -destination mock_env/mock_env.go github.com/int128/goxzst/adaptors/env Interface

type Interface interface {
	LookupEnv(key string) (string, bool)
	Exec(in Exec) error
}

type Exec struct {
	Name     string
	Args     []string
	ExtraEnv []string
}

type Env struct{}

func (*Env) Getwd() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	return dir, nil
}

func (*Env) LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

// Exec runs and waits for a process.
// It inherits env vars of the current process.
// It sets Stdout and Stderr to the os defaults.
func (*Env) Exec(e Exec) error {
	cmd := exec.Command(e.Name, e.Args...)
	cmd.Env = append(os.Environ(), e.ExtraEnv...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error while exec: %w", err)
	}
	return nil
}
