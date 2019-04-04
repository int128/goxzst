package adaptors

import (
	"io"
	"os"
)

//go:generate mockgen -package mock_adaptors -destination ../mock_adaptors/mock_adaptors.go github.com/int128/goxzst/adaptors/interfaces Env,FileSystem

type Cmd interface {
	Run(args []string, version string) int
}

type Logger interface {
	Logf(format string, v ...interface{})
}

type Env interface {
	LookupEnv(key string) (string, bool)
	Exec(in ExecIn) error
}

type ExecIn struct {
	Name     string
	Args     []string
	ExtraEnv []string
}

type FileSystem interface {
	Open(name string) (io.ReadCloser, error)
	Create(name string) (io.WriteCloser, error)
	Remove(name string) error
	Stat(name string) (os.FileInfo, error)
	MkdirAll(path string) error
}
