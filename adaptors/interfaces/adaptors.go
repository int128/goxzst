package adaptors

import (
	"io"
	"os"
)

//go:generate mockgen -package mock_adaptors -destination ../mock_adaptors/mock_adaptors.go github.com/int128/goxzst/adaptors/interfaces Env,Filesystem

type Cmd interface {
	Run(args []string) int
}

type Logger interface {
	Logf(format string, v ...interface{})
}

type Env interface {
	Getwd() (string, error)
	Exec(in ExecIn) error
}

type ExecIn struct {
	Name     string
	Args     []string
	ExtraEnv []string
}

type Filesystem interface {
	Open(name string) (io.ReadCloser, error)
	Create(name string) (io.WriteCloser, error)
	GetMode(name string) (os.FileMode, error)
	MkdirAll(path string) error
}
